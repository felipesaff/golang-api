package models

type User struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Document  string `json:"document"`
	CreatedAt string `json:"created_at"`
}
