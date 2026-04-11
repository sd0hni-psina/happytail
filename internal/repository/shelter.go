package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sd0hni-psina/happytail/internal/models"
)

type ShelterRepository struct {
	pool *pgxpool.Pool
}

func NewShelterRepository(pool *pgxpool.Pool) *ShelterRepository {
	return &ShelterRepository{pool: pool}
}

func (r *ShelterRepository) GetAll(ctx context.Context) ([]models.Shelter, error) {
	query := `SELECT id, name, 
				address, email,
				phone_number, 
				created_at 
				FROM shelters`
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var shelters []models.Shelter
	for rows.Next() {
		var s models.Shelter
		err := rows.Scan(&s.ID, &s.Name, &s.Address, &s.Email, &s.Phone, &s.CreatedAt)
		if err != nil {
			return nil, err
		}
		shelters = append(shelters, s)
	}
	return shelters, nil
}

func (r *ShelterRepository) GetByID(ctx context.Context, id int) (*models.Shelter, error) {
	query := `SELECT id, name, 
				address, email,
				phone_number, 
				created_at 
				FROM shelters WHERE id = $1`
	row := r.pool.QueryRow(ctx, query, id)

	var s models.Shelter
	err := row.Scan(&s.ID, &s.Name, &s.Address, &s.Email, &s.Phone, &s.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, models.ErrNotFound
		}
		return nil, err
	}
	return &s, nil
}

func (r *ShelterRepository) Create(ctx context.Context, input models.CreateShelterInput) (*models.Shelter, error) {
	query := `INSERT INTO shelters (name, address, email, phone_number) 
				VALUES ($1, $2, $3, $4) RETURNING id, name, address, email, phone_number, created_at`
	row := r.pool.QueryRow(ctx, query, input.Name, input.Address, input.Email, input.Phone)

	s := models.Shelter{}
	err := row.Scan(&s.ID, &s.Name, &s.Address, &s.Email, &s.Phone, &s.CreatedAt)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, models.ErrConflict
		}
		return nil, err
	}
	return &s, nil
}
