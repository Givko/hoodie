package domain

import "time"

type User struct {
	Id        string
	Username  string
	Password  Password
	IsAdmin   bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewPassword(password string, salt string) *Password {
	pwd := Password{}
	pwd.password = password
	pwd.salt = salt

	return &pwd
}

type Password struct {
	password string
	salt     string
}
