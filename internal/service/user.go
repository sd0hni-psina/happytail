package service

import (
	"context"

	"github.com/sd0hni-psina/happytail/internal/models"
)

type UserService struct {
	repo      UserRepository
	jwtSecret string
}

func NewUserService(repo UserRepository, jwtSecret string) *UserService {
	return &UserService{repo: repo, jwtSecret: jwtSecret}
}

func (s *UserService) GetAllUsers(ctx context.Context) ([]models.UserPublic, error) {
	return s.repo.GetAll(ctx)
}

func (s *UserService) GetUserByID(ctx context.Context, id int) (*models.UserPublic, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *UserService) CreateUser(ctx context.Context, input models.CreateUserInput) (*models.User, error) {
	return s.repo.Create(ctx, input)
}
