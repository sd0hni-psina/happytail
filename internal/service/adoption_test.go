package service

import (
	"context"

	"github.com/sd0hni-psina/happytail/internal/models"
)

type mockAdoptedRepo struct {
	adopted *models.Adoption
	err     error
}

func (m *mockAdoptedRepo) Create(ctx context.Context, userID, animalID int) (*models.Adoption, error) {
	return m.adopted, m.err
}
