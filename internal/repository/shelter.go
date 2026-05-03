package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

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

func (r *ShelterRepository) GetAll(ctx context.Context, limit, offset int) ([]models.Shelter, int, error) {
	var total int
	err := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM shelters`).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query := `
		SELECT id, name, address, email, phone_number, latitude, longitude, created_at
		FROM shelters
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`
	rows, err := r.pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var shelters []models.Shelter
	for rows.Next() {
		var s models.Shelter
		err := rows.Scan(&s.ID, &s.Name, &s.Address, &s.Email, &s.Phone, &s.Latitude, &s.Longitude, &s.CreatedAt)
		if err != nil {
			return nil, 0, err
		}
		shelters = append(shelters, s)
	}
	return shelters, total, rows.Err()
}

func (r *ShelterRepository) GetByID(ctx context.Context, id int) (*models.Shelter, error) {
	query := `SELECT id, name, address, email, phone_number, latitude, longitude, created_at 
	          FROM shelters WHERE id = $1 AND deleted_at IS NULL`

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
				AND latitude  BETWEEN $1 - ($3 / 111.0)
				                  AND $1 + ($3 / 111.0)
				AND longitude BETWEEN $2 - ($3 / (111.0 * cos(radians($1))))
				                  AND $2 + ($3 / (111.0 * cos(radians($1))))
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

func (r *ShelterRepository) Update(ctx context.Context, id int, input models.UpdateShelterInput) (*models.Shelter, error) {
	setClauses := []string{}
	args := []any{}

	if input.Name != nil {
		args = append(args, *input.Name)
		setClauses = append(setClauses, fmt.Sprintf("name = $%d", len(args)))
	}
	if input.Address != nil {
		args = append(args, *input.Address)
		setClauses = append(setClauses, fmt.Sprintf("address = $%d", len(args)))
	}
	if input.Email != nil {
		args = append(args, *input.Email)
		setClauses = append(setClauses, fmt.Sprintf("email = $%d", len(args)))
	}
	if input.Phone != nil {
		args = append(args, *input.Phone)
		setClauses = append(setClauses, fmt.Sprintf("phone_number = $%d", len(args)))
	}
	if input.Latitude != nil {
		args = append(args, *input.Latitude)
		setClauses = append(setClauses, fmt.Sprintf("latitude = $%d", len(args)))
	}
	if input.Longitude != nil {
		args = append(args, *input.Longitude)
		setClauses = append(setClauses, fmt.Sprintf("longitude = $%d", len(args)))
	}

	if len(setClauses) == 0 {
		return r.GetByID(ctx, id)
	}

	args = append(args, id)
	query := fmt.Sprintf(
		`UPDATE shelters SET %s WHERE id = $%d
		 RETURNING id, name, address, email, phone_number, latitude, longitude, created_at`,
		strings.Join(setClauses, ", "), len(args),
	)

	s := models.Shelter{}
	err := r.pool.QueryRow(ctx, query, args...).Scan(
		&s.ID, &s.Name, &s.Address, &s.Email, &s.Phone, &s.Latitude, &s.Longitude, &s.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, models.ErrNotFound
		}
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, models.ErrConflict
		}
		return nil, err
	}
	return &s, nil
}

func (r *ShelterRepository) Delete(ctx context.Context, id int) error {
	query := `UPDATE shelters SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL`
	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return nil
	}
	if result.RowsAffected() == 0 {
		return models.ErrNotFound
	}
	return nil
}
