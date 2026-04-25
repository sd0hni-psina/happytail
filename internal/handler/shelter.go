package handler

import (
	"encoding/json"
	"errors"
	"log/slog"
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

// GetAllShelters godoc
// @Summary Получить все приюты
// @Tags shelters
// @Accept json
// @Produce json
// @Success 200 {array} models.Shelter
// @Router /shelters [get]
func (h *ShelterHandler) GetAllShelters(w http.ResponseWriter, r *http.Request) {
	shelters, err := h.svc.GetAllShelters(r.Context())
	if err != nil {
		http.Error(w, "Failed to fetch shelters", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(shelters)
}

// GetShelterByID godoc
// @Summary Получить приют по ID
// @Tags shelters
// @Accept json
// @Produce json
// @Param id path int true "ID приюта"
// @Success 200 {object} models.Shelter
// @Failure 404 {object} map[string]string
// @Router /shelters/{id} [get]
func (h *ShelterHandler) GetShelterByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid shelter ID", http.StatusBadRequest)
		return
	}
	shelter, err := h.svc.GetShelterByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			http.Error(w, "Shelter not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to fetch shelter", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(shelter)
}

// CreateShelter godoc
// @Summary Создать приют
// @Tags shelters
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param input body models.CreateShelterInput true "данные приюта"
// @Success 201 {object} models.Shelter
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /shelters [post]
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
		slog.Error("failed to create shelter", "error", err)
		if errors.Is(err, models.ErrConflict) {
			http.Error(w, "Shelter with the same name or email already exists", http.StatusConflict)
			return
		}
		http.Error(w, "Failed to create shelter", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(shelter)
}

// FindNearby godoc
// @Summary Найти ближайшие приюты
// @Tags shelters
// @Produce json
// @Param lat query number true "Широта"
// @Param lon query number true "Долгота"
// @Param radius query number false "Радиус поиска в км (default: 10, max: 500)"
// @Success 200 {array} models.ShelterWithDistance
// @Router /shelters/nearby [get]
func (h *ShelterHandler) FindNearby(w http.ResponseWriter, r *http.Request) {
	latStr := r.URL.Query().Get("lat")
	lonStr := r.URL.Query().Get("lon")

	if latStr == "" || lonStr == "" {
		http.Error(w, "lat and lon are required", http.StatusBadRequest)
		return
	}

	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		http.Error(w, "invalid lat", http.StatusBadRequest)
		return
	}

	lon, err := strconv.ParseFloat(lonStr, 64)
	if err != nil {
		http.Error(w, "invalid lon", http.StatusBadRequest)
		return
	}

	var radius float64
	if radiusStr := r.URL.Query().Get("radius"); radiusStr != "" {
		radius, err = strconv.ParseFloat(radiusStr, 64)
		if err != nil {
			http.Error(w, "invalid radius", http.StatusBadRequest)
			return
		}
	}

	params := models.NearbyParams{
		Latitude:  lat,
		Longitude: lon,
		RadiusKm:  radius,
	}

	shelters, err := h.svc.FindNearby(r.Context(), params)
	if err != nil {
		slog.Error("failed to find nearby shelters", "error", err)
		http.Error(w, "Failed to find nearby shelters", http.StatusInternalServerError)
		return
	}
	if shelters == nil {
		shelters = []models.ShelterWithDistance{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(shelters)
}
