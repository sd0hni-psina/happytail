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

const animalCacheTTL = 5 * time.Minute

func strVal(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func shelterIDVal(id *int) int {
	if id == nil {
		return 0
	}
	return *id
}

type AnimalService struct {
	repo  AnimalRepository
	cache *cache.Cache
}

func NewAnimalService(repo AnimalRepository, cache *cache.Cache) *AnimalService {
	return &AnimalService{repo: repo, cache: cache}
}

func (s *AnimalService) GetAllAnimals(ctx context.Context, params models.PaginationParams, filter models.FilterParams) ([]models.Animal, int, error) {
	// Name добавлен в cache key — поиск "барсик" и поиск "мурзик" кэшируются отдельно
	cacheKey := fmt.Sprintf("animals:page=%d:limit=%d:name=%s:type=%s:status=%s:breed=%s:shelter=%d",
		params.Page, params.Limit,
		strVal(filter.Name),
		strVal(filter.Type),
		strVal(filter.Status),
		strVal(filter.Breed),
		shelterIDVal(filter.ShelterID),
	)

	type cacheResult struct {
		Animals []models.Animal `json:"animals"`
		Total   int             `json:"total"`
	}

	if s.cache != nil {
		if cached, err := s.cache.Get(ctx, cacheKey); err == nil {
			var result cacheResult
			if err := json.Unmarshal([]byte(cached), &result); err == nil {
				return result.Animals, result.Total, nil
			}
		}
	}

	animals, total, err := s.repo.GetAll(ctx, params.Limit, params.Offset(), filter)
	if err != nil {
		return nil, 0, err
	}

	if s.cache != nil {
		go func() {
			data, err := json.Marshal(cacheResult{Animals: animals, Total: total})
			if err != nil {
				return
			}
			s.cache.Set(context.Background(), cacheKey, string(data), animalCacheTTL)
		}()
	}

	return animals, total, nil
}

func (s *AnimalService) GetAnimalByID(ctx context.Context, id int) (*models.Animal, error) {
	cacheKey := fmt.Sprintf("animals:id:%d", id)

	if s.cache != nil {
		cached, err := s.cache.Get(ctx, cacheKey)
		if err == nil {
			var animal models.Animal
			if err = json.Unmarshal([]byte(cached), &animal); err == nil {
				return &animal, nil
			}
		}
	}

	animal, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if s.cache != nil {
		go func() {
			data, err := json.Marshal(animal)
			if err != nil {
				return
			}
			s.cache.Set(context.Background(), cacheKey, string(data), animalCacheTTL)
		}()
	}

	return animal, nil
}

func (s *AnimalService) CreateAnimal(ctx context.Context, input models.CreateAnimalInput) (*models.Animal, error) {
	animal, err := s.repo.Create(ctx, input)
	if err != nil {
		return nil, err
	}
	if s.cache != nil {
		if err := s.cache.DeleteByPattern(ctx, "animals:*"); err != nil {
			slog.Error("failed to invalidate animals cache", "error", err)
			// _ = err // не кретичсно удалить кеш, просто логируется в middleware.Logger
		} else {
			slog.Info("animals cache invalidated")
		}
	}
	return animal, nil
}

func (s *AnimalService) UpdateAnimal(ctx context.Context, id int, input models.UpdateAnimalInput) (*models.Animal, error) {
	animal, err := s.repo.Update(ctx, id, input)
	if err != nil {
		return nil, err
	}

	if s.cache != nil {
		specificKey := fmt.Sprintf("animals:id:%d", id)
		if err := s.cache.Delete(ctx, specificKey); err != nil {
			slog.Error("failed to invalidate animal cache", "error", err, "animal_id", id)
		}
		if err := s.cache.DeleteByPattern(ctx, "animals:page=*"); err != nil {
			slog.Error("failed to invalidate animals list cache", "error", err)
		}
	}
	return animal, nil
}

func (s *AnimalService) DeleteAnimal(ctx context.Context, id int) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return nil
	}
	if s.cache != nil {
		if err := s.cache.Delete(ctx, fmt.Sprintf("animals:id:%d", id)); err != nil {
			slog.Error("failed to invalidate animal cache", "error", err)
		}
		if err := s.cache.DeleteByPattern(ctx, "animals:page=*"); err != nil {
			slog.Error("failed to invalidate animals list cache", "error", err)
		}
	}
	return nil
}
