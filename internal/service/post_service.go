package service

import (
	"belajar-backend-golang/internal/model"
	"belajar-backend-golang/internal/repository"
	"context"
)

type PostService struct {
	Repository *repository.PostRepository
}

func NewPostService(repository *repository.PostRepository) *PostService {
	return &PostService{Repository: repository}
}

func (s *PostService) GetAllPosts(ctx context.Context) ([]model.Post, error) {
	return s.Repository.GetAllPosts(ctx)
}

func (s *PostService) GetPostById(ctx context.Context, id int) (model.Post, error) {
	return s.Repository.GetPostById(ctx, id)
}

func (s *PostService) CreatePost(ctx context.Context, title, content string) (model.Post, error) {
	return s.Repository.CreatePost(ctx, title, content)
}

func (s *PostService) UpdatePost(ctx context.Context, title, content string, id int) (model.Post, error) {
	return s.Repository.UpdatePost(ctx, title, content, id)
}

func (s *PostService) DeletePost(ctx context.Context, id int) error {
	return s.Repository.DeletePost(ctx, id)
}
