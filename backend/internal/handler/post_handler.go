package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/AzafaDev/golang-web-server-and-rest-api.git/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5"
)

type PostHandler struct {
	service  service.PostService
	validate *validator.Validate
}

func NewPostHandler(svc *service.PostService) *PostHandler {
	return &PostHandler{
		service:  *svc,
		validate: validator.New(),
	}
}

func (h *PostHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	posts, err := h.service.GetAll(r.Context())
	if err != nil {
		http.Error(w, "Failed to fetch posts", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"data": posts})
}

func (h *PostHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	post, err := h.service.GetPostByID(r.Context(), id)
	if err != nil {
		if err == pgx.ErrNoRows {
			http.Error(w, "Post not found", http.StatusNotFound)
		} else {
			http.Error(w, "Error fetching post", http.StatusInternalServerError)
		}
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"data": post})
}

func (h *PostHandler) Post(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title   string `json:"title" validate:"required,min=5,max=100"`
		Content string `json:"content" validate:"required,min=10"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}
	if err := h.validate.Struct(input); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{
			"error":  "Invalid JSON format",
			"detail": err.Error(),
		})
		return
	}
	newPost, err := h.service.CreatePost(r.Context(), input.Title, input.Content)
	if err != nil {
		http.Error(w, "Error in creating post", http.StatusInternalServerError)
		return
	}
	writeJSON(w, 201, map[string]any{"data": newPost})
}

func (h *PostHandler) Update(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title   *string `json:"title" validate:"omitempty,min=5"`
		Content *string `json:"content" validate:"omitempty,min=10"`
	}
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	p, err := h.service.GetPostByID(r.Context(), id)
	if err != nil {
		if err == pgx.ErrNoRows {
			http.Error(w, "Post not found", http.StatusNotFound)
		} else {
			http.Error(w, "Error fetching post", http.StatusInternalServerError)
		}
		return
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}
	if err := h.validate.Struct(input); err != nil {

	}
	if input.Title != nil {
		p.Title = *input.Title
	}
	if input.Content != nil {
		p.Content = *input.Content
	}
	updatedPost, err := h.service.UpdatePost(r.Context(), p)
	if err != nil {
		http.Error(w, "Error in updating post", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"data": updatedPost})
}

func (h *PostHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	_, err := h.service.GetPostByID(r.Context(), id)
	if err != nil {
		if err == pgx.ErrNoRows {
			http.Error(w, "Post not found", http.StatusNotFound)
		} else {
			http.Error(w, "Error fetching post", http.StatusInternalServerError)
		}
		return
	}
	if err := h.service.DeletePost(r.Context(), id); err != nil {
		http.Error(w, "Error in deleting post", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"message": "Deleted post successfully"})
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
