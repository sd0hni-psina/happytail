package service

import (
	"context"
	"errors"
	"testing"

	"github.com/sd0hni-psina/happytail/internal/models"
)

type mockShelterRepo struct {
	shelter  *models.Shelter
	shelters []models.Shelter
	nearby   []models.ShelterWithDistance
	err      error

	lastNearbyParams models.NearbyParams
}

func (m *mockShelterRepo) GetAll(ctx context.Context, limit, ofset int) ([]models.Shelter, int, error) {
	return m.shelters, len(m.shelters), m.err
}

func (m *mockShelterRepo) GetByID(ctx context.Context, id int) (*models.Shelter, error) {
	return m.shelter, m.err
}

func (m *mockShelterRepo) Create(ctx context.Context, input models.CreateShelterInput) (*models.Shelter, error) {
	return m.shelter, m.err
}

func (m *mockShelterRepo) Update(ctx context.Context, id int, input models.UpdateShelterInput) (*models.Shelter, error) {
	return m.shelter, m.err
}

func (m *mockShelterRepo) Delete(ctx context.Context, id int) error {
	return m.err
}

func (m *mockShelterRepo) FindNearby(ctx context.Context, params models.NearbyParams) ([]models.ShelterWithDistance, error) {
	m.lastNearbyParams = params
	return m.nearby, m.err
}

func TestGetShelterByID_Success(t *testing.T) {
	expected := &models.Shelter{ID: 1, Name: "Приют Надежда"}
	repo := &mockShelterRepo{shelter: expected}
	svc := NewShelterService(repo, nil)

	shelter, err := svc.GetShelterByID(context.Background(), 1)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if shelter == nil {
		t.Fatal("expected shelter, got nil")
	}
	if shelter.ID != expected.ID {
		t.Errorf("expected ID %d, got %d", expected.ID, shelter.ID)
	}
	if shelter.Name != expected.Name {
		t.Errorf("expected name %s, got %s", expected.Name, shelter.Name)
	}
}

func TestGetShelterByID_NotFound(t *testing.T) {
	repo := &mockShelterRepo{err: models.ErrNotFound}
	svc := NewShelterService(repo, nil)

	shelter, err := svc.GetShelterByID(context.Background(), 999)

	if !errors.Is(err, models.ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
	if shelter != nil {
		t.Errorf("expected nil shelter, got %v", shelter)
	}
}

func TestCreateShelter_Success(t *testing.T) {
	email := "info@shelter.kz"
	expected := &models.Shelter{ID: 1, Name: "Приют Надежда", Email: &email}
	repo := &mockShelterRepo{shelter: expected}
	svc := NewShelterService(repo, nil)

	input := models.CreateShelterInput{
		Name:    "Приют Надежда",
		Address: "ул. Примерная 1",
		Email:   &email,
	}

	shelter, err := svc.CreateShelter(context.Background(), input)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if shelter == nil {
		t.Fatal("expected shelter, got nil")
	}
	if shelter.Name != expected.Name {
		t.Errorf("expected name %s, got %s", expected.Name, shelter.Name)
	}
}

func TestCreateShelter_Conflict(t *testing.T) {
	repo := &mockShelterRepo{err: models.ErrConflict}
	svc := NewShelterService(repo, nil)

	email := "duplicate@shelter.kz"
	input := models.CreateShelterInput{
		Name:    "Дубль",
		Address: "ул. Дублирующая 2",
		Email:   &email,
	}

	shelter, err := svc.CreateShelter(context.Background(), input)

	if !errors.Is(err, models.ErrConflict) {
		t.Errorf("expected ErrConflict, got %v", err)
	}
	if shelter != nil {
		t.Errorf("expected nil shelter, got %v", shelter)
	}
}

func TestCreateShelter_DBError(t *testing.T) {
	repo := &mockShelterRepo{err: errors.New("connection refused")}
	svc := NewShelterService(repo, nil)

	input := models.CreateShelterInput{Name: "Test", Address: "Test"}

	_, err := svc.CreateShelter(context.Background(), input)

	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestFindNearby_Success(t *testing.T) {
	expected := []models.ShelterWithDistance{
		{Shelter: models.Shelter{ID: 1, Name: "Ближайший"}, Distance: 2.5},
	}
	repo := &mockShelterRepo{nearby: expected}
	svc := NewShelterService(repo, nil)

	params := models.NearbyParams{Latitude: 47.1, Longitude: 51.9, RadiusKm: 10}
	result, err := svc.FindNearby(context.Background(), params)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(result) != 1 {
		t.Fatalf("expected 1 shelter, got %d", len(result))
	}
	if result[0].Name != "Ближайший" {
		t.Errorf("expected name Ближайший, got %s", result[0].Name)
	}
}

func TestFindNearby_DefaultRadius(t *testing.T) {
	repo := &mockShelterRepo{nearby: []models.ShelterWithDistance{}}
	svc := NewShelterService(repo, nil)

	params := models.NearbyParams{Latitude: 47.1, Longitude: 51.9, RadiusKm: 0}
	_, err := svc.FindNearby(context.Background(), params)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if repo.lastNearbyParams.RadiusKm != 10 {
		t.Errorf("expected default radius 10, got %f", repo.lastNearbyParams.RadiusKm)
	}
}

func TestFindNearby_MaxRadius(t *testing.T) {
	repo := &mockShelterRepo{nearby: []models.ShelterWithDistance{}}
	svc := NewShelterService(repo, nil)

	params := models.NearbyParams{Latitude: 47.1, Longitude: 51.9, RadiusKm: 9999}
	_, err := svc.FindNearby(context.Background(), params)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if repo.lastNearbyParams.RadiusKm != 500 {
		t.Errorf("expected capped radius 500, got %f", repo.lastNearbyParams.RadiusKm)
	}
}

func TestFindNearby_Empty(t *testing.T) {
	repo := &mockShelterRepo{nearby: []models.ShelterWithDistance{}}
	svc := NewShelterService(repo, nil)

	params := models.NearbyParams{Latitude: 47.1, Longitude: 51.9, RadiusKm: 1}
	result, err := svc.FindNearby(context.Background(), params)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(result) != 0 {
		t.Errorf("expected empty result, got %d shelters", len(result))
	}
}

func TestFindNearby_RepoError(t *testing.T) {
	repo := &mockShelterRepo{err: errors.New("db timeout")}
	svc := NewShelterService(repo, nil)

	params := models.NearbyParams{Latitude: 47.1, Longitude: 51.9, RadiusKm: 10}
	_, err := svc.FindNearby(context.Background(), params)

	if err == nil {
		t.Error("expected error, got nil")
	}
}
