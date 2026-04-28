package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/sd0hni-psina/happytail/internal/cache"
	"github.com/sd0hni-psina/happytail/internal/models"
)

const shelterCacheTTL = 10 * time.Minute
const shelterNearbyCacheTTL = 10 * time.Minute

type ShelterService struct {
	repo  ShelterRepository
	cache *cache.Cache
}

func NewShelterService(repo ShelterRepository, cache *cache.Cache) *ShelterService {
	return &ShelterService{repo: repo, cache: cache}
}

func (s *ShelterService) GetAllShelters(ctx context.Context) ([]models.Shelter, error) {
	cacheKey := "shelters:all"

	if s.cache != nil {
		cached, err := s.cache.Get(ctx, cacheKey)
		if err == nil {
			var shelters []models.Shelter
			if err := json.Unmarshal([]byte(cached), &shelters); err == nil {
				return shelters, nil
			}
		}
	}

	shelters, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	if s.cache != nil {
		go func() {
			data, err := json.Marshal(shelters)
			if err != nil {
				return
			}
			s.cache.Set(context.Background(), cacheKey, string(data), shelterCacheTTL)
		}()
	}
	return shelters, nil
}

func (s *ShelterService) GetShelterByID(ctx context.Context, id int) (*models.Shelter, error) {
	cacheKey := fmt.Sprintf("shelters:id:%d", id)

	if s.cache != nil {
		cached, err := s.cache.Get(ctx, cacheKey)
		if err == nil {
			var shelter models.Shelter
			if err := json.Unmarshal([]byte(cached), &shelter); err == nil {
				return &shelter, nil
			}
		}
	}

	shelter, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if s.cache != nil {
		go func() {
			data, err := json.Marshal(shelter)
			if err != nil {
				return
			}
			s.cache.Set(context.Background(), cacheKey, string(data), shelterCacheTTL)
		}()
	}
	return shelter, nil
}

func (s *ShelterService) CreateShelter(ctx context.Context, input models.CreateShelterInput) (*models.Shelter, error) {
	shelter, err := s.repo.Create(ctx, input)
	if err != nil {
		return nil, err
	}

	if s.cache != nil {
		if err := s.cache.DeleteByPattern(ctx, "shelters:*"); err != nil {
			slog.Error("failed to invalidate shelters cache", "error", err)
		} else {
			slog.Info("shelters cache invalidated")
		}
	}

	return shelter, nil
}

func (s *ShelterService) FindNearby(ctx context.Context, params models.NearbyParams) ([]models.ShelterWithDistance, error) {
	if params.RadiusKm <= 0 {
		params.RadiusKm = 10
	}
	if params.RadiusKm > 500 {
		params.RadiusKm = 500
	}

	cacheKey := fmt.Sprintf("shelters:nearby:lat=%.4f:lon=%.4f:radius=%.1f", params.Latitude, params.Longitude, params.RadiusKm)

	if s.cache != nil {
		cached, err := s.cache.Get(ctx, cacheKey)
		if err == nil {
			var result []models.ShelterWithDistance
			if err := json.Unmarshal([]byte(cached), &result); err == nil {
				return result, nil
			}
		}
	}

	shelters, err := s.repo.FindNearby(ctx, params)
	if err != nil {
		return nil, err
	}

	if s.cache != nil {
		go func() {
			data, err := json.Marshal(shelters)
			if err != nil {
				return
			}
			s.cache.Set(context.Background(), cacheKey, string(data), shelterNearbyCacheTTL)
			slog.Info("nearby shelters cache set", "cache_key", cacheKey)
		}()
	}

	return shelters, nil
}
