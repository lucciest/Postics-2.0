package models

import (
	"time"
)

type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Password  string    `json:"-"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func NewUser(username, password, email string) *User {
	now := time.Now()
	return &User{
		Username:  username,
		Password:  password,
		Email:     email,
		CreatedAt: now,
		UpdatedAt: now,
	}
}
