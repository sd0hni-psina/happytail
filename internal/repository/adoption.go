package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
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
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	query := `SELECT status FROM animals WHERE id = $1 FOR UPDATE`
	query2 := `INSERT INTO adoptions (user_id, animal_id)
				VALUES ($1, $2)
				RETURNING id, user_id, animal_id, created_at, updated_at`
	query3 := `UPDATE animals SET status = 'adopted' WHERE id = $1`
	var status string
	err = tx.QueryRow(ctx, query, animalID).Scan(&status)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, models.ErrNotFound
		}
		return nil, err
	}

	if status != "available" {
		return nil, models.ErrNotAvailable
	}

	adoption := &models.Adoption{}
	err = tx.QueryRow(ctx, query2, userID, animalID).Scan(&adoption.ID, &adoption.UserID, &adoption.AnimalID, &adoption.CreatedAt, &adoption.UpdatedAt)
	if err != nil {
		return nil, err
	}

	_, err = tx.Exec(ctx, query3, animalID)
	if err != nil {
		return nil, err
	}
	if err = tx.Commit(ctx); err != nil {
		return nil, err
	}
	return adoption, nil
}

// Пусть пока тут лежит
// query := `INSERT INTO adoptions (user_id, animal_id) VALUES ($1, $2) RETURNING id, user_id, animal_id, created_at, updated_at`
// query2 := `UPDATE animals SET status = 'adopted' WHERE id = $1`

// tx, err := r.pool.Begin(ctx)
// if err != nil {
// 	return nil, err
// }
// defer tx.Rollback(ctx)

// row := tx.QueryRow(ctx, query, userID, animalID)
// adoption := &models.Adoption{}
// err = row.Scan(&adoption.ID, &adoption.UserID, &adoption.AnimalID, &adoption.CreatedAt, &adoption.UpdatedAt)
// if err != nil {
// 	return nil, err
// }

// _, err = tx.Exec(ctx, query2, animalID)
// if err != nil {
// 	return nil, err
// }

// if err = tx.Commit(ctx); err != nil {
// 	return nil, err
// }
// return adoption, nil
