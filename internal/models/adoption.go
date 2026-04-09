package models

import "time"

type Adoption struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	AnimalID  int       `json:"animal_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateAdoptionInput struct {
	AnimalID int `json:"animal_id"`
}
