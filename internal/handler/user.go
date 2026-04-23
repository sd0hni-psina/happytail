package handler

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/sd0hni-psina/happytail/internal/middleware"
	"github.com/sd0hni-psina/happytail/internal/models"
)

type UserHandler struct {
	svc UserService
}

func NewUserHandler(svc UserService) *UserHandler {
	return &UserHandler{svc: svc}
}

// GetAllUsers godoc
// @Summary Получить всех пользователей
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {array} models.User
// @Router /users [get]
func (h *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	user, err := h.svc.GetAllUsers(r.Context())
	if err != nil {
		http.Error(w, "Failed to fetch users", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// GetUserByID godoc
// @Summary Получить пользователя по ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "ID пользователя"
// @Success 200 {object} models.User
// @Failure 404 {object} map[string]string
// @Router /users/{id} [get]
func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	user, err := h.svc.GetUserByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to fetch user", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// GetMe godoc
// @Summary Получить свой профиль
// @Tags users
// @Security ApiKeyAuth
// @Produce json
// @Success 200 {object} models.UserPublic
// @Failure 401 {object} map[string]string
// @Router /users/me [get]
func (h *UserHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	user, err := h.svc.GetUserByID(r.Context(), userID)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to fetch user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// CreateUser godoc
// @Summary Создать пользователя
// @Tags users
// @Accept json
// @Produce json
// @Param input body models.CreateUserInput true "данные пользователя"
// @Success 201 {object} models.User
// @Failure 400 {object} map[string]string
// @Router /users [post]
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var input models.CreateUserInput

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if err := input.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := h.svc.CreateUser(r.Context(), input)
	if err != nil {
		if errors.Is(err, models.ErrConflict) {
			http.Error(w, "Email already exists", http.StatusConflict)
			return
		}
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

// Login godoc
// @Summary Логин пользователя
// @Tags auth
// @Accept json
// @Produce json
// @Param input body models.LoginInput true "логин и пароль"
// @Success 200 {object} models.AuthResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /auth/login [post]
func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var input models.LoginInput

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if input.Email == "" || input.Password == "" {
		http.Error(w, "Email and Password are required", http.StatusBadRequest)
		return
	}

	auth, err := h.svc.Login(r.Context(), input.Email, input.Password)
	if err != nil {
		slog.Warn("login failed", "error", err)
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(auth)
}

// Refresh godoc
// @Summary Обновить access token
// @Tags auth
// @Accept json
// @Produce json
// @Param input body models.RefreshInput true "refresh token"
// @Success 200 {object} models.AuthResponse
// @Failure 401 {object} map[string]string
// @Router /auth/refresh [post]
func (h *UserHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var input models.RefreshInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if input.RefreshToken == "" {
		http.Error(w, "Refresh token is required", http.StatusBadRequest)
		return
	}

	response, err := h.svc.Refresh(r.Context(), input.RefreshToken)
	if err != nil {
		http.Error(w, "Invalid or expired refresh token", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Logout godoc
// @Summary Выход из системы
// @Tags auth
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param input body models.RefreshInput true "refresh token"
// @Success 204
// @Router /auth/logout [post]
func (h *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
	var input models.RefreshInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.svc.Logout(r.Context(), input.RefreshToken); err != nil {
		http.Error(w, "Failed to logout", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
