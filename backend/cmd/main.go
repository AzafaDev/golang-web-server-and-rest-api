package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/AzafaDev/golang-web-server-and-rest-api.git/internal/config"
	"github.com/AzafaDev/golang-web-server-and-rest-api.git/internal/handler"
	"github.com/AzafaDev/golang-web-server-and-rest-api.git/internal/repository"
	"github.com/AzafaDev/golang-web-server-and-rest-api.git/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	cfg := config.LoadConfig()
	dbConfig, err := pgxpool.ParseConfig(cfg.DBURL)
	if err != nil {
		log.Fatal("Failed to parsing config:", err)
	}
	db, err := pgxpool.NewWithConfig(context.Background(), dbConfig)
	if err != nil {
		log.Fatal("Failed to connect DB", err)
	}
	defer db.Close()

	postRepo := repository.NewPostRepository(db)
	postService := service.NewPostService(postRepo)
	postHandler := handler.NewPostHandler(postService)

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

	addr := ":" + cfg.Port
	fmt.Printf("🚀 Server running on %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, r))
}
