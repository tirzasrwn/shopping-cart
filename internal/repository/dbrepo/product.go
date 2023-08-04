package dbrepo

import (
	"context"

	"github.com/tirzasrwn/shopping-cart/internal/models"
)

func (m *PostgresDBRepo) GetProductByCategory(id int) ([]*models.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select id, category_id, name, description, price, created_at, updated_at from public.product where category_id = $1`

	var products []*models.Product
	rows, err := m.DB.QueryContext(ctx, query, id)
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
