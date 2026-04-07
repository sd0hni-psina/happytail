package handler

import (
	"context"

	"github.com/sd0hni-psina/happytail/internal/models"
)

type AnimalRepository interface {
	GetAll(ctx context.Context) ([]models.Animal, error)
	GetByID(ctx context.Context, id int) (*models.Animal, error)
	Create(ctx context.Context, input models.CreateAnimalInput) (*models.Animal, error)
}
