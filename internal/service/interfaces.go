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

type ShelterRepository interface {
	GetAll(ctx context.Context) ([]models.Shelter, error)
	GetByID(ctx context.Context, id int) (*models.Shelter, error)
	Create(ctx context.Context, input models.CreateShelterInput) (*models.Shelter, error)
}

func (s *ShelterService) GetAllShelters(ctx context.Context) ([]models.Shelter, error) {
	return s.repo.GetAll(ctx)
}

func (s *ShelterService) GetShelterByID(ctx context.Context, id int) (*models.Shelter, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *ShelterService) CreateShelter(ctx context.Context, input models.CreateShelterInput) (*models.Shelter, error) {
	return s.repo.Create(ctx, input)
}

type UserRepository interface {
	GetAll(ctx context.Context) ([]models.User, error)
	GetByID(ctx context.Context, id int) (*models.User, error)
	Create(ctx context.Context, input models.CreateUserInput) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
}

func (s *UserService) GetAllUsers(ctx context.Context) ([]models.User, error) {
	return s.repo.GetAll(ctx)
}

func (s *UserService) GetUserByID(ctx context.Context, id int) (*models.User, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *UserService) CreateUser(ctx context.Context, input models.CreateUserInput) (*models.User, error) {
	return s.repo.Create(ctx, input)
}
