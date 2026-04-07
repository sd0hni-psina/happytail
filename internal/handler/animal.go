package handler

import (
	"encoding/json"
	"net/http"
)

type AnimalHandler struct {
	repo AnimalRepository
}

func NewAnimalHandler(repo AnimalRepository) *AnimalHandler {
	return &AnimalHandler{repo: repo}
}

func (h *AnimalHandler) GetAllAnimals(w http.ResponseWriter, r *http.Request) {
	animals, err := h.repo.GetAll(r.Context())
	if err != nil {
		http.Error(w, "Failed to fetch animals", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(animals)
}
