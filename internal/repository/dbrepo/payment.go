package dbrepo

import (
	"context"

	"github.com/tirzasrwn/shopping-cart/internal/models"
)

func (m *PostgresDBRepo) InsertSuccessPayment(userID int, productID int, quantity int) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `insert into public.payment(user_id, product_id, quantity) values ($1, $2, $3) returning id`

	var paymentID int
	row := m.DB.QueryRowContext(ctx, query, userID, productID, quantity)
	err := row.Scan(&paymentID)
	if err != nil {
		return 0, err
	}

	return paymentID, nil
}

func (m *PostgresDBRepo) GetPaymentByUserEmail(email string) ([]*models.ProductPayment, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select p.id, p.name, p.price, m.id, m.quantity, m.created_at, m.updated_at
from public.payment m inner join public.product p on (m.product_id = p.id) 
where user_id in (
    select id from public.user where email = $1
)`
	rows, err := m.DB.QueryContext(ctx, query, email)
	if err != nil {
		return nil, err
	}
	var productPayments []*models.ProductPayment
	for rows.Next() {
		var productPayment models.ProductPayment
		err := rows.Scan(
			&productPayment.ID,
			&productPayment.Name,
			&productPayment.Price,
			&productPayment.PaymentID,
			&productPayment.Quantity,
			&productPayment.CreatedAt,
			&productPayment.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		productPayments = append(productPayments, &productPayment)
	}
	return productPayments, nil
}
