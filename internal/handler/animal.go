package handler

import (
	"encoding/json"
	"errors"
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

// GetAllAnimals godoc
// @Summary Получить всех животных
// @Tags animals
// @Produce json
// @Success 200 {array} models.Animal
// @Router /animals [get]
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
// GetAnimalByID godoc
// @Summary Получить животное по ID
// @Tags animals
// @Accept json
// @Produce json
// @Param id path int true "ID животного"
// @Success 200 {object} models.Animal
// @Failure 404 {object} map[string]string
// @Router /animals/{id} [get]
func (h *AnimalHandler) GetAnimalByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid animal ID", http.StatusBadRequest)
		return
	}
	animal, err := h.svc.GetAnimalByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			http.Error(w, "Animal not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to fetch animal", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(animal)
}

// CreateAnimal godoc
// @Summary Создать животное
// @Tags animals
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param input body models.CreateAnimalInput true "данные животного"
// @Success 201 {object} models.Animal
// @Router /animals [post]
func (h *AnimalHandler) CreateAnimal(w http.ResponseWriter, r *http.Request) {
	var input models.CreateAnimalInput

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if err := input.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
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
