package models

import "time"

type User struct {
	ID        int64
	Name      string
	Email     string
	Document  string
	CreatedAt time.Time
}
