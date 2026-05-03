package models

import (
	"net/http"
	"strconv"
)

type FilterParams struct {
	Name           *string `json:"name"`
	Type           *string `json:"type"`
	Breed          *string `json:"breed"`
	Color          *string `json:"color"`
	IsVaccinated   *bool   `json:"is_vaccinated"`
	HasVetPassport *bool   `json:"has_vet_passport"`
	Status         *string `json:"status"`
	ShelterID      *int    `json:"shelter_id"`
}

func ParseFilter(r *http.Request) FilterParams {
	params := FilterParams{}
	if v := r.URL.Query().Get("name"); v != "" {
		wrapped := "%" + v + "%"
		params.Name = &wrapped
	}

	animalType := r.URL.Query().Get("type")
	if animalType != "" {
		params.Type = &animalType
	}
	breed := r.URL.Query().Get("breed")
	if breed != "" {
		params.Breed = &breed
	}
	color := r.URL.Query().Get("color")
	if color != "" {
		params.Color = &color
	}
	IsVaccinated, err := strconv.ParseBool(r.URL.Query().Get("is_vaccinated"))
	if err == nil {
		params.IsVaccinated = &IsVaccinated
	}
	HasVetPassport, err := strconv.ParseBool(r.URL.Query().Get("has_vet_passport"))
	if err == nil {
		params.HasVetPassport = &HasVetPassport
	}
	status := r.URL.Query().Get("status")
	if status != "" {
		params.Status = &status
	}
	shelterID, err := strconv.Atoi(r.URL.Query().Get("shelter_id"))
	if err == nil {
		params.ShelterID = &shelterID
	}
	return params
}
