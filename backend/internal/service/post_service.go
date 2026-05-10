package service

import (
	"context"

	"github.com/AzafaDev/golang-web-server-and-rest-api.git/internal/models"
	"github.com/AzafaDev/golang-web-server-and-rest-api.git/internal/repository"
)

type PostService struct {
	repo repository.PostRepository
}

func NewPostService(repo *repository.PostRepository) *PostService {
	return &PostService{repo: *repo} 
}

func (s *PostService) GetAll(ctx context.Context) ([]models.Post, error) {
	return s.repo.GetAll(ctx)
}

func (s *PostService) GetPostByID(ctx context.Context, id int) (models.Post, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *PostService) CreatePost(ctx context.Context, title, content string) (models.Post, error) {
	return s.repo.Create(ctx, title, content)
}

func (s *PostService) UpdatePost(ctx context.Context, p models.Post) (models.Post, error) {
	return s.repo.Update(ctx, p)
}

func (s *PostService) DeletePost(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}
