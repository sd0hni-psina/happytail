package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/sd0hni-psina/happytail/internal/models"
	"golang.org/x/crypto/bcrypt"
)

// ---- mock репозитория пользователей ----

type mockUserRepoFull struct {
	user   *models.User
	public *models.UserPublic
	users  []models.UserPublic
	err    error
}

func (m *mockUserRepoFull) GetAll(ctx context.Context) ([]models.UserPublic, error) {
	return m.users, m.err
}

func (m *mockUserRepoFull) GetByID(ctx context.Context, id int) (*models.UserPublic, error) {
	return m.public, m.err
}

func (m *mockUserRepoFull) Create(ctx context.Context, input models.CreateUserInput) (*models.User, error) {
	return m.user, m.err
}

func (m *mockUserRepoFull) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	return m.user, m.err
}

// ---- mock репозитория токенов ----

type mockTokenRepo struct {
	token *models.RefreshToken
	err   error
}

func (m *mockTokenRepo) Create(ctx context.Context, userID int, token string, expiresAt time.Time) error {
	return m.err
}

func (m *mockTokenRepo) GetByToken(ctx context.Context, token string) (*models.RefreshToken, error) {
	return m.token, m.err
}

func (m *mockTokenRepo) Revoke(ctx context.Context, token string) error {
	return m.err
}

func (m *mockTokenRepo) RevokeAllForUser(ctx context.Context, userID int) error {
	return m.err
}

// ---- вспомогательная функция — создаём хеш пароля для тестов ----

func mustHashPassword(t *testing.T, password string) string {
	t.Helper()
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}
	return string(hash)
}

// ---- тесты GetAllUsers ----

func TestGetAllUsers_Success(t *testing.T) {
	expected := []models.UserPublic{
		{ID: 1, FullName: "Айгерим"},
		{ID: 2, FullName: "Данияр"},
	}
	repo := &mockUserRepoFull{users: expected}
	svc := NewUserService(repo, &mockTokenRepo{}, "secret-32-chars-minimum-length!!", nil)

	users, err := svc.GetAllUsers(context.Background())

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(users) != 2 {
		t.Errorf("expected 2 users, got %d", len(users))
	}
}

func TestGetAllUsers_DBError(t *testing.T) {
	repo := &mockUserRepoFull{err: errors.New("db error")}
	svc := NewUserService(repo, &mockTokenRepo{}, "secret-32-chars-minimum-length!!", nil)

	_, err := svc.GetAllUsers(context.Background())

	if err == nil {
		t.Error("expected error, got nil")
	}
}

// ---- тесты GetUserByID ----

func TestGetUserByID_Success(t *testing.T) {
	expected := &models.UserPublic{ID: 1, FullName: "Айгерим", Email: "ai@test.kz"}
	repo := &mockUserRepoFull{public: expected}
	svc := NewUserService(repo, &mockTokenRepo{}, "secret-32-chars-minimum-length!!", nil)

	user, err := svc.GetUserByID(context.Background(), 1)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if user == nil {
		t.Fatal("expected user, got nil")
	}
	if user.FullName != expected.FullName {
		t.Errorf("expected %s, got %s", expected.FullName, user.FullName)
	}
}

