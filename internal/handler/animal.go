package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
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

func (h *AnimalHandler) GetAnimalByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid animal ID", http.StatusBadRequest)
		return
	}
	animal, err := h.repo.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, "Animal not found", http.StatusNotFound)
		return
	}
	if animal == nil {
		http.Error(w, "Animal not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(animal)
}
