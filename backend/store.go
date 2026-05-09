package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Store struct {
	db *pgxpool.Pool
}

func NewStore(ctx context.Context) (*Store, error) {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return nil, fmt.Errorf("DATABASE_URL tidak diset")
	}

	config, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		return nil, fmt.Errorf("gagal parse config: %w", err)
	}

	config.MaxConns = 10
	config.MinConns = 2

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("gagal membuat connection pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("gagal ping database: %w", err)
	}

	if err := createTable(ctx, pool); err != nil {
		pool.Close()
		return nil, err
	}

	return &Store{db: pool}, nil
}

func (s *Store) Close() {
	s.db.Close()
}

func createTable(ctx context.Context, pool *pgxpool.Pool) error {
	query := `
	CREATE TABLE IF NOT EXISTS posts (
		id SERIAL PRIMARY KEY,
		title TEXT NOT NULL,
		content TEXT NOT NULL DEFAULT '',
		created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);`

	_, err := pool.Exec(ctx, query)
	if err != nil {
		return fmt.Errorf("gagal membuat tabel posts: %w", err)
	}
	return nil
}

func (s *Store) CreatePost(ctx context.Context, input CreatePostInput) (Post, error) {
	var p Post
	query := `
		INSERT INTO posts (title, content)
		VALUES ($1, $2)
		RETURNING id, title, content, created_at, updated_at
	`
	row := s.db.QueryRow(ctx, query, input.Title, input.Content)
	err := row.Scan(&p.ID, &p.Title, &p.Content, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return Post{}, fmt.Errorf("gagal insert post: %w", err)
	}
	return p, nil
}

func (s *Store) GetPosts(ctx context.Context) ([]Post, error) {
	query := `SELECT id, title, content, created_at, updated_at FROM posts ORDER BY created_at DESC`
	rows, err := s.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("gagal query posts: %w", err)
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var p Post
		if err := rows.Scan(&p.ID, &p.Title, &p.Content, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, fmt.Errorf("gagal scan post: %w", err)
		}
		posts = append(posts, p)
	}
	return posts, rows.Err()
}

func (s *Store) GetPostByID(ctx context.Context, id int) (Post, error) {
	var p Post
	query := `SELECT id, title, content, created_at, updated_at FROM posts WHERE id = $1`
	row := s.db.QueryRow(ctx, query, id)
	err := row.Scan(&p.ID, &p.Title, &p.Content, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return Post{}, fmt.Errorf("post tidak ditemukan")
		}
		return Post{}, fmt.Errorf("gagal query post: %w", err)
	}
	return p, nil
}

func (s *Store) UpdatePost(ctx context.Context, id int, input UpdatePostInput) (Post, error) {
	existing, err := s.GetPostByID(ctx, id)
	if err != nil {
		return Post{}, err
	}

	if input.Title != nil {
		existing.Title = *input.Title
	}
	if input.Content != nil {
		existing.Content = *input.Content
	}

	query := `
		UPDATE posts
		SET title = $1, content = $2, updated_at = NOW()
		WHERE id = $3
		RETURNING id, title, content, created_at, updated_at
	`
	row := s.db.QueryRow(ctx, query, existing.Title, existing.Content, id)
	var updated Post
	err = row.Scan(&updated.ID, &updated.Title, &updated.Content, &updated.CreatedAt, &updated.UpdatedAt)
	if err != nil {
		return Post{}, fmt.Errorf("gagal update post: %w", err)
	}
	return updated, nil
}

func (s *Store) DeletePost(ctx context.Context, id int) error {
	query := `DELETE FROM posts WHERE id = $1`
	tag, err := s.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("gagal hapus post: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("post tidak ditemukan")
	}
	return nil
}
