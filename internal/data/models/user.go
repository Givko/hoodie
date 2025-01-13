package models

import "time"

type User struct {
	Id        string
	Username  string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}
