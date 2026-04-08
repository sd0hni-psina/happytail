package models

import "time"

type User struct {
	ID           int       `json:"id"`
	FullName     string    `json:"full_name"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	PhoneNumber  string    `json:"phone_number"`
	City         *string   `json:"city"`
	Points       int       `json:"points"`
	CreatedAt    time.Time `json:"created_at"`
}

type CreateUserInput struct {
	FullName    string  `json:"full_name"`
	Email       string  `json:"email"`
	Password    string  `json:"password"`
	PhoneNumber string  `json:"phone_number"`
	City        *string `json:"city"`
}

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
