package service

import (
	"context"

	"github.com/sd0hni-psina/happytail/internal/models"
)

type AnimalRepository interface {
	GetAll(ctx context.Context) ([]models.Animal, error)
	GetByID(ctx context.Context, id int) (*models.Animal, error)
	Create(ctx context.Context, input models.CreateAnimalInput) (*models.Animal, error)
}

func (s *AnimalService) GetAllAnimals(ctx context.Context) ([]models.Animal, error) {
	return s.repo.GetAll(ctx)
}

func (s *AnimalService) GetAnimalByID(ctx context.Context, id int) (*models.Animal, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *AnimalService) CreateAnimal(ctx context.Context, input models.CreateAnimalInput) (*models.Animal, error) {
	return s.repo.Create(ctx, input)
}
