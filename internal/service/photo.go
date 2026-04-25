package service

import (
	"context"
	"mime/multipart"

	"github.com/sd0hni-psina/happytail/internal/models"
)

type AnimalPhotoService struct {
	repo    AnimalPhotoRepository
	storage PhotoStorage
}

func NewAnimalPhotoService(repo AnimalPhotoRepository, storage PhotoStorage) *AnimalPhotoService {
	return &AnimalPhotoService{repo: repo, storage: storage}
}

func (s *AnimalPhotoService) AddPhoto(ctx context.Context, animalID int, file multipart.File, header *multipart.FileHeader, isMain bool) (*models.AnimalPhoto, error) {
	url, err := s.storage.Upload(ctx, file, header)
	if err != nil {
		return nil, err
	}

	input := models.AnimalPhotoInput{
		AnimalID: animalID,
		URL:      url,
		IsMain:   isMain,
	}
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
