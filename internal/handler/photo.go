package handler

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/sd0hni-psina/happytail/internal/models"
)

type AnimalPhotoHandler struct {
	svc AnimalPhotoService
}

func NewAnimalPhotoHandler(svc AnimalPhotoService) *AnimalPhotoHandler {
	return &AnimalPhotoHandler{svc: svc}
}

// AddPhoto godoc
// @Summary Добавить фото животного
// @Description Добавляет новое фото для конкретного животного
// @Tags AnimalPhotos
// @Accept json
// @Produce json
// @Param id path int true "ID животного"
// @Param input body models.AnimalPhotoInput true "данные фото"
// @Success 201 {object} models.AnimalPhoto
// @Failure 400 {string} string "Invalid input"
// @Failure 500 {string} string "Failed to add photo"
// @Router /animals/{id}/photos [post]
func (h *AnimalPhotoHandler) AddPhoto(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	if err = r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "filte too large", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("photo")
	if err != nil {
		http.Error(w, "photo field is required", http.StatusBadRequest)
		return
	}
	defer file.Close()

	ext := filepath.Ext(header.Filename)
	allowed := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".webp": true}
	if !allowed[ext] {
		http.Error(w, "only jpg, png, webp allowed", http.StatusBadRequest)
		return
	}

	isMain := r.FormValue("is_main") == "true"

	photo, err := h.svc.AddPhoto(r.Context(), id, file, header, isMain)
	if err != nil {
		slog.Error("failed to add photo", "error", err)
		http.Error(w, "Failed to add photo", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(photo)
}

// DeletePhoto godoc
// @Summary Удалить фото
// @Description Удаляет фото по ID
// @Tags AnimalPhotos
// @Param photo_id path int true "ID фото"
// @Success 204 "No Content"
// @Failure 400 {string} string "Invalid id"
// @Failure 500 {string} string "Failed to delete photo"
// @Router /photos/{photo_id} [delete]
func (h *AnimalPhotoHandler) DeletePhoto(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("photo_id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid id", http.StatusBadRequest)
		return
	}
	err = h.svc.DeletePhoto(r.Context(), id)
	if err != nil {
		http.Error(w, "Failed to delete photo", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// MakeMainPhoto godoc
// @Summary Сделать фото главным
// @Description Устанавливает фото как главное для животного
// @Tags AnimalPhotos
// @Param id path int true "ID животного"
// @Param photo_id path int true "ID фото"
// @Success 204 "No Content"
// @Failure 400 {string} string "Invalid id"
// @Failure 404 {string} string "Photo not available"
// @Failure 500 {string} string "Internal error"
// @Router /animals/{id}/photos/{photo_id}/main [patch]
func (h *AnimalPhotoHandler) MakeMainPhoto(w http.ResponseWriter, r *http.Request) {
	animalIDStr := r.PathValue("id")
	animalID, err := strconv.Atoi(animalIDStr)
	if err != nil {
		http.Error(w, "Invalid animal id", http.StatusBadRequest)
		return
	}
	photoIDStr := r.PathValue("photo_id")
	photoID, err := strconv.Atoi(photoIDStr)
	if err != nil {
		http.Error(w, "Invalid photo id", http.StatusBadRequest)
		return
	}
	err = h.svc.MakeMainPhoto(r.Context(), animalID, photoID)
	if err != nil {
		if errors.Is(err, models.ErrNotAvailable) {
			http.Error(w, "Photo not available", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to update photo", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// GetAllPhotos godoc
// @Summary Получить все фото животного
// @Description Возвращает список фото для конкретного животного (главное фото первым)
// @Tags AnimalPhotos
// @Produce json
// @Param id path int true "ID животного"
// @Success 200 {array} models.AnimalPhoto
// @Failure 400 {string} string "Invalid id"
// @Failure 404 {string} string "Not found"
// @Failure 500 {string} string "Failed to fetch photos"
// @Router /animals/{id}/photos [get]
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
