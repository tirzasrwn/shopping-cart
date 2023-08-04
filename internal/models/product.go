package models

import "time"

type Product struct {
	ID          int        `json:"id"`
	CategoryID  int        `json:"category_id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Price       float64    `json:"price"`
	CreatedAt   *time.Time `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
}
