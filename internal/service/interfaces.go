package service

import (
	"context"

	"github.com/sd0hni-psina/happytail/internal/models"
)

type AnimalRepository interface {
	GetAll(ctx context.Context, limit, offset int, filter models.FilterParams) ([]models.Animal, int, error)
	GetByID(ctx context.Context, id int) (*models.Animal, error)
	Create(ctx context.Context, input models.CreateAnimalInput) (*models.Animal, error)
}

type ShelterRepository interface {
	GetAll(ctx context.Context) ([]models.Shelter, error)
	GetByID(ctx context.Context, id int) (*models.Shelter, error)
	Create(ctx context.Context, input models.CreateShelterInput) (*models.Shelter, error)
}

type UserRepository interface {
	GetAll(ctx context.Context) ([]models.User, error)
	GetByID(ctx context.Context, id int) (*models.User, error)
	Create(ctx context.Context, input models.CreateUserInput) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
}

type AdoptionRepository interface {
	Create(ctx context.Context, userID, animalID int) (*models.Adoption, error)
}

type PostRepository interface {
	GetAll(ctx context.Context) ([]models.Post, error)
	GetByID(ctx context.Context, id int) (*models.Post, error)
	Create(ctx context.Context, input models.CreatePostInput) (*models.Post, error)
}
