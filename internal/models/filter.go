package models

import (
	"net/http"
	"strconv"
)

type FilterParams struct {
	Type           *string
	Breed          *string
	Color          *string
	IsVaccinated   *bool
	HasVetPassport *bool
	Status         *string
}

func ParseFilter(r *http.Request) FilterParams {
	params := FilterParams{}
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
	IsVaccinated, err := strconv.ParseBool(r.URL.Query().Get("IsVaccinated"))
	if err == nil {
		params.IsVaccinated = &IsVaccinated
	}
	HasVetPasswort, err := strconv.ParseBool(r.URL.Query().Get("HasVetPasswort"))
	if err == nil {
		params.HasVetPassport = &HasVetPasswort
	}
	status := r.URL.Query().Get("status")
	if status != "" {
		params.Status = &status
	}
	return params
}
