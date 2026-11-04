package models

type Order struct {
	ID         int    `json:"id"`
	UserID     int    `json:"user_id"`
	TotalCents int64  `json:"total_cents"`
	Status     string `json:"status"`
	CreatedAt  string `json:"created_at"`
}
