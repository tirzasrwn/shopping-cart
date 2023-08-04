package dbrepo

import (
	"context"

	"github.com/tirzasrwn/shopping-cart/internal/models"
)

func (m *PostgresDBRepo) GetCategories() ([]*models.Category, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select id, name, created_at, updated_at from public.category`

	var categories []*models.Category
	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var category models.Category
		err := rows.Scan(
			&category.ID,
			&category.Name,
			&category.CreatedAt,
			&category.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		categories = append(categories, &category)
	}

	return categories, nil
}
