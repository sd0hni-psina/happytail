package models

import (
	"errors"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"
)

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

type UserPublic struct {
	ID          int       `json:"id"`
	FullName    string    `json:"full_name"`
	Email       string    `json:"email"`
	PhoneNumber string    `json:"phone_number"`
	City        *string   `json:"city"`
	Points      int       `json:"points"`
	CreatedAt   time.Time `json:"created_at"`
}

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (cui *CreateUserInput) Validate() error {
	ValidationErrors := make(map[string]string)
	hasDigit := false
	hasLetter := false

	if cui.FullName == "" {
		ValidationErrors["full name"] = "name is required"
	}
	if cui.Email == "" || !strings.Contains(cui.Email, "@") || !strings.Contains(cui.Email, ".") {
		ValidationErrors["email"] = "email is required"
	}
	if utf8.RuneCountInString(cui.Password) < 8 {
		ValidationErrors["password"] = "password must be have > 8 symbols"
	}
	for _, ch := range cui.Password {
		if unicode.IsDigit(ch) {
			hasDigit = true
		}
		if unicode.IsLetter(ch) {
			hasLetter = true
		}
	}
	if !hasDigit || !hasLetter {
		ValidationErrors["password"] = "password must contain at least one digit and one letter"
	}

	if cui.PhoneNumber == "" {
		ValidationErrors["phone number"] = "phone number is required"
	}

	if len(ValidationErrors) > 0 {
		msgs := make([]string, 0, len(ValidationErrors))
		for fields, msg := range ValidationErrors {
			msgs = append(msgs, fields+":", msg)
		}
		return errors.New(strings.Join(msgs, ", "))
	}
	return nil
}

func (u *User) ToPublic() UserPublic {
	return UserPublic{
		ID:          u.ID,
		FullName:    u.FullName,
		Email:       u.Email,
		PhoneNumber: u.PhoneNumber,
		City:        u.City,
		Points:      u.Points,
		CreatedAt:   u.CreatedAt,
	}
}
