package main

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
)

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, APIError{Error: message})
}

type Handler struct {
	store *Store
}

func NewHandler(store *Store) *Handler {
	return &Handler{store: store}
}

func (h *Handler) CreatePost(w http.ResponseWriter, r *http.Request) {
	var input CreatePostInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	if strings.TrimSpace(input.Title) == "" {
		writeError(w, http.StatusUnprocessableEntity, "Title harus diisi")
		return
	}

	post, err := h.store.CreatePost(r.Context(), input)
	if err != nil {
		slog.Error("CreatePost error", "err", err)
		writeError(w, http.StatusInternalServerError, "Gagal membuat post")
		return
	}

	writeJSON(w, http.StatusCreated, post)
}

func (h *Handler) GetPosts(w http.ResponseWriter, r *http.Request) {
	posts, err := h.store.GetPosts(r.Context())
	if err != nil {
		slog.Error("GetPosts error", "err", err)
		writeError(w, http.StatusInternalServerError, "Gagal mengambil posts")
		return
	}

	if posts == nil {
		posts = []Post{}
	}

	writeJSON(w, http.StatusOK, posts)
}

func (h *Handler) GetPost(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "ID harus angka")
		return
	}

	post, err := h.store.GetPostByID(r.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "tidak ditemukan") {
			writeError(w, http.StatusNotFound, "Post tidak ditemukan")
			return
		}
		slog.Error("GetPost error", "err", err)
		writeError(w, http.StatusInternalServerError, "Gagal mengambil post")
		return
	}

	writeJSON(w, http.StatusOK, post)
}

func (h *Handler) UpdatePost(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "ID harus angka")
		return
	}

	var input UpdatePostInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	post, err := h.store.UpdatePost(r.Context(), id, input)
	if err != nil {
		if strings.Contains(err.Error(), "tidak ditemukan") {
			writeError(w, http.StatusNotFound, "Post tidak ditemukan")
			return
		}
		slog.Error("UpdatePost error", "err", err)
		writeError(w, http.StatusInternalServerError, "Gagal update post")
		return
	}

	writeJSON(w, http.StatusOK, post)
}

func (h *Handler) DeletePost(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "ID harus angka")
		return
	}

	if err := h.store.DeletePost(r.Context(), id); err != nil {
		if strings.Contains(err.Error(), "tidak ditemukan") {
			writeError(w, http.StatusNotFound, "Post tidak ditemukan")
			return
		}
		slog.Error("DeletePost error", "err", err)
		writeError(w, http.StatusInternalServerError, "Gagal hapus post")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
