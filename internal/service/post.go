package service

import (
	"context"

	"github.com/sd0hni-psina/happytail/internal/models"
)

type PostService struct {
	repo PostRepository
}

func NewPostService(repo PostRepository) *PostService {
	return &PostService{repo: repo}
}

func (s *PostService) GetAllPost(ctx context.Context) ([]models.Post, error) {
	return s.repo.GetAll(ctx)
}

func (s *PostService) GetPostByID(ctx context.Context, id int) (*models.Post, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *PostService) CreatePost(ctx context.Context, input models.CreatePostInput) (*models.Post, error) {
	return s.repo.Create(ctx, input)
}
