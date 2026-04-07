package models

import "time"

type Animal struct {
	ID             int     `json:"id"`
	Type           string  `json:"type"`
	Name           string  `json:"name"`
	Age            *int    `json:"age"`
	Breed          *string `json:"breed"`
	Color          *string `json:"color"`
	IsVaccinated   bool    `json:"is_vaccinated"`
	HasVetPassport bool    `json:"has_vet_passport"`

	Description *string   `json:"description"`
	ShelterID   *int      `json:"shelter_id"`
	Status      string    `json:"status"`
	ShareCount  int       `json:"share_count"`
	CreatedAt   time.Time `json:"created_at"`
}
