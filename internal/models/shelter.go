package models

import "time"

type Shelter struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Address   string    `json:"address"`
	Email     *string   `json:"email"`
	Phone     *string   `json:"phone_number"`
	Latitude  *float64  `json:"latitude"`
	Longitude *float64  `json:"longitude"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateShelterInput struct {
	Name      string   `json:"name"`
	Address   string   `json:"address"`
	Email     *string  `json:"email"`
	Phone     *string  `json:"phone_number"`
	Latitude  *float64 `json:"latitude"`
	Longitude *float64 `json:"longitude"`
}

type ShelterWithDistance struct {
	Shelter
	Distance float64 `json:"distance_km"`
}

type NearbyParams struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	RadiusKm  float64 `json:"radius_km"`
}
