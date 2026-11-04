package user

import "time"

type UserResponse struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Document  string    `json:"document"`
	CreatedAt time.Time `json:"created_at"`
}
