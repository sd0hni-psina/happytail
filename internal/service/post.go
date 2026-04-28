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

func (s *PostService) GetAllPost(ctx context.Context) ([]models.Post, error) {
	cacheKey := "posts:all"

	if s.cache != nil {
		cached, err := s.cache.Get(ctx, cacheKey)
		if err == nil {
			var posts []models.Post
			if err := json.Unmarshal([]byte(cached), &posts); err == nil {
				return posts, nil
			}
		}
	}
	posts, err := s.repo.GetAll(ctx)
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
