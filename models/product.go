package models

type Product struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	PriceCents  int64  `json:"price_cents"`
	Stock       int    `json:"stock"`
	CreatedAt   string `json:"created_at"`
}
