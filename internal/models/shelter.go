package models

import "time"

type Shelter struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Address   string    `json:"address"`
	Email     *string   `json:"email"`
	Phone     *string   `json:"phone_number"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateShelterInput struct {
	Name    string  `json:"name"`
	Address string  `json:"address"`
	Email   *string `json:"email"`
	Phone   *string `json:"phone_number"`
}
