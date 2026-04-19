package service

import (
	"context"

	"github.com/sd0hni-psina/happytail/internal/models"
)

type AnimalPhotoService struct {
	repo AnimalPhotoRepository
}

func NewAnimalPhotoService(repo AnimalPhotoRepository) *AnimalPhotoService {
	return &AnimalPhotoService{repo: repo}
}

func (s *AnimalPhotoService) AddPhoto(ctx context.Context, input models.AnimalPhotoInput) (*models.AnimalPhoto, error) {
	return s.repo.Add(ctx, input)
}

func (s *AnimalPhotoService) DeletePhoto(ctx context.Context, photoID int) error {
	return s.repo.Delete(ctx, photoID)
}

func (s *AnimalPhotoService) MakeMainPhoto(ctx context.Context, animalID, photoID int) error {
	return s.repo.MakeMain(ctx, animalID, photoID)
}

func (s *AnimalPhotoService) GetAllPhotos(ctx context.Context, animalID int) ([]models.AnimalPhoto, error) {
	return s.repo.GetAll(ctx, animalID)
}
