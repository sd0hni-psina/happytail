package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/sd0hni-psina/happytail/internal/models"
)

type AnimalPhotoHandler struct {
	svc AnimalPhotoService
}

func NewAnimalPhotoHandler(svc AnimalPhotoService) *AnimalPhotoHandler {
	return &AnimalPhotoHandler{svc: svc}
}

func (h *AnimalPhotoHandler) AddPhoto(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("animal_id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	var input models.AnimalPhotoInput

	err = json.NewDecoder(r.Body).Decode(&input)
	input.AnimalID = id

	photo, err := h.svc.AddPhoto(r.Context(), input)
	if err != nil {
		http.Error(w, "Failed to add photo", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(photo)
}

func (h *AnimalPhotoHandler) DeletePhoto(w http.ResponseWriter, r *http.Request) {}

func (h *AnimalPhotoHandler) MakeMainPhoto(w http.ResponseWriter, r *http.Request) {}

func (h *AnimalPhotoHandler) GetAllPhotos(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid id", http.StatusBadRequest)
		return
	}
	photos, err := h.svc.GetAllPhotos(r.Context(), id)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to fetch photos", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(photos)

}
