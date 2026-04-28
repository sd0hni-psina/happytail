package service

import (
	"context"
	"fmt"
	"log/slog"
	"mime/multipart"

	"github.com/sd0hni-psina/happytail/internal/cache"
	"github.com/sd0hni-psina/happytail/internal/models"
)

type AnimalPhotoService struct {
	repo    AnimalPhotoRepository
	storage PhotoStorage
	cache   *cache.Cache
}

func NewAnimalPhotoService(repo AnimalPhotoRepository, storage PhotoStorage, cache *cache.Cache) *AnimalPhotoService {
	return &AnimalPhotoService{repo: repo, storage: storage, cache: cache}
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

	if s.cache != nil {
		if err := s.cache.DeleteByPattern(ctx, "photos:*"); err != nil {
			slog.Error("failed to invalidate photos cache", "error", err)
		} else {
			slog.Info("photos cache invalidated")
		}
	}

	return s.repo.Add(ctx, input)
}

func (s *AnimalPhotoService) DeletePhoto(ctx context.Context, photoID int, animalID int) error {
	if err := s.repo.Delete(ctx, photoID); err != nil {
		return err
	}

	if s.cache != nil {
		if err := s.cache.DeleteByPattern(ctx, fmt.Sprintf("photos:animal:%d", animalID)); err != nil {
			slog.Error("failed to invalidate photos cache", "error", err)
		} else {
			slog.Info("photos cache invalidated")
		}
	}

	return nil
}

func (s *AnimalPhotoService) MakeMainPhoto(ctx context.Context, animalID, photoID int) error {
	return s.repo.MakeMain(ctx, animalID, photoID)
}

func (s *AnimalPhotoService) GetAllPhotos(ctx context.Context, animalID int) ([]models.AnimalPhoto, error) {
	return s.repo.GetAll(ctx, animalID)
}
