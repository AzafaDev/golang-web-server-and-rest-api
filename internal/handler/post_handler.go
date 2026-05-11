package handler

import (
	"belajar-backend-golang/internal/service"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator"
)

type PostHandler struct {
	service    *service.PostService
	validation *validator.Validate
}

func NewPostHandler(service *service.PostService) *PostHandler {
	return &PostHandler{
		service:    service,
		validation: validator.New(),
	}
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func errorJSON(w http.ResponseWriter, status int, data any) {
	writeJSON(w, status, data)
}

func (h *PostHandler) GetAllPosts(w http.ResponseWriter, r *http.Request) {
	posts, err := h.service.GetAllPosts(r.Context())
	if err != nil {
		errorJSON(w, http.StatusInternalServerError, "Failed to fetch posts")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"message": "Fetching posts successfully",
		"data":    posts,
	})
}

func (h *PostHandler) GetPostById(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		errorJSON(w, http.StatusBadRequest, "Id must be number")
		return
	}
	post, err := h.service.GetPostById(r.Context(), id)
	if err != nil {
		errorJSON(w, http.StatusInternalServerError, "Failed to fetch post")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"message": "Fetching post successfully",
		"data":    post,
	})
}

func (h *PostHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(int)
	if !ok {
		errorJSON(w, http.StatusUnauthorized, map[string]any{
			"message": "Failed to get user's identity",
		})
		return
	}
	fmt.Printf("User ID %d is now trying to make a post\n", userID)
	var input struct {
		Title   string `json:"title" validate:"required,min=5,max=30"`
		Content string `json:"content"`
	}
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		errorJSON(w, http.StatusBadRequest, "Invalid JSON format")
		return
	}
	err = h.validation.Struct(input)
	if err != nil {
		errorJSON(w, http.StatusBadRequest, map[string]any{
			"status": "error",
			"errors": FormatValidationError(err),
		})
		return
	}

	post, err := h.service.CreatePost(r.Context(), input.Title, input.Content)
	if err != nil {
		errorJSON(w, http.StatusInternalServerError, "Failed to create post")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"message": "Creating post successfully",
		"data":    post,
	})
}

func (h *PostHandler) UpdatePost(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(int)
	if !ok {
		errorJSON(w, http.StatusUnauthorized, map[string]any{
			"message": "Failed to get user's identity",
		})
		return
	}
	fmt.Printf("User ID %d is now trying to update a post\n", userID)
	var input struct {
		Title   string `json:"title" validate:"omitempty,min=5,max=30"`
		Content string `json:"content"`
	}
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		errorJSON(w, http.StatusBadRequest, "Invalid JSON format")
		return
	}
	err = h.validation.Struct(input)
	if err != nil {
		errorJSON(w, http.StatusBadRequest, map[string]any{
			"status": "error",
			"errors": FormatValidationError(err),
		})
		return
	}
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		errorJSON(w, http.StatusBadRequest, "Id must be number")
		return
	}
	post, err := h.service.UpdatePost(r.Context(), input.Title, input.Content, id)
	if err != nil {
		errorJSON(w, http.StatusInternalServerError, "Failed to update post")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"message": "Updating post successfully",
		"data":    post,
	})
}

func (h *PostHandler) DeletePost(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(int)
	if !ok {
		errorJSON(w, http.StatusUnauthorized, map[string]any{
			"message": "Failed to get user's identity",
		})
		return
	}
	fmt.Printf("User ID %d is now trying to delete a post\n", userID)
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		errorJSON(w, http.StatusBadRequest, "Id must be number")
		return
	}
	err = h.service.DeletePost(r.Context(), id)
	if err != nil {
		errorJSON(w, http.StatusInternalServerError, "Failed to delete post")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"message": "Deleting post successfully",
	})
}
