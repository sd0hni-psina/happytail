package models

type CreateAnimalInput struct {
	Type           string  `json:"animal_type"`
	Name           string  `json:"name"`
	Age            *int    `json:"age"`
	Breed          *string `json:"breed"`
	Color          *string `json:"color"`
	IsVaccinated   bool    `json:"is_vaccinated"`
	HasVetPassport bool    `json:"has_vet_passport" `

	Description *string `json:"description"`
	ShelterID   *int    `json:"shelter_id"`
}
