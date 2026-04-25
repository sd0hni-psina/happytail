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
	query := `SELECT id, name, address, email, phone_number, latitude, longitude, created_at FROM shelters`
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var shelters []models.Shelter
	for rows.Next() {
		var s models.Shelter
		err := rows.Scan(&s.ID, &s.Name, &s.Address, &s.Email, &s.Phone, &s.Latitude, &s.Longitude, &s.CreatedAt)
		if err != nil {
			return nil, err
		}
		shelters = append(shelters, s)
	}
	return shelters, nil
}

func (r *ShelterRepository) GetByID(ctx context.Context, id int) (*models.Shelter, error) {
	query := `SELECT id, name, address, email, phone_number, latitude, longitude, created_at 
	          FROM shelters WHERE id = $1`

	var s models.Shelter
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&s.ID, &s.Name, &s.Address, &s.Email, &s.Phone, &s.Latitude, &s.Longitude, &s.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, models.ErrNotFound
		}
		return nil, err
	}
	return &s, nil
}

func (r *ShelterRepository) Create(ctx context.Context, input models.CreateShelterInput) (*models.Shelter, error) {
	query := `INSERT INTO shelters (name, address, email, phone_number, latitude, longitude) 
	          VALUES ($1, $2, $3, $4, $5, $6) 
	          RETURNING id, name, address, email, phone_number, latitude, longitude, created_at`

	s := models.Shelter{}
	err := r.pool.QueryRow(ctx, query,
		input.Name, input.Address, input.Email, input.Phone, input.Latitude, input.Longitude,
	).Scan(&s.ID, &s.Name, &s.Address, &s.Email, &s.Phone, &s.Latitude, &s.Longitude, &s.CreatedAt)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, models.ErrConflict
		}
		return nil, err
	}
	return &s, nil
}

func (r *ShelterRepository) FindNearby(ctx context.Context, params models.NearbyParams) ([]models.ShelterWithDistance, error) {
	// ТРИДАЦАТЬ ТРИ РАЗА ПЕРЕПРОВЕРИТЬ
	query := `
		SELECT 
			id, name, address, email, phone_number, 
			latitude, longitude, created_at,
			distance_km
		FROM (
			SELECT 
				id, name, address, email, phone_number, latitude, longitude, created_at,
				(6371 * acos(
					LEAST(1.0, 
						cos(radians($1)) * cos(radians(latitude)) *
						cos(radians(longitude) - radians($2)) +
						sin(radians($1)) * sin(radians(latitude))
					)
				)) AS distance_km
			FROM shelters
			WHERE 
				latitude IS NOT NULL 
				AND longitude IS NOT NULL
				AND latitude  BETWEEN $1 - ($3 / 111.0) AND $1 + ($3 / 111.0)
				AND longitude BETWEEN $2 - ($3 / 111.0) AND $2 + ($3 / 111.0)
		) AS shelters_with_distance
		WHERE distance_km <= $3
		ORDER BY distance_km ASC
	`

	rows, err := r.pool.Query(ctx, query, params.Latitude, params.Longitude, params.RadiusKm)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var shelters []models.ShelterWithDistance
	for rows.Next() {
		var s models.ShelterWithDistance
		err := rows.Scan(
			&s.ID, &s.Name, &s.Address, &s.Email, &s.Phone,
			&s.Latitude, &s.Longitude, &s.CreatedAt, &s.Distance,
		)
		if err != nil {
			return nil, err
		}
		shelters = append(shelters, s)
	}
	return shelters, rows.Err()
}
