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

func (m *mockAnimalRepo) GetAll(ctx context.Context, limit, ofser int, filter models.FilterParams) ([]models.Animal, int, error) {
	return nil, 0, nil
}

func (m *mockAnimalRepo) Create(ctx context.Context, input models.CreateAnimalInput) (*models.Animal, error) {
	return m.animal, m.err
}

func (m *mockAnimalRepo) GetShelterIDByAnimalID(ctx context.Context, animalID int) (*int, error) {
	return nil, nil
}

func (m *mockAnimalRepo) Update(ctx context.Context, id int, input models.UpdateAnimalInput) (*models.Animal, error) {
	return m.animal, m.err
}

func TestGetAnimalByID_NotFound(t *testing.T) {
	repo := &mockAnimalRepo{
		animal: nil,
		err:    models.ErrNotFound,
	}
	svc := NewAnimalService(repo, nil)

	animal, err := svc.GetAnimalByID(context.Background(), 1)

	if !errors.Is(err, models.ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
	if animal != nil {
		t.Errorf("expected nil animal, got %v", animal)
	}
}

func TestGetAnimalByID_Success(t *testing.T) {
	expected := &models.Animal{ID: 42, Name: "Murzik"}
	repo := &mockAnimalRepo{animal: expected, err: nil}
	svc := NewAnimalService(repo, nil)

	animal, err := svc.GetAnimalByID(context.Background(), 42)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if animal == nil {
		t.Fatal("expected animal, got nil")
	}
	if animal.ID != expected.ID {
		t.Errorf("expected ID %d, got %d", expected.ID, animal.ID)
	}
	if animal.Name != expected.Name {
		t.Errorf("expected name %s, got %s", expected.Name, animal.Name)
	}
}

func TestCreateAnimal_Succes(t *testing.T) {
	expected := &models.Animal{ID: 1, Name: "Barsik"}
	repo := &mockAnimalRepo{animal: expected, err: nil}
	svc := NewAnimalService(repo, nil)

	input := models.CreateAnimalInput{Name: "Barsik", Type: models.AnimalTypeCat}

	animal, err := svc.CreateAnimal(context.Background(), input)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if animal == nil {
		t.Fatal("expected animal, got nil")
	}
	if animal.Name != expected.Name {
		t.Errorf("expected name %s, got %s", expected.Name, animal.Name)
	}
}

func TestCreateAnimal_RepoError(t *testing.T) {
	repoErr := errors.New("database unavailable")
	repo := &mockAnimalRepo{animal: nil, err: repoErr}
	svc := NewAnimalService(repo, nil)

	input := models.CreateAnimalInput{Name: "Barsik", Type: models.AnimalTypeCat}

	animal, err := svc.CreateAnimal(context.Background(), input)

	if err == nil {
		t.Error("expected error, got nil")
	}
	if animal != nil {
		t.Errorf("expected nil animal on error, got %v", animal)
	}
}
