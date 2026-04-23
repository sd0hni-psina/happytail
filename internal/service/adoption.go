package service

import (
	"context"

	"github.com/sd0hni-psina/happytail/internal/models"
)

type AdoptionService struct {
	repo AdoptionRepository
}

func NewAdoptionService(repo AdoptionRepository) *AdoptionService {
	return &AdoptionService{repo: repo}
}

func (s *AdoptionService) CreateAdoption(ctx context.Context, userID, animalID int) (*models.Adoption, error) {
	return s.repo.Create(ctx, userID, animalID)
}
