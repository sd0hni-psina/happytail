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
	svc       ShelterService
	animalSvc AnimalService
}

func NewShelterHandler(svc ShelterService, animalSvc AnimalService) *ShelterHandler {
	return &ShelterHandler{svc: svc, animalSvc: animalSvc}
}

// GetAllShelters godoc
// @Summary Получить все приюты
// @Tags shelters
// @Accept json
// @Produce json
// @Success 200 {array} models.Shelter
// @Router /shelters [get]
func (h *ShelterHandler) GetAllShelters(w http.ResponseWriter, r *http.Request) {
	params := models.ParsePagination(r)

	shelters, total, err := h.svc.GetAllShelters(r.Context(), params)
	if err != nil {
		http.Error(w, "Failed to fetch shelters", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.NewPaginatedResponse(shelters, total, params))
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
		http.Error(w, "Name and Address are required", http.StatusBadRequest)
		return
	}
	if err := input.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
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
	// нет верхней границы валидации на уровне хендлера, пользователь может передать любой радиус
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

// UpdateShelter godoc
// @Summary Обновить данные приюта
// @Tags shelters
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param id path int true "ID приюта"
// @Param input body models.UpdateShelterInput true "поля для обновления"
// @Success 200 {object} models.Shelter
// @Failure 400 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /shelters/{id} [patch]
func (h *ShelterHandler) UpdateShelter(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid shelter ID", http.StatusBadRequest)
		return
	}

	var input models.UpdateShelterInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	if err := input.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	shelter, err := h.svc.UpdateShelter(r.Context(), id, input)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			http.Error(w, "Shelter not found", http.StatusNotFound)
			return
		}
		if errors.Is(err, models.ErrConflict) {
			http.Error(w, "Shelter with this email already exists", http.StatusConflict)
			return
		}
		http.Error(w, "Failed to update shelter", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(shelter)
}

// GetShelterAnimals godoc
// @Summary Получить животных приюта
// @Tags shelters
// @Produce json
// @Param id path int true "ID приюта"
// @Param page query int false "Страница"
// @Param limit query int false "Лимит"
// @Param status query string false "Фильтр по статусу"
// @Param type query string false "Фильтр по типу"
// @Success 200 {object} models.PaginatedResponse[models.Animal]
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /shelters/{id}/animals [get]
func (h *ShelterHandler) GetShelterAnimals(w http.ResponseWriter, r *http.Request) {
	shelterID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid shelter ID", http.StatusBadRequest)
		return
	}

	if _, err := h.svc.GetShelterByID(r.Context(), shelterID); err != nil {
		if errors.Is(err, models.ErrNotFound) {
			http.Error(w, "Shelter not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to fetch shelter", http.StatusInternalServerError)
		return
	}

	pagination := models.ParsePagination(r)
	filter := models.ParseFilter(r)
	filter.ShelterID = &shelterID

	animals, total, err := h.animalSvc.GetAllAnimals(r.Context(), pagination, filter)
	if err != nil {
		http.Error(w, "Failed to fetch animals", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.NewPaginatedResponse(animals, total, pagination))
}
