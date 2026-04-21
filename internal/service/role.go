package service

import (
	"context"

	"github.com/sd0hni-psina/happytail/internal/models"
)

type RoleService struct {
	repo RoleRepository
}

func NewRoleService(repo RoleRepository) *RoleService {
	return &RoleService{repo: repo}
}

func (rs *RoleService) AppointRole(ctx context.Context, input models.RoleInput) (*models.Role, error) {
	return rs.repo.Appoint(ctx, input)
}

func (rs *RoleService) RemoveRole(ctx context.Context, roleID int) error {
	return rs.repo.Remove(ctx, roleID)
}

func (rs *RoleService) HasRole(ctx context.Context, userID int, role models.RoleType, shelterID *int) (bool, error) {
	return rs.repo.HasRole(ctx, userID, role, shelterID)
}
