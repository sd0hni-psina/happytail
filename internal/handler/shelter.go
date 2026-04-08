package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/sd0hni-psina/happytail/internal/models"
)

type ShelterHandler struct {
	svc ShelterService
}

func NewShelterHandler(svc ShelterService) *ShelterHandler {
	return &ShelterHandler{svc: svc}
}

func (h *ShelterHandler) GetAllShelters(w http.ResponseWriter, r *http.Request) {
	shelters, err := h.svc.GetAllShelters(r.Context())
	if err != nil {
		http.Error(w, "Failed to fetch shelters", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(shelters)
}

func (h *ShelterHandler) GetShelterByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid shelter ID", http.StatusBadRequest)
		return
	}
	shelter, err := h.svc.GetShelterByID(r.Context(), id)
	if err != nil {
		http.Error(w, "Failed to fetch shelter", http.StatusInternalServerError)
		return
	}
	if shelter == nil {
		http.Error(w, "Shelter not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(shelter)
}

func (h *ShelterHandler) CreateShelter(w http.ResponseWriter, r *http.Request) {
	var input models.CreateShelterInput

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if input.Name == "" || input.Address == "" {
		http.Error(w, "Name, Address and Email are required", http.StatusBadRequest)
		return
	}

	shelter, err := h.svc.CreateShelter(r.Context(), input)
	if err != nil {
		http.Error(w, "Failed to create shelter", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(shelter)
}
