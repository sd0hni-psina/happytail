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
	filter := models.ParseFilter(r)

	animals, total, err := h.svc.GetAllAnimals(r.Context(), params, filter)
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

// Любой авторизованный пользователь может добавить животное (для приюта, которого он не является администратором). Администратор приюта может редактировать и удалять животных своего приюта.

// UpdateAnimal godoc
// @Summary Update animal
// @Description Update existing animal by ID
// @Tags animals
// @Accept json
// @Produce json
// @Param id path int true "Animal ID"
// @Param input body models.UpdateAnimalInput true "Update animal input"
// @Success 200 {object} models.Animal
// @Failure 400 {string} string "Invalid input or invalid ID"
// @Failure 404 {string} string "Animal not found"
// @Failure 500 {string} string "Internal server error"
// @Router /animals/{id} [put]
func (h *AnimalHandler) UpdateAnimal(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid animal ID", http.StatusBadRequest)
		return
	}

	var input models.UpdateAnimalInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	if err := input.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	animal, err := h.svc.UpdateAnimal(r.Context(), id, input)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			http.Error(w, "Animal not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to update animal", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(animal)
}

// DeleteAnimal godoc
// @Summary Удалить животное
// @Description Удаляет животное по ID (soft delete)
// @Tags animals
// @Param id path int true "ID животного"
// @Success 204 "Успешно удалено"
// @Failure 400 {string} string "invalid animal ID"
// @Failure 404 {string} string "Animal not found"
// @Failure 500 {string} string "Failed to delete animal"
// @Router /animals/{id} [delete]
func (h *AnimalHandler) DeleteAnimal(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid animal ID", http.StatusBadRequest)
		return
	}

	if err := h.svc.DeleteAnimal(r.Context(), id); err != nil {
		if errors.Is(err, models.ErrNotFound) {
			http.Error(w, "Animal not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to delete animal", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// ShareAnimal godoc
// @Summary Поделиться животным
// @Tags animals
// @Security ApiKeyAuth
// @Produce json
// @Param id path int true "ID животного"
// @Success 200 {object} models.Animal
// @Failure 404 {object} map[string]string
// @Router /animals/{id}/share [post]
func (h *AnimalHandler) ShareAnimal(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid animal ID", http.StatusBadRequest)
		return
	}
	animal, err := h.svc.ShareAnimal(r.Context(), id)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			http.Error(w, "Animal not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to share animal", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(animal)
}
