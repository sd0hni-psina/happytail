package service

import (
	"context"
	"errors"
	"testing"

	"github.com/sd0hni-psina/happytail/internal/models"
)

type mockAnimalRepo struct {
	animal *models.Animal
	err    error
}

func (m *mockAnimalRepo) GetByID(ctx context.Context, id int) (*models.Animal, error) {
	return m.animal, m.err
}

func (m *mockAnimalRepo) GetAll(ctx context.Context, limit, ofser int) ([]models.Animal, int, error) {
	return nil, 0, nil
}

func (m *mockAnimalRepo) Create(ctx context.Context, input models.CreateAnimalInput) (*models.Animal, error) {
	return m.animal, m.err
}

func TestGetAnimalByID(t *testing.T) {
	repo := &mockAnimalRepo{
		animal: nil,
		err:    models.ErrNotFound,
	}
	svc := NewAnimalService(repo)

	animal, err := svc.GetAnimalByID(context.Background(), 1)

	if !errors.Is(err, models.ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
	if animal != nil {
		t.Errorf("expected nil animal, got %v", animal)
	}
}

func TestCreateAnimal(t *testing.T) {
	var input models.CreateAnimalInput
	repo := &mockAnimalRepo{
		animal: &models.Animal{
			Name: "Murzik",
			Type: "Cat",
		},
		err: nil,
	}
	svc := NewAnimalService(repo)

	animal, err := svc.CreateAnimal(context.Background(), input)

	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if animal == nil {
		t.Errorf("expected animal, got nil")
	}
}
