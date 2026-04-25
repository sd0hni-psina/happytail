package service

import (
	"context"
	"log/slog"

	"github.com/sd0hni-psina/happytail/internal/models"
)

type AdoptionService struct {
	repo       AdoptionRepository
	userRepo   UserRepository
	AnimalRepo AnimalRepository
	notifier   Notifier
}

func NewAdoptionService(repo AdoptionRepository, userRepo UserRepository, animalRepo AnimalRepository, notifier Notifier) *AdoptionService {
	return &AdoptionService{repo: repo, userRepo: userRepo, AnimalRepo: animalRepo, notifier: notifier}
}

func (s *AdoptionService) CreateAdoption(ctx context.Context, userID, animalID int) (*models.Adoption, error) {
	adoption, err := s.repo.Create(ctx, userID, animalID)
	if err != nil {
		return nil, err
	}

	go func() {
		bgCtx := context.Background()
		user, err := s.userRepo.GetByID(bgCtx, userID)
		if err != nil {
			slog.Error("failed to get user got notification", "error", err, "user_id", userID)
			return
		}
		animal, err := s.AnimalRepo.GetByID(bgCtx, animalID)
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
