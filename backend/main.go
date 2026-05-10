package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

var db *pgxpool.Pool

func main() {
	// 1. Load .env
	if err := godotenv.Load(); err != nil {
		log.Fatal("There is no .env file!")
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL is not set!")
	}

	// 2. Koneksi ke Pool
	ctx := context.Background()
	var err error
	db, err = pgxpool.New(ctx, dbURL) // Menggunakan '=' untuk mengisi variabel global [cite: 109]
	if err != nil {
		log.Fatal("Failed to connect to Neon:", err)
	}
	defer db.Close()

	// Cek koneksi
	if err := db.Ping(ctx); err != nil {
		log.Fatal("Database unreachable:", err)
	}
	fmt.Println("✅ Connected to Neon PostgreSQL (via Pool)")

	createTable(ctx)

	// 3. Routes
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/posts", func(r chi.Router) {
		r.Get("/", getPosts)
		r.Post("/", createPost)
		r.Put("/{id}", updatePost)
		r.Delete("/{id}", deletePost)
	})

	r.Route("/hello", func(r chi.Router) {
		r.Get("/", helloHandler)
	})

	fmt.Println("🚀 Server is running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// Handler Hello dengan format JSON yang benar
func helloHandler(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"message": "Hello World!"})
}

func createTable(ctx context.Context) {
	query := `
	CREATE TABLE IF NOT EXISTS posts (
		id SERIAL PRIMARY KEY,
		title VARCHAR(255) NOT NULL,
		content TEXT NOT NULL DEFAULT '', -- Petik satu untuk Postgres [cite: 131]
		created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);`
	if _, err := db.Exec(ctx, query); err != nil {
		log.Fatal("Failed to make a table:", err)
	}
	fmt.Println("📦 Table posts is ready")
}

func deletePost(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	stringId := r.URL.Query().Get("id")

	postId, err := strconv.Atoi(stringId)
	if err != nil {
		http.Error(w, "Id must be a number", http.StatusBadRequest)
		return // FIXED: Sekarang fungsi akan berhenti di sini jika error
	}

	// Menggunakan db.Exec untuk efisiensi
	query := `DELETE FROM posts WHERE id = $1`
	result, err := db.Exec(ctx, query, postId)

	if err != nil {
		http.Error(w, "Failed to delete post", http.StatusInternalServerError)
		return
	}

	// Cek apakah ada baris yang benar-benar dihapus
	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "Deleted post successfully"})
}

func updatePost(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context() // Gunakan context dari request
	idStr := r.URL.Query().Get("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID harus berupa angka", http.StatusBadRequest)
		return
	}

	// Cari data lama
	existing, err := getPostById(ctx, id)
	if err != nil {
		if err == pgx.ErrNoRows {
			http.Error(w, "Post tidak ditemukan", http.StatusNotFound)
		} else {
			http.Error(w, "Gagal mengambil data", http.StatusInternalServerError)
		}
		return
	}

	// Decode input partial
	var input struct {
		Title   *string `json:"title"`
		Content *string `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Format JSON tidak valid", http.StatusBadRequest)
		return
	}

	// Update field jika dikirim
	if input.Title != nil {
		existing.Title = *input.Title
	}
	if input.Content != nil {
		existing.Content = *input.Content
	}

	var updated Post
	query := `UPDATE posts SET title=$1, content=$2, updated_at=NOW() WHERE id=$3 
	          RETURNING id, title, content, created_at, updated_at`

	err = db.QueryRow(ctx, query, existing.Title, existing.Content, existing.Id).
		Scan(&updated.Id, &updated.Title, &updated.Content, &updated.CreatedAt, &updated.UpdatedAt)

	if err != nil {
		http.Error(w, "Gagal mengupdate database", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, updated)
}

func createPost(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if input.Title == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	}

	var p Post
	query := `INSERT INTO posts (title, content) VALUES ($1, $2) 
	          RETURNING id, title, content, created_at, updated_at`

	err := db.QueryRow(r.Context(), query, input.Title, input.Content).
		Scan(&p.Id, &p.Title, &p.Content, &p.CreatedAt, &p.UpdatedAt)

	if err != nil {
		http.Error(w, "Failed to create post", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusCreated, p)
}

func getPosts(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query(r.Context(), "SELECT id, title, content, created_at, updated_at FROM posts ORDER BY created_at DESC")
	if err != nil {
		http.Error(w, "Failed to fetch data", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var p Post
		if err := rows.Scan(&p.Id, &p.Title, &p.Content, &p.CreatedAt, &p.UpdatedAt); err != nil {
			http.Error(w, "Scan error", http.StatusInternalServerError)
			return
		}
		posts = append(posts, p)
	}

	writeJSON(w, http.StatusOK, posts)
}

func getPostById(ctx context.Context, id int) (Post, error) {
	var p Post
	query := `SELECT id, title, content, created_at, updated_at FROM posts WHERE id=$1`
	err := db.QueryRow(ctx, query, id).Scan(&p.Id, &p.Title, &p.Content, &p.CreatedAt, &p.UpdatedAt)
	return p, err
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
