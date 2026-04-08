package handler

import (
	"context"

	"github.com/sd0hni-psina/happytail/internal/models"
)

type AnimalService interface {
	GetAllAnimals(ctx context.Context) ([]models.Animal, error)
	GetAnimalByID(ctx context.Context, id int) (*models.Animal, error)
	CreateAnimal(ctx context.Context, input models.CreateAnimalInput) (*models.Animal, error)
}

type ShelterService interface {
	GetAllShelters(ctx context.Context) ([]models.Shelter, error)
	GetShelterByID(ctx context.Context, id int) (*models.Shelter, error)
	CreateShelter(ctx context.Context, input models.CreateShelterInput) (*models.Shelter, error)
}	

type UserService interface {
	GetAllUsers(ctx context.Context) ([]models.User, error)
	GetUserByID(ctx context.Context, id int) (*models.User, error)
	CreateUser(ctx context.Context, input models.CreateUserInput) (*models.User, error)
}	