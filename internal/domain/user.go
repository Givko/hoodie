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
	pwd.Password = password
	pwd.Salt = salt

	return &pwd
}

type Password struct {
	Password string
	Salt     string
}
