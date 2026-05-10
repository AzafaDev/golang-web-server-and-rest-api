package repository

import (
	"context"

	"github.com/AzafaDev/golang-web-server-and-rest-api.git/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostRepository struct {
	db *pgxpool.Pool
}

func NewPostRepository(db *pgxpool.Pool) *PostRepository {
	return &PostRepository{db: db}
}

func (r *PostRepository) GetAll(ctx context.Context) ([]models.Post, error) {
	rows, err := r.db.Query(ctx, "SELECT id, title, content, created_at, updated_at FROM posts ORDER BY created_at DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var p models.Post
		if err := rows.Scan(&p.Id, &p.Title, &p.Content, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	return posts, nil
}

func (r *PostRepository) GetByID(ctx context.Context, id int) (models.Post, error) {
	var p models.Post
	query := `SELECT id, title, content, created_at, updated_at FROM posts WHERE id=$1`
	err := r.db.QueryRow(ctx, query, id).Scan(&p.Id, &p.Title, &p.Content, &p.CreatedAt, &p.UpdatedAt)
	return p, err
}

func (r *PostRepository) Create(ctx context.Context, title, content string) (models.Post, error) {
	var p models.Post
	query := `INSERT INTO posts (title, content) VALUES ($1, $2) RETURNING id, title, content, created_at, updated_at`
	err := r.db.QueryRow(ctx, query, title, content).Scan(&p.Id, &p.Title, &p.Content, &p.CreatedAt, &p.UpdatedAt)
	return p, err
}

func (r *PostRepository) Update(ctx context.Context, p models.Post) (models.Post, error) {
	var updated models.Post
	query := `UPDATE posts SET title=$1, content=$2, updated_at=NOW() WHERE id=$3 RETURNING id, title, content, created_at, updated_at`
	err := r.db.QueryRow(ctx, query, p.Title, p.Content, p.Id).Scan(&updated.Id, &updated.Title, &updated.Content, &updated.CreatedAt, &updated.UpdatedAt)
	return updated, err
}

func (r *PostRepository) Delete(ctx context.Context, id int) error {
	_, err := r.db.Exec(ctx, "DELETE FROM posts WHERE id = $1", id)
	return err
}
