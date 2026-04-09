package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sd0hni-psina/happytail/internal/models"
)

type AdoptionRepository struct {
	pool *pgxpool.Pool
}

func NewAdoptionRepository(pool *pgxpool.Pool) *AdoptionRepository {
	return &AdoptionRepository{pool: pool}
}

func (r *AdoptionRepository) Create(ctx context.Context, userID, animalID int) (*models.Adoption, error) {
	query := `INSERT INTO adoptions (user_id, animal_id) VALUES ($1, $2) RETURNING id, user_id, animal_id, created_at, updated_at`
	query2 := `UPDATE animals SET status = 'adopted' WHERE id = $1`
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	row := tx.QueryRow(ctx, query, userID, animalID)
	adoption := &models.Adoption{}
	err = row.Scan(&adoption.ID, &adoption.UserID, &adoption.AnimalID, &adoption.CreatedAt, &adoption.UpdatedAt)
	if err != nil {
		return nil, err
	}

	_, err = tx.Exec(ctx, query2, animalID)
	if err != nil {
		return nil, err
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, err
	}
	return adoption, nil
}
