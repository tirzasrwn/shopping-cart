package dbrepo

import (
	"context"
	"time"

	"github.com/tirzasrwn/shopping-cart/internal/models"
)

// InsertOrder creates new order.
func (m *PostgresDBRepo) InsertOrder(cartID int, productID int, quantity int) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `insert into public.order (cart_id, product_id, quantity) values ($1, $2, $3)
  returning id`

	var orderID int
	row := m.DB.QueryRowContext(ctx, query, cartID, productID, quantity)
	err := row.Scan(&orderID)
	if err != nil {
		return 0, err
	}

	return orderID, nil
}

// CheckOrderExist checks order exist based on card id and product id.
func (m *PostgresDBRepo) CheckOrderExist(cartID int, productID int) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select exists(select * from public.order where cart_id = $1 and product_id = $2)`

	var isExist bool
	row := m.DB.QueryRowContext(ctx, query, cartID, productID)
	err := row.Scan(&isExist)
	if err != nil {
		return false, err
	}

	return isExist, nil
}

// UpdateQuantity adds more quantity from the same order.
func (m *PostgresDBRepo) UpdateQuantity(cartID int, productID int, quantityToAdd int) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `update public.order set quantity = quantity + $1, updated_at = $2
  where cart_id = $3 and product_id = $4
  returning id`

	var orderID int
	row := m.DB.QueryRowContext(ctx, query, quantityToAdd, time.Now(), cartID, productID)
	err := row.Scan(&orderID)
	if err != nil {
		return 0, err
	}

	return orderID, nil
}

// GetOrderByUserEmail returns products in specific user cart.
func (m *PostgresDBRepo) GetOrderByUserEmail(email string) ([]*models.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select p.id, p.category_id, p.name, p.description, p.price, p.created_at, p.updated_at
  from public.order o inner join product p on (o.product_id = p.id) 
  where cart_id in (
    select id from cart where id in (
      select id from public.user where email = $1
    )
  )`

	var products []*models.Product
	rows, err := m.DB.QueryContext(ctx, query, email)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var product models.Product
		err := rows.Scan(
			&product.ID,
			&product.CategoryID,
			&product.Name,
			&product.Description,
			&product.Price,
			&product.CreatedAt,
			&product.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		products = append(products, &product)
	}
	return products, nil
}
