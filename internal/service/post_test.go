package service

import (
	"context"
	"errors"
	"testing"

	"github.com/sd0hni-psina/happytail/internal/models"
)

// ---- mock репозитория постов ----

type mockPostRepo struct {
	post  *models.Post
	posts []models.Post
	total int
	err   error
	// сохраняем последний вызов UpdateStatus чтобы проверить аргументы
	lastUpdatedStatus models.PostStatus
}

func (m *mockPostRepo) GetAll(ctx context.Context, limit, offset int) ([]models.Post, int, error) {
	return m.posts, m.total, m.err
}

func (m *mockPostRepo) GetByID(ctx context.Context, id int) (*models.Post, error) {
	return m.post, m.err
}

func (m *mockPostRepo) GetByUserID(ctx context.Context, userID, limit, offset int) ([]models.Post, int, error) {
	return m.posts, m.total, m.err
}

func (m *mockPostRepo) Create(ctx context.Context, input models.CreatePostInput) (*models.Post, error) {
	return m.post, m.err
}

func (m *mockPostRepo) UpdateStatus(ctx context.Context, postID int, status models.PostStatus) error {
	m.lastUpdatedStatus = status
	return m.err
}

// ---- mock репозитория ролей для PostService ----

type mockRoleRepoForPost struct {
	isAdmin bool
	err     error
}

func (m *mockRoleRepoForPost) Appoint(ctx context.Context, input models.RoleInput) (*models.Role, error) {
	return nil, nil
}

func (m *mockRoleRepoForPost) Remove(ctx context.Context, roleID int) error {
	return nil
}

func (m *mockRoleRepoForPost) HasRole(ctx context.Context, userID int, role models.RoleType, shelterID *int) (bool, error) {
	return m.isAdmin, m.err
}

// ---- тесты GetPostByID ----

func TestGetPostByID_Success(t *testing.T) {
	expected := &models.Post{ID: 1, UserID: 10, AnimalID: 5}
	repo := &mockPostRepo{post: expected}
	svc := NewPostService(repo, &mockRoleRepoForPost{}, nil)

	post, err := svc.GetPostByID(context.Background(), 1)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if post == nil {
		t.Fatal("expected post, got nil")
	}
	if post.ID != expected.ID {
		t.Errorf("expected ID %d, got %d", expected.ID, post.ID)
	}
}

func TestGetPostByID_NotFound(t *testing.T) {
	repo := &mockPostRepo{err: models.ErrNotFound}
	svc := NewPostService(repo, &mockRoleRepoForPost{}, nil)

	post, err := svc.GetPostByID(context.Background(), 999)

	if !errors.Is(err, models.ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
	if post != nil {
		t.Errorf("expected nil post, got %v", post)
	}
}

// ---- тесты CreatePost ----

func TestCreatePost_Success(t *testing.T) {
	expected := &models.Post{ID: 1, UserID: 10, AnimalID: 5, Status: "active"}
	repo := &mockPostRepo{post: expected}
	svc := NewPostService(repo, &mockRoleRepoForPost{}, nil)

	input := models.CreatePostInput{
		UserID:      10,
		AnimalID:    5,
		ListingType: models.ListingTypeGive,
		ContactInfo: "test@test.com",
	}

	post, err := svc.CreatePost(context.Background(), input)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if post == nil {
		t.Fatal("expected post, got nil")
	}
	if post.UserID != expected.UserID {
		t.Errorf("expected UserID %d, got %d", expected.UserID, post.UserID)
	}
}

func TestCreatePost_DBError(t *testing.T) {
	repo := &mockPostRepo{err: errors.New("db error")}
	svc := NewPostService(repo, &mockRoleRepoForPost{}, nil)

	_, err := svc.CreatePost(context.Background(), models.CreatePostInput{})

	if err == nil {
		t.Error("expected error, got nil")
	}
}

// ---- тесты UpdateStatus ----

func TestUpdateStatus_AuthorCanDeactivate(t *testing.T) {
	// Автор может менять active → inactive
	post := &models.Post{ID: 1, UserID: 10}
	repo := &mockPostRepo{post: post}
	roleRepo := &mockRoleRepoForPost{isAdmin: false}
	svc := NewPostService(repo, roleRepo, nil)

	err := svc.UpdateStatus(context.Background(), 1, 10, models.PostStatusInactive)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if repo.lastUpdatedStatus != models.PostStatusInactive {
		t.Errorf("expected status inactive, got %s", repo.lastUpdatedStatus)
	}
}

func TestUpdateStatus_AuthorCannotDelete(t *testing.T) {
	// Автор НЕ может ставить deleted — только admin
	post := &models.Post{ID: 1, UserID: 10}
	repo := &mockPostRepo{post: post}
	roleRepo := &mockRoleRepoForPost{isAdmin: false}
	svc := NewPostService(repo, roleRepo, nil)

	err := svc.UpdateStatus(context.Background(), 1, 10, models.PostStatusDeleted)

	if !errors.Is(err, models.ErrForbidden) {
		t.Errorf("expected ErrForbidden, got %v", err)
	}
}

func TestUpdateStatus_StrangerCannotUpdate(t *testing.T) {
	// Чужой пост — userID 99 пытается изменить пост userID 10
	post := &models.Post{ID: 1, UserID: 10}
	repo := &mockPostRepo{post: post}
	roleRepo := &mockRoleRepoForPost{isAdmin: false}
	svc := NewPostService(repo, roleRepo, nil)

	err := svc.UpdateStatus(context.Background(), 1, 99, models.PostStatusInactive)

	if !errors.Is(err, models.ErrForbidden) {
		t.Errorf("expected ErrForbidden, got %v", err)
	}
}

func TestUpdateStatus_AdminCanDelete(t *testing.T) {
	// Admin может ставить любой статус включая deleted
	post := &models.Post{ID: 1, UserID: 10}
	repo := &mockPostRepo{post: post}
	roleRepo := &mockRoleRepoForPost{isAdmin: true}
	svc := NewPostService(repo, roleRepo, nil)

	err := svc.UpdateStatus(context.Background(), 1, 42, models.PostStatusDeleted)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if repo.lastUpdatedStatus != models.PostStatusDeleted {
		t.Errorf("expected status deleted, got %s", repo.lastUpdatedStatus)
	}
}

func TestUpdateStatus_AdminCanUpdateAnyPost(t *testing.T) {
	// Admin может менять чужие посты
	post := &models.Post{ID: 1, UserID: 10}
	repo := &mockPostRepo{post: post}
	roleRepo := &mockRoleRepoForPost{isAdmin: true}
	svc := NewPostService(repo, roleRepo, nil)

	// userID 42 — не автор, но admin
	err := svc.UpdateStatus(context.Background(), 1, 42, models.PostStatusInactive)

	if err != nil {
		t.Errorf("expected no error for admin, got %v", err)
	}
}

func TestUpdateStatus_PostNotFound(t *testing.T) {
	// Пост не существует — GetByID вернул ErrNotFound
	repo := &mockPostRepo{err: models.ErrNotFound}
	roleRepo := &mockRoleRepoForPost{isAdmin: false}
	svc := NewPostService(repo, roleRepo, nil)

	err := svc.UpdateStatus(context.Background(), 999, 10, models.PostStatusInactive)

	if !errors.Is(err, models.ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}
