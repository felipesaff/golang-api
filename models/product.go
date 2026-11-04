package models

import "time"

type Product struct {
	ID          int64
	Name        string
	Description string
	PriceCents  int64
	Stock       int
	CreatedAt   time.Time
}
