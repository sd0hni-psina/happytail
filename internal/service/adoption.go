package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/sd0hni-psina/happytail/internal/cache"
	"github.com/sd0hni-psina/happytail/internal/models"
)

type AdoptionService struct {
	repo       AdoptionRepository
	userRepo   UserRepository
	animalRepo AnimalRepository
	notifier   Notifier
	cache      *cache.Cache
}

func NewAdoptionService(repo AdoptionRepository, userRepo UserRepository, animalRepo AnimalRepository, notifier Notifier, cache *cache.Cache) *AdoptionService {
	return &AdoptionService{repo: repo, userRepo: userRepo, animalRepo: animalRepo, notifier: notifier, cache: cache}
}

func (s *AdoptionService) CreateAdoption(ctx context.Context, userID, animalID int) (*models.Adoption, error) {
	adoption, err := s.repo.Create(ctx, userID, animalID)
	if err != nil {
		return nil, err
	}
	bgCtx := context.Background()

	specificKey := fmt.Sprintf("animals:id:%d", animalID)
	if err := s.cache.Delete(bgCtx, specificKey); err != nil {
		slog.Error("failed to invalidate animal cache", "error", err, "animal_id", animalID)
	}
	if err := s.cache.DeleteByPattern(bgCtx, "animals:page=*"); err != nil {
		slog.Error("failed to invalidate animals list cache", "error", err)
	}
	// реализовать очередь задач (outbox pattern)
	go func() {
		user, err := s.userRepo.GetByID(bgCtx, userID)
		if err != nil {
			slog.Error("failed to get user got notification", "error", err, "user_id", userID)
			return
		}
		animal, err := s.animalRepo.GetByID(bgCtx, animalID)
		if err != nil {
			slog.Error("failed to get animal for notification", "error", err, "animal_id", animalID)
			return
		}

		if err := s.notifier.SendAdoptionConfirmation(user.Email, user.FullName, animal.Name); err != nil {
			slog.Error("failed to send adoption email", "error", err, "user_id", userID)
			return
		}

		slog.Info("adoption confirmation email sent", "user_id", userID, "animal_id", animalID)
	}()
	return adoption, nil
}

func (s *AdoptionService) GetByUserID(ctx context.Context, userID int) ([]models.Adoption, error) {
	return s.repo.GetByUserID(ctx, userID)
}
