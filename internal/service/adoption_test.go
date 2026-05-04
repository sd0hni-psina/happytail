package service

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/sd0hni-psina/happytail/internal/models"
)

type mockAdoptionRepo struct {
	adoption *models.Adoption
	err      error
}

func (m *mockAdoptionRepo) Create(ctx context.Context, userID, animalID int) (*models.Adoption, error) {
	return m.adoption, m.err
}

func (m *mockAdoptionRepo) GetByUserID(ctx context.Context, userID int) ([]models.Adoption, error) {
	return nil, m.err
}

type mockUserRepo struct {
	user *models.UserPublic
	err  error
}

func (m *mockUserRepo) GetAll(ctx context.Context) ([]models.UserPublic, error) { return nil, nil }
func (m *mockUserRepo) GetByID(ctx context.Context, id int) (*models.UserPublic, error) {
	return m.user, m.err
}
func (m *mockUserRepo) Create(ctx context.Context, input models.CreateUserInput) (*models.User, error) {
	return nil, nil
}
func (m *mockUserRepo) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	return nil, nil
}

type mockNotifier struct{}

func (m *mockNotifier) SendAdoptionConfirmation(toEmail, userName, animalName string) error {
	return nil
}

func TestCreateAdoption_Success(t *testing.T) {
	t.Log("before CreateAdoption")
	expected := &models.Adoption{ID: 1, UserID: 10, AnimalID: 5}

	svc := NewAdoptionService(
		&mockAdoptionRepo{adoption: expected},
		&mockUserRepo{user: &models.UserPublic{ID: 10, Email: "test@test.com", FullName: "Test"}},
		&mockAnimalRepo{animal: &models.Animal{ID: 5, Name: "Barsik"}},
		&mockNotifier{},
		nil,
	)
	fmt.Println("step 1")
	adoption, err := svc.CreateAdoption(context.Background(), 10, 5)
	fmt.Println("step 2")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if adoption == nil {
		t.Fatal("expected adoption, got nil")
	}
	if adoption.ID != expected.ID {
		t.Errorf("expected ID %d, got %d", expected.ID, adoption.ID)
	}
}

func TestCreateAdoption_NotFound(t *testing.T) {
	svc := NewAdoptionService(
		&mockAdoptionRepo{err: models.ErrNotFound},
		&mockUserRepo{},
		&mockAnimalRepo{},
		&mockNotifier{},
		nil,
	)

	_, err := svc.CreateAdoption(context.Background(), 10, 999)

	if !errors.Is(err, models.ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestCreateAdoption_NotAvailable(t *testing.T) {
	svc := NewAdoptionService(
		&mockAdoptionRepo{err: models.ErrNotAvailable},
		&mockUserRepo{},
		&mockAnimalRepo{},
		&mockNotifier{},
		nil,
	)

	_, err := svc.CreateAdoption(context.Background(), 10, 5)

	if !errors.Is(err, models.ErrNotAvailable) {
		t.Errorf("expected ErrNotAvailable, got %v", err)
	}
}

func TestCreateAdoption_DBError(t *testing.T) {
	svc := NewAdoptionService(
		&mockAdoptionRepo{err: errors.New("connection refused")},
		&mockUserRepo{},
		&mockAnimalRepo{},
		&mockNotifier{},
		nil,
	)

	adoption, err := svc.CreateAdoption(context.Background(), 10, 5)

	if err == nil {
		t.Error("expected error, got nil")
	}
	if adoption != nil {
		t.Errorf("expected nil adoption, got %v", adoption)
	}
}
