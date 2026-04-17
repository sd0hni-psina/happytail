package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/sd0hni-psina/happytail/internal/middleware"
	"github.com/sd0hni-psina/happytail/internal/models"
)

type AdoptionHandler struct {
	svc AdoptionService
}

func NewAdoptionHandler(svc AdoptionService) *AdoptionHandler {
	return &AdoptionHandler{svc: svc}
}

// CreateAdoption godoc
// @Summary Создать заявку на усыновление
// @Tags adoptions
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param input body models.CreateAdoptionInput true "данные заявки"
// @Success 201 {object} models.Adoption
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /adoptions [post]
func (h *AdoptionHandler) CreateAdoption(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		http.Error(w, "Failed to get user ID", http.StatusUnauthorized)
		return
	}

	var input models.CreateAdoptionInput
	json.NewDecoder(r.Body).Decode(&input)
	if input.AnimalID == 0 {
		http.Error(w, "Animal ID is required", http.StatusBadRequest)
		return
	}

	adoption, err := h.svc.CreateAdoption(r.Context(), userID, input.AnimalID)
	if err != nil {
		if errors.Is(err, models.ErrNotAvailable) {
			http.Error(w, "Animal is not available for adoption", http.StatusConflict)
			return
		}
		http.Error(w, "Failed to create adoption", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(adoption)
}
