package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/sd0hni-psina/happytail/internal/middleware"
	"github.com/sd0hni-psina/happytail/internal/models"
)

type PostHandler struct {
	svc PostService
}

func NewPostHandler(svc PostService) *PostHandler {
	return &PostHandler{svc: svc}
}

// GetAllPost godoc
// @Summary Получить все посты
// @Description Возвращает список всех постов
// @Tags Post
// @Produce json
// @Success 200 {array} models.Post
// @Failure 500 {string} string "Failed to fetch posts"
// @Router /posts [get]
func (h *PostHandler) GetAllPost(w http.ResponseWriter, r *http.Request) {
	params := models.ParsePagination(r)

	posts, total, err := h.svc.GetAllPost(r.Context(), params)
	if err != nil {
		http.Error(w, "Failed to fetch posts", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.NewPaginatedResponse(posts, total, params))
}

// GetPostByID godoc
// @Summary Получить пост по ID
// @Description Возвращает один пост по его ID
// @Tags Post
// @Produce json
// @Param id path int true "ID поста"
// @Success 200 {object} models.Post
// @Failure 400 {string} string "Invalid post ID"
// @Failure 404 {string} string "Post not found"
// @Failure 500 {string} string "Failed to fetch post"
// @Router /posts/{id} [get]
func (h *PostHandler) GetPostByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}
	post, err := h.svc.GetPostByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			http.Error(w, "Post not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to fetch post", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(post)
}

// CreatePost godoc
// @Summary Создать пост
// @Description Создаёт новый пост
// @Tags Post
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param input body models.CreatePostInput true "данные поста"
// @Success 201 {object} models.Post
// @Failure 400 {string} string "Invalid input"
// @Failure 409 {string} string "Conflict"
// @Failure 500 {string} string "Failed to create post"
// @Router /posts [post]
func (h *PostHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	var input models.CreatePostInput

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	input.UserID = userID

	if err := input.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	post, err := h.svc.CreatePost(r.Context(), input)
	if err != nil {
		if errors.Is(err, models.ErrConflict) {
			http.Error(w, "Something wrong with input", http.StatusConflict)
			return
		}
		http.Error(w, "Failed to create post", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(post)
}

// UpdatePostStatus godoc
// @Summary Обновить статус поста
// @Tags posts
// @Security ApiKeyAuth
// @Accept json
// @Param id path int true "ID поста"
// @Param input body models.UpdatePostStatusInput true "новый статус"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /posts/{id}/status [patch]
func (h *PostHandler) UpdatePostStatus(w http.ResponseWriter, r *http.Request) {
	postID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var input models.UpdatePostStatusInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	if err := input.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.svc.UpdateStatus(r.Context(), postID, userID, input.Status); err != nil {
		if errors.Is(err, models.ErrNotFound) {
			http.Error(w, "Post not found", http.StatusNotFound)
			return
		}
		if errors.Is(err, models.ErrForbidden) {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		http.Error(w, "Failed to update post status", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
