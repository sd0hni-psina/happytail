package service

import (
	"context"

	"github.com/sd0hni-psina/happytail/internal/models"
)

type ShelterService struct {
	repo ShelterRepository
}

func NewShelterService(repo ShelterRepository) *ShelterService {
	return &ShelterService{repo: repo}
}

func (s *ShelterService) GetAllShelters(ctx context.Context) ([]models.Shelter, error) {
	return s.repo.GetAll(ctx)
}

func (s *ShelterService) GetShelterByID(ctx context.Context, id int) (*models.Shelter, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *ShelterService) CreateShelter(ctx context.Context, input models.CreateShelterInput) (*models.Shelter, error) {
	return s.repo.Create(ctx, input)
}
