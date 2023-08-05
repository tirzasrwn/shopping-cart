package models

import "time"

type Product struct {
	ID          int        `json:"id,omitempty"`
	CategoryID  int        `json:"category_id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
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
