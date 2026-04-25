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

func (s *ShelterService) FindNearby(ctx context.Context, params models.NearbyParams) ([]models.ShelterWithDistance, error) {
	if params.RadiusKm <= 0 {
		params.RadiusKm = 10
	}
	if params.RadiusKm > 500 {
		params.RadiusKm = 500
	}
	return s.repo.FindNearby(ctx, params)
}
