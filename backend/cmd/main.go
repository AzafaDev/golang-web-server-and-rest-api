package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/AzafaDev/golang-web-server-and-rest-api.git/internal/handler"
	"github.com/AzafaDev/golang-web-server-and-rest-api.git/internal/repository"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(".env file not found:", err)
	}
	dbUrl := os.Getenv("DATABASE_URL")
	config, err := pgxpool.ParseConfig(dbUrl)
	if err != nil {
		log.Fatal("Gagal parsing config:", err)
	}
	db, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Fatal("Failed to connect DB:", err)
	}
	defer db.Close()

	postRepo := repository.NewPostRepository(db)
	postHandler := handler.NewPostHandler(postRepo)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/posts", func(r chi.Router) {
		r.Get("/", postHandler.GetAll)
		r.Post("/", postHandler.Post)
		r.Get("/{id}", postHandler.GetByID)
		r.Put("/{id}", postHandler.Update)
		r.Delete("/{id}", postHandler.Delete)
	})

	fmt.Println("🚀 Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
