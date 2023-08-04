package dbrepo

import (
	"context"

	"github.com/tirzasrwn/shopping-cart/internal/models"
)

// GetUserByEmail returns one user, by email.
func (m *PostgresDBRepo) GetUserByEmail(email string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select id, email, password, created_at, updated_at from public.user where email = $1`

	var user models.User
	row := m.DB.QueryRowContext(ctx, query, email)

	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// InsertUser creates new user.
func (m *PostgresDBRepo) InsertUser(user *models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `insert into public.user (email, password)
			values ($1, $2)`

	_, err := m.DB.ExecContext(ctx, stmt,
		user.Email,
		user.Password,
	)

	if err != nil {
		return err
	}

	return nil
}
