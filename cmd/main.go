package main

import (
	"belajar-backend-golang/internal/config"
	"belajar-backend-golang/internal/handler"
	"belajar-backend-golang/internal/middlewares"
	"belajar-backend-golang/internal/repository"
	"belajar-backend-golang/internal/service"
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	cfg := config.LoadEnv()
	dbConfig, err := pgxpool.ParseConfig(cfg.DBURL)
	if err != nil {
		log.Fatal("Failed to parse config:", err)
	}
	db, err := pgxpool.NewWithConfig(context.Background(), dbConfig)
	if err != nil {
		log.Fatal("Failed to connect DB:", err)
	}
	defer db.Close()
	runMigrations(cfg.DBURL)

	postRepo := repository.NewPostRepository(db)
	postService := service.NewPostService(postRepo)
	postHandler := handler.NewPostHandler(postService)

	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo, cfg)
	userHandler := handler.NewUserHandler(userService)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/post", func(r chi.Router) {
		r.Get("/", postHandler.GetAllPosts)
		r.Get("/{id}", postHandler.GetPostById)
		r.Group(func(r chi.Router) {
			r.Use(middlewares.AuthMiddleware(cfg.JWTSECRET))
			r.Post("/", postHandler.CreatePost)
			r.Put("/{id}", postHandler.UpdatePost)
			r.Delete("/{id}", postHandler.DeletePost)
		})
	})

	r.Route("/user", func(r chi.Router) {
		r.Post("/register", userHandler.Register)
		r.Post("/login", userHandler.Login)
	})

	addr := ":" + cfg.PORT
	fmt.Println("Server is running on port:", cfg.PORT)
	log.Fatal(http.ListenAndServe(addr, r))
}

func runMigrations(dbUrl string) {
	m, err := migrate.New("file://internal/migrations", dbUrl)
	if err != nil {
		log.Fatal("Failed to init migration:", err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal("Failed to run migration:", err)
	}
	fmt.Println("Database migration running successfully")
}
