package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/sd0hni-psina/happytail/internal/models"
)

type RoleHandler struct {
	svc RoleService
}

func NewRoleHandler(svc RoleService) *RoleHandler {
	return &RoleHandler{svc: svc}
}

// AppointRole godoc
// @Summary Назначить роль пользователю
// @Description Назначает роль пользователю (например admin, moderator)
// @Tags Roles
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param input body models.RoleInput true "данные для назначения роли"
// @Success 201 {object} models.Role
// @Failure 400 {string} string "Invalid input"
// @Failure 500 {string} string "Cannot appoint role"
// @Router /roles [post]
func (h *RoleHandler) AppointRole(w http.ResponseWriter, r *http.Request) {
	var input models.RoleInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	if err := input.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	info, err := h.svc.AppointRole(r.Context(), input)
	if err != nil {
		http.Error(w, "Cannot appoint", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(info)
}

// RemoveRole godoc
// @Summary Удалить роль
// @Description Удаляет роль по ID
// @Tags Roles
// @Security ApiKeyAuth
// @Param id path int true "ID роли"
// @Success 204 "No Content"
// @Failure 400 {string} string "Invalid id"
// @Failure 500 {string} string "Failed to remove role"
// @Router /roles/{id} [delete]
func (h *RoleHandler) RemoveRole(w http.ResponseWriter, r *http.Request) {
	roleIDStr := r.PathValue("id")
	roleID, err := strconv.Atoi(roleIDStr)
	if err != nil {
		http.Error(w, "Invalid id", http.StatusBadRequest)
		return
	}
	err = h.svc.RemoveRole(r.Context(), roleID)
	if err != nil {
		http.Error(w, "Failed to remove role", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
