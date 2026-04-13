package models

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

type AnimalType string

const (
	AnimalTypeDog   AnimalType = "dog"
	AnimalTypeCat   AnimalType = "cat"
	AnimalTypeBunny AnimalType = "bunny"
)

var validAnimalTypes = map[AnimalType]struct{}{
	AnimalTypeDog:   {},
	AnimalTypeCat:   {},
	AnimalTypeBunny: {},
}

func (at AnimalType) Normalize() AnimalType {
	return AnimalType(strings.ToLower(string(at)))
}
func (at AnimalType) IsValid() bool {
	_, ok := validAnimalTypes[at]
	return ok
}

type Animal struct {
	ID             int        `json:"id"`
	Type           AnimalType `json:"type"`
	Name           string     `json:"name"`
	Age            *int       `json:"age"`
	Breed          *string    `json:"breed"`
	Color          *string    `json:"color"`
	IsVaccinated   bool       `json:"is_vaccinated"`
	HasVetPassport bool       `json:"has_vet_passport"`

	Description *string   `json:"description"`
	ShelterID   *int      `json:"shelter_id"`
	Status      string    `json:"status"`
	ShareCount  int       `json:"share_count"`
	CreatedAt   time.Time `json:"created_at"`
}

type CreateAnimalInput struct {
	Type           AnimalType `json:"animal_type"`
	Name           string     `json:"name"`
	Age            *int       `json:"age"`
	Breed          *string    `json:"breed"`
	Color          *string    `json:"color"`
	IsVaccinated   bool       `json:"is_vaccinated"`
	HasVetPassport bool       `json:"has_vet_passport"`

	Description *string `json:"description"`
	ShelterID   *int    `json:"shelter_id"`
}

func (cai *CreateAnimalInput) Validate() error {
	cai.Type = cai.Type.Normalize()
	validationErrors := make(map[string]string)

	if cai.Type == "" {
		validationErrors["type"] = "type is required"
	} else if !cai.Type.IsValid() {
		validationErrors["type"] = fmt.Sprintf("invalid type: %s", cai.Type)
	}
	if cai.Name == "" {
		validationErrors["name"] = "name is required"
	}
	if cai.Age != nil && *cai.Age < 0 {
		validationErrors["age"] = "age cannot be negative"
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
