package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/sd0hni-psina/happytail/internal/models"
)

type AnimalHandler struct {
	svc AnimalService
}

func NewAnimalHandler(svc AnimalService) *AnimalHandler {
	return &AnimalHandler{svc: svc}
}

func (h *AnimalHandler) GetAllAnimals(w http.ResponseWriter, r *http.Request) {
	params := models.ParsePagination(r)

	animals, total, err := h.svc.GetAllAnimals(r.Context(), params)
	if err != nil {
		http.Error(w, "Failed to fetch animals", http.StatusInternalServerError)
		return
	}
	response := models.NewPaginatedResponse(animals, total, params)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *AnimalHandler) GetAnimalByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid animal ID", http.StatusBadRequest)
		return
	}
	animal, err := h.svc.GetAnimalByID(r.Context(), id)
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

func (h *AnimalHandler) CreateAnimal(w http.ResponseWriter, r *http.Request) {
	var input models.CreateAnimalInput

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if input.Name == "" || input.Type == "" {
		http.Error(w, "Name and Type are required", http.StatusBadRequest)
		return
	}

	animal, err := h.svc.CreateAnimal(r.Context(), input)
	if err != nil {
		http.Error(w, "Failed to create animal", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(animal)
}
