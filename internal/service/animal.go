package service

import (
	"context"

	"github.com/sd0hni-psina/happytail/internal/models"
)

type AnimalService struct {
	repo AnimalRepository
}

func NewAnimalService(repo AnimalRepository) *AnimalService {
	return &AnimalService{repo: repo}
}

func (s *AnimalService) GetAllAnimals(ctx context.Context, params models.PaginationParams) ([]models.Animal, int, error) {
	return s.repo.GetAll(ctx, params.Limit, params.Offset())
}

func (s *AnimalService) GetAnimalByID(ctx context.Context, id int) (*models.Animal, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *AnimalService) CreateAnimal(ctx context.Context, input models.CreateAnimalInput) (*models.Animal, error) {
	return s.repo.Create(ctx, input)
}
