package handler

import (
	"context"

	"github.com/sd0hni-psina/happytail/internal/models"
)

type AnimalService interface {
	GetAllAnimals(ctx context.Context, params models.PaginationParams, filter models.FilterParams) ([]models.Animal, int, error)
	GetAnimalByID(ctx context.Context, id int) (*models.Animal, error)
	CreateAnimal(ctx context.Context, input models.CreateAnimalInput) (*models.Animal, error)
}

type ShelterService interface {
	GetAllShelters(ctx context.Context) ([]models.Shelter, error)
	GetShelterByID(ctx context.Context, id int) (*models.Shelter, error)
	CreateShelter(ctx context.Context, input models.CreateShelterInput) (*models.Shelter, error)
	FindNearby(ctx context.Context, params models.NearbyParams) ([]models.ShelterWithDistance, error)
}

type UserService interface {
	GetAllUsers(ctx context.Context) ([]models.UserPublic, error)
	GetUserByID(ctx context.Context, id int) (*models.UserPublic, error)
	CreateUser(ctx context.Context, input models.CreateUserInput) (*models.User, error)
	Login(ctx context.Context, email, password string) (*models.AuthResponse, error)
	Refresh(ctx context.Context, refreshToken string) (*models.AuthResponse, error)
	Logout(ctx context.Context, refreshToken string) error
}

type AdoptionService interface {
	CreateAdoption(ctx context.Context, userID, animalID int) (*models.Adoption, error)
}

type PostService interface {
	GetAllPost(ctx context.Context) ([]models.Post, error)
	GetPostByID(ctx context.Context, id int) (*models.Post, error)
	CreatePost(ctx context.Context, input models.CreatePostInput) (*models.Post, error)
}

type AnimalPhotoService interface {
	AddPhoto(ctx context.Context, input models.AnimalPhotoInput) (*models.AnimalPhoto, error)
	DeletePhoto(ctx context.Context, photoID int) error
	MakeMainPhoto(ctx context.Context, animalID, photoID int) error
	GetAllPhotos(ctx context.Context, animalID int) ([]models.AnimalPhoto, error)
}

type RoleService interface {
	AppointRole(ctx context.Context, input models.RoleInput) (*models.Role, error)
	RemoveRole(ctx context.Context, roleID int) error
	HasRole(ctx context.Context, userID int, role models.RoleType, shelterID *int) (bool, error)
}
