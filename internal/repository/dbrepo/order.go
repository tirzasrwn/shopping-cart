package dbrepo

import (
	"context"
	"time"
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
