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

const postCacheTTL = 10 * time.Minute

type PostService struct {
	repo  PostRepository
	cache *cache.Cache
}

func NewPostService(repo PostRepository, cache *cache.Cache) *PostService {
	return &PostService{repo: repo, cache: cache}
}

func (s *PostService) GetAllPost(ctx context.Context, params models.PaginationParams) ([]models.Post, int, error) {
	cacheKey := fmt.Sprintf("posts:page=%d:limit=%d", params.Page, params.Limit)

	type cacheResult struct {
		Posts []models.Post `json:"posts"`
		Total int           `json:"total"`
	}

	if s.cache != nil {
		if cached, err := s.cache.Get(ctx, cacheKey); err == nil {
			var result cacheResult
			if err := json.Unmarshal([]byte(cached), &result); err == nil {
				return result.Posts, result.Total, nil
			}
		}
	}

	posts, total, err := s.repo.GetAll(ctx, params.Limit, params.Offset())
	if err != nil {
		return nil, 0, err
	}

	if s.cache != nil {
		go func() {
			data, err := json.Marshal(cacheResult{Posts: posts, Total: total})
			if err != nil {
				return
			}
			s.cache.Set(context.Background(), cacheKey, string(data), postCacheTTL)
		}()
	}
	return posts, total, nil
}

func (s *PostService) GetPostByID(ctx context.Context, id int) (*models.Post, error) {
	cacheKey := fmt.Sprintf("posts:id:%d", id)
	if s.cache != nil {
		cached, err := s.cache.Get(ctx, cacheKey)
		if err == nil {
			var post models.Post
			if err := json.Unmarshal([]byte(cached), &post); err == nil {
				return &post, nil
			}
		}
	}

	posts, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if s.cache != nil {
		go func() {
			data, err := json.Marshal(posts)
			if err != nil {
				return
			}
			s.cache.Set(context.Background(), cacheKey, string(data), postCacheTTL)
		}()
	}
	return posts, nil
}

func (s *PostService) CreatePost(ctx context.Context, input models.CreatePostInput) (*models.Post, error) {
	post, err := s.repo.Create(ctx, input)
	if err != nil {
		return nil, err
	}
	if s.cache != nil {
		if err := s.cache.DeleteByPattern(ctx, "posts:*"); err != nil {
			slog.Error("failed to invalidate posts cache", "error", err)
		} else {
			slog.Info("posts cache invalidated")
		}
	}
	return post, nil
}
