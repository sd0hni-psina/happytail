package service

import (
	"context"

	"github.com/sd0hni-psina/happytail/internal/models"
)

type AdoptionService struct {
	repo       AdoptionRepository
	animalRepo AnimalRepository
}

func NewAdoptionService(repo AdoptionRepository, animalRepo AnimalRepository) *AdoptionService {
	return &AdoptionService{repo: repo, animalRepo: animalRepo}
}

func (s *AdoptionService) CreateAdoption(ctx context.Context, userID, animalID int) (*models.Adoption, error) {
	animal, err := s.animalRepo.GetByID(ctx, animalID)
	if err != nil {
		return nil, err
	}
	if animal == nil || animal.Status != "available" {
		return nil, models.ErrNotAvailable
	}
	return s.repo.Create(ctx, userID, animalID)
}
