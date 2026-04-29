package service

import (
	"context"
	"fmt"
	"log/slog"
	"mime/multipart"
	"path"

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

	photo, err := s.repo.Add(ctx, input)
	if err != nil {
		objectName := path.Base(url)
		if delErr := s.storage.Delete(ctx, objectName); delErr != nil {
			slog.Error("failed to delete uploaded photo after DB error", "error", delErr, "object", objectName)
		} else {
			slog.Info("uploaded photo deleted after DB error", "object", objectName)
		}
		return nil, err
	}

	if s.cache != nil {
		cacheKey := fmt.Sprintf("photos:animal:%d", animalID)
		if err := s.cache.Delete(ctx, cacheKey); err != nil {
			slog.Error("failed to invalidate photos cache", "error", err)
		} else {
			slog.Info("photos cache invalidated")
		}
	}

	return photo, nil
}

func (s *AnimalPhotoService) DeletePhoto(ctx context.Context, photoID int, animalID int) error {
	photo, err := s.repo.GetByID(ctx, photoID)
	if err != nil {
		return err
	}

	objectName := path.Base(photo.URL)
	if err := s.storage.Delete(ctx, objectName); err != nil {
		slog.Error("failed to delete photo from storage", "error", err, "object", objectName)
	}

	if err := s.repo.Delete(ctx, photoID); err != nil {
		return err
	}

	if s.cache != nil {
		cacheKey := fmt.Sprintf("photos:animal:%d", animalID)
		if err := s.cache.Delete(ctx, cacheKey); err != nil {
			slog.Error("failed to invalidate photos cache", "error", err)
		} else {
			slog.Info("photos cache invalidated")
		}
	}

	return nil
}

func (s *AnimalPhotoService) MakeMainPhoto(ctx context.Context, animalID, photoID int) error {
	photo, err := s.repo.GetByID(ctx, photoID)
	if err != nil {
		return err
	}

	if photo.AnimalID != animalID {
		return models.ErrNotFound
	}
	return s.repo.MakeMain(ctx, animalID, photoID)
}

func (s *AnimalPhotoService) GetAllPhotos(ctx context.Context, animalID int) ([]models.AnimalPhoto, error) {
	return s.repo.GetAll(ctx, animalID)
}
