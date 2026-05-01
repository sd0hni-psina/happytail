package service

import (
	"context"
	"mime/multipart"
	"time"

	"github.com/sd0hni-psina/happytail/internal/models"
)

type AnimalRepository interface {
	GetAll(ctx context.Context, limit, offset int, filter models.FilterParams) ([]models.Animal, int, error)
	GetByID(ctx context.Context, id int) (*models.Animal, error)
	Create(ctx context.Context, input models.CreateAnimalInput) (*models.Animal, error)
	Update(ctx context.Context, id int, input models.UpdateAnimalInput) (*models.Animal, error)
	GetShelterIDByAnimalID(ctx context.Context, animalID int) (*int, error)
}

type ShelterRepository interface {
	GetAll(ctx context.Context, limit, offset int, filter models.FilterParams) ([]models.Shelter, int, error)
	GetByID(ctx context.Context, id int) (*models.Shelter, error)
	Create(ctx context.Context, input models.CreateShelterInput) (*models.Shelter, error)
	FindNearby(ctx context.Context, params models.NearbyParams) ([]models.ShelterWithDistance, error)
	Update(ctx context.Context, id int, input models.UpdateShelterInput) (*models.Shelter, error)
}

type UserRepository interface {
	GetAll(ctx context.Context) ([]models.UserPublic, error)
	GetByID(ctx context.Context, id int) (*models.UserPublic, error)
	Create(ctx context.Context, input models.CreateUserInput) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
}

type AdoptionRepository interface {
	Create(ctx context.Context, userID, animalID int) (*models.Adoption, error)
	GetByUserID(ctx context.Context, userID int) ([]models.Adoption, error)
}

type PostRepository interface {
	GetAll(ctx context.Context) ([]models.Post, error)
	GetByID(ctx context.Context, id int) (*models.Post, error)
	Create(ctx context.Context, input models.CreatePostInput) (*models.Post, error)
}

type AnimalPhotoRepository interface {
	Add(ctx context.Context, input models.AnimalPhotoInput) (*models.AnimalPhoto, error)
	Delete(ctx context.Context, photoID int) error
	MakeMain(ctx context.Context, animalID, photoID int) error
	GetAll(ctx context.Context, animalID int) ([]models.AnimalPhoto, error)
	GetByID(ctx context.Context, photoID int) (*models.AnimalPhoto, error)
}

type RoleRepository interface {
	Appoint(ctx context.Context, input models.RoleInput) (*models.Role, error)
	Remove(ctx context.Context, roleID int) error
	HasRole(ctx context.Context, userID int, role models.RoleType, shelterID *int) (bool, error)
}

type RefreshTokenRepository interface {
	Create(ctx context.Context, userID int, token string, expiresAt time.Time) error
	GetByToken(ctx context.Context, token string) (*models.RefreshToken, error)
	Revoke(ctx context.Context, token string) error
	RevokeAllForUser(ctx context.Context, userID int) error
}

type Notifier interface {
	SendAdoptionConfirmation(toEmail, userName, animalName string) error
}

type PhotoStorage interface {
	Upload(ctx context.Context, file multipart.File, header *multipart.FileHeader) (string, error)
	Delete(ctx context.Context, objectName string) error
}
