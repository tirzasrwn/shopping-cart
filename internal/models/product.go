package models

import "time"

type Product struct {
	ID          int        `json:"id,omitempty"`
	CategoryID  int        `json:"category_id,omitempty"`
	Name        string     `json:"name"`
	Description string     `json:"description,omitempty"`
	Price       float64    `json:"price"`
	CreatedAt   *time.Time `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
}

type ProductOrder struct {
	Product
	ID        int        `json:"product_id"`
	OrderID   int        `json:"order_id"`
	Quantity  int        `json:"quantity"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

type ProductPayment struct {
	Product
	ID        int        `json:"product_id"`
	PaymentID int        `json:"payment_id"`
	Quantity  int        `json:"quantity"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}