func TestGetUserByID_NotFound(t *testing.T) {
	repo := &mockUserRepoFull{err: models.ErrNotFound}
	svc := NewUserService(repo, &mockTokenRepo{}, "secret-32-chars-minimum-length!!", nil)

	user, err := svc.GetUserByID(context.Background(), 999)

	if !errors.Is(err, models.ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
	if user != nil {
		t.Errorf("expected nil user, got %v", user)
	}
}

// ---- тесты CreateUser ----

func TestCreateUser_Success(t *testing.T) {
	expected := &models.User{ID: 1, FullName: "Данияр", Email: "d@test.kz"}
	repo := &mockUserRepoFull{user: expected}
	svc := NewUserService(repo, &mockTokenRepo{}, "secret-32-chars-minimum-length!!", nil)

	input := models.CreateUserInput{
		FullName:    "Данияр",
		Email:       "d@test.kz",
		Password:    "password123",
		PhoneNumber: "+77001234567",
	}

	user, err := svc.CreateUser(context.Background(), input)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if user == nil {
		t.Fatal("expected user, got nil")
	}
	if user.Email != expected.Email {
		t.Errorf("expected email %s, got %s", expected.Email, user.Email)
	}
}

func TestCreateUser_EmailConflict(t *testing.T) {
	repo := &mockUserRepoFull{err: models.ErrConflict}
	svc := NewUserService(repo, &mockTokenRepo{}, "secret-32-chars-minimum-length!!", nil)

	input := models.CreateUserInput{
		FullName:    "Данияр",
		Email:       "exists@test.kz",
		Password:    "password123",
		PhoneNumber: "+77001234567",
	}

	_, err := svc.CreateUser(context.Background(), input)

	if !errors.Is(err, models.ErrConflict) {
		t.Errorf("expected ErrConflict, got %v", err)
	}
}

// ---- тесты Login ----

func TestLogin_Success(t *testing.T) {
	password := "correctPassword1"
	hash := mustHashPassword(t, password)

	user := &models.User{ID: 1, Email: "ai@test.kz", PasswordHash: hash}
	repo := &mockUserRepoFull{user: user}
	tokenRepo := &mockTokenRepo{}
	svc := NewUserService(repo, tokenRepo, "secret-32-chars-minimum-length!!", nil)

	auth, err := svc.Login(context.Background(), "ai@test.kz", password)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if auth == nil {
		t.Fatal("expected auth response, got nil")
	}
	if auth.AccessToken == "" {
		t.Error("expected non-empty access token")
	}
	if auth.RefreshToken == "" {
		t.Error("expected non-empty refresh token")
	}
}

func TestLogin_UserNotFound(t *testing.T) {
	// GetByEmail возвращает nil, nil — пользователь не найден
	repo := &mockUserRepoFull{user: nil, err: nil}
	svc := NewUserService(repo, &mockTokenRepo{}, "secret-32-chars-minimum-length!!", nil)

	_, err := svc.Login(context.Background(), "notfound@test.kz", "password123")

	if err == nil {
		t.Error("expected error for unknown user, got nil")
	}
}

func TestLogin_WrongPassword(t *testing.T) {
	hash := mustHashPassword(t, "correctPassword1")
	user := &models.User{ID: 1, Email: "ai@test.kz", PasswordHash: hash}
	repo := &mockUserRepoFull{user: user}
	svc := NewUserService(repo, &mockTokenRepo{}, "secret-32-chars-minimum-length!!", nil)

	_, err := svc.Login(context.Background(), "ai@test.kz", "wrongPassword!")

	if err == nil {
		t.Error("expected error for wrong password, got nil")
	}
}

// ---- тесты Logout ----

func TestLogout_Success(t *testing.T) {
	repo := &mockUserRepoFull{}
	tokenRepo := &mockTokenRepo{}
	svc := NewUserService(repo, tokenRepo, "secret-32-chars-minimum-length!!", nil)

	// Генерируем валидный access токен для теста
	user := &models.User{ID: 1}
	accessToken, err := generateAccessToken(user, "secret-32-chars-minimum-length!!")
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	err = svc.Logout(context.Background(), accessToken, "some-refresh-token")

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestLogout_TokenRepoError(t *testing.T) {
	repo := &mockUserRepoFull{}
	tokenRepo := &mockTokenRepo{err: errors.New("db error")}
	svc := NewUserService(repo, tokenRepo, "secret-32-chars-minimum-length!!", nil)

	err := svc.Logout(context.Background(), "any-token", "refresh-token")

	if err == nil {
		t.Error("expected error when tokenRepo fails, got nil")
	}
}
