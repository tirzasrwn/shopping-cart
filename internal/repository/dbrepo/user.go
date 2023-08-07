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
func (m *PostgresDBRepo) InsertUser(user *models.User) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `insert into public.user (email, password) values ($1, $2) returning id`

	row := m.DB.QueryRowContext(ctx, query,
		user.Email,
		user.Password,
	)
	var userID int
	err := row.Scan(&userID)
	if err != nil {
		return 0, err
	}

	return userID, nil
}

func (m *PostgresDBRepo) GetUserCartByEmail(email string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select c.id from public.user u inner join public.cart c on (u.id = c.user_id) where email = $1`

	var cartID int
	row := m.DB.QueryRowContext(ctx, query, email)

	err := row.Scan(
		&cartID,
	)

	if err != nil {
		return 0, err
	}

	return cartID, nil

}

func (m *PostgresDBRepo) CreateCartForNewUser(userID int) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `insert into public.cart (user_id) values ($1) returning id`

	var cartID int
	row := m.DB.QueryRowContext(ctx, query, userID)
	err := row.Scan(
		&cartID,
	)
	if err != nil {
		return 0, err
	}

	return cartID, nil
}

// GetUserCartIDByEmail return user cart_id.
func (m *PostgresDBRepo) GetUserCartIDByEmail(email string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select (c.id) from public.cart c inner join public.user u on (c.user_id = u.id) where email = $1`

	var cartID int
	row := m.DB.QueryRowContext(ctx, query, email)
	err := row.Scan(
		&cartID,
	)
	if err != nil {
		return 0, err
	}

	return cartID, nil
}
