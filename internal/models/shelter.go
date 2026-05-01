package models

import (
	"errors"
	"strings"
	"time"
)

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

type UpdateShelterInput struct {
	Name      *string  `json:"name"`
	Address   *string  `json:"address"`
	Email     *string  `json:"email"`
	Phone     *string  `json:"phone_number"`
	Latitude  *float64 `json:"latitude"`
	Longitude *float64 `json:"longitude"`
}

func (usi UpdateShelterInput) Validate() error {
	validationErrors := make(map[string]string)

	if usi.Name != nil && *usi.Name == "" {
		validationErrors["name"] = "cannot be empty"
	}
	if usi.Address != nil && *usi.Address == "" {
		validationErrors["address"] = "cannot be empty"
	}
	if usi.Longitude != nil && (*usi.Longitude < -180 || *usi.Longitude > 180) {
		validationErrors["longitude"] = "must be between -180 and 180"
	}
	if usi.Latitude != nil && (*usi.Latitude < -90 || *usi.Latitude > 90) {
		validationErrors["latitude"] = "must be between -90 and 90"
	}

	if len(validationErrors) > 0 {
		msgs := make([]string, 0, len(validationErrors))
		for field, msg := range validationErrors {
			msgs = append(msgs, field+": "+msg)
		}
		return errors.New(strings.Join(msgs, ", "))
	}
	return nil
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

func (csi CreateShelterInput) Validate() error {
	validationErrors := make(map[string]string)
	if csi.Longitude != nil && (*csi.Longitude < -180 || *csi.Longitude > 180) {
		validationErrors["longitude"] = "must be between -180 and 180"
	}
	if csi.Latitude != nil && (*csi.Latitude < -90 || *csi.Latitude > 90) {
		validationErrors["latitude"] = "must be between -90 and 90"
	}
	if len(validationErrors) > 0 {
		msgs := make([]string, 0, len(validationErrors))
		for field, msg := range validationErrors {
			msgs = append(msgs, field+": "+msg)
		}
		return errors.New(strings.Join(msgs, ", "))
	}
	return nil
}
