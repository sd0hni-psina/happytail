package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/sd0hni-psina/happytail/internal/models"
)

type PostHandler struct {
	svc PostService
}

func NewPostHandler(svc PostService) *PostHandler {
	return &PostHandler{svc: svc}
}

func (h *PostHandler) GetAllPost(w http.ResponseWriter, r *http.Request) {
	posts, err := h.svc.GetAllPost(r.Context())
	if err != nil {
		http.Error(w, "Failed to fetch posts", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}

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

func (h *PostHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	var input models.PostInput

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
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
