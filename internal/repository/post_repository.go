package repository

import (
	"belajar-backend-golang/internal/model"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostRepository struct {
	db *pgxpool.Pool
}

func NewPostRepository(db *pgxpool.Pool) *PostRepository {
	return &PostRepository{db: db}
}

func (r *PostRepository) GetAllPosts(ctx context.Context) ([]model.Post, error) {
	query := `SELECT id, title, content, created_at, updated_at FROM posts ORDER BY created_at DESC`
	var posts []model.Post
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var post model.Post
		err := rows.Scan(&post.Id, &post.Title, &post.Content, &post.CreatedAt, &post.UpdatedAt)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func (r *PostRepository) GetPostById(ctx context.Context, id int) (model.Post, error) {
	query := `SELECT id, title, content, created_at, updated_at FROM posts WHERE id=$1`
	var post model.Post
	if err := r.db.QueryRow(ctx, query, id).Scan(&post.Id, &post.Title, &post.Content, &post.CreatedAt, &post.UpdatedAt); err != nil {
		return model.Post{}, err
	}
	return post, nil
}

func (r *PostRepository) CreatePost(ctx context.Context, title, content string) (model.Post, error) {
	query := `INSERT INTO posts (title, content) VALUES ($1, $2) RETURNING id, title, content, created_at, updated_at`
	var newPost model.Post
	if err := r.db.QueryRow(ctx, query, title, content).Scan(&newPost.Id, &newPost.Title, &newPost.Content, &newPost.CreatedAt, &newPost.UpdatedAt); err != nil {
		return model.Post{}, err
	}
	return newPost, nil
}

func (r *PostRepository) UpdatePost(ctx context.Context, title, content string, id int) (model.Post, error) {
	query := `UPDATE posts SET title=$1, content=$2, updated_at=NOW() WHERE id=$3 RETURNING id, title, content, created_at, updated_at`
	var updatedPost model.Post
	if err := r.db.QueryRow(ctx, query, title, content, id).Scan(&updatedPost.Id, &updatedPost.Title, &updatedPost.Content, &updatedPost.CreatedAt, &updatedPost.UpdatedAt); err != nil {
		return model.Post{}, err
	}
	return updatedPost, nil
}

func (r *PostRepository) DeletePost(ctx context.Context, id int) error {
	query := `DELETE FROM posts WHERE id=$1 `
	_, err := r.db.Exec(ctx, query, id)
	return err
}
