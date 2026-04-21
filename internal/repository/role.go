package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sd0hni-psina/happytail/internal/models"
)

type RoleRepository struct {
	pool *pgxpool.Pool
}

func NewRoleRepository(pool *pgxpool.Pool) *RoleRepository {
	return &RoleRepository{pool: pool}
}

func (r *RoleRepository) Appoint(ctx context.Context, input models.RoleInput) (*models.Role, error) {
	query := `INSERT INTO user_roles (user_id, role, shelter_id)
              VALUES ($1, $2, $3)
              RETURNING id, user_id, role, shelter_id`
	role := &models.Role{}
	err := r.pool.QueryRow(ctx, query, input.UserID, input.RoleType, input.ShelterID).Scan(&role.ID, &role.UserID, &role.RoleType, &role.ShelterID)
	var pgErr *pgconn.PgError
	if err != nil {
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, models.ErrConflict
		}
		return nil, err
	}
	return role, nil
}

func (r *RoleRepository) Remove(ctx context.Context, roleID int) error {
	query := `DELETE FROM user_roles WHERE id = $1`

	_, err := r.pool.Exec(ctx, query, roleID)
	if err != nil {
		return err
	}
	return nil
}

func (r *RoleRepository) HasRole(ctx context.Context, userID int, role models.RoleType, shelterID *int) (bool, error) {
	query := `SELECT EXISTS(
    SELECT 1 FROM user_roles 
    WHERE user_id = $1 
    AND role = $2
    AND shelter_id IS NOT DISTINCT FROM $3
)`
	var exists bool
	err := r.pool.QueryRow(ctx, query, userID, role, shelterID).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}
