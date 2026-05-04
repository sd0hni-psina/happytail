package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sd0hni-psina/happytail/internal/models"
)

type AnimalRepository struct {
	pool *pgxpool.Pool
}

func NewAnimalRepository(pool *pgxpool.Pool) *AnimalRepository {
	return &AnimalRepository{pool: pool}
}

func (r *AnimalRepository) GetAll(ctx context.Context, limit, offset int, filter models.FilterParams) ([]models.Animal, int, error) {
	conditions := []string{}
	args := []any{}

	if filter.Name != nil {
		args = append(args, *filter.Name)
		conditions = append(conditions, fmt.Sprintf("name ILIKE $%d", len(args)))
	}
	if filter.Type != nil {
		args = append(args, *filter.Type)
		conditions = append(conditions, fmt.Sprintf("animal_type = $%d", len(args)))
	}
	if filter.Breed != nil {
		args = append(args, *filter.Breed)
		conditions = append(conditions, fmt.Sprintf("breed = $%d", len(args)))
	}
	if filter.Color != nil {
		args = append(args, *filter.Color)
		conditions = append(conditions, fmt.Sprintf("color = $%d", len(args)))
	}
	if filter.IsVaccinated != nil {
		args = append(args, *filter.IsVaccinated)
		conditions = append(conditions, fmt.Sprintf("is_vaccinated = $%d", len(args)))
	}
	if filter.HasVetPassport != nil {
		args = append(args, *filter.HasVetPassport)
		conditions = append(conditions, fmt.Sprintf("has_vet_passport = $%d", len(args)))
	}
	if filter.Status != nil {
		args = append(args, *filter.Status)
		conditions = append(conditions, fmt.Sprintf("status = $%d", len(args)))
	} else {
		args = append(args, "deleted")
		conditions = append(conditions, fmt.Sprintf("status != $%d", len(args)))
	}
	if filter.ShelterID != nil {
		args = append(args, *filter.ShelterID)
		conditions = append(conditions, fmt.Sprintf("shelter_id = $%d", len(args)))
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	var total int
	if err := r.pool.QueryRow(ctx, "SELECT COUNT(*) FROM animals "+whereClause, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := `SELECT id, animal_type, name, age, breed, color,
	       is_vaccinated, has_vet_passport, description,
	       shelter_id, status, share_count, created_at
		   FROM animals ` + whereClause +
		` ORDER BY created_at DESC LIMIT $` + fmt.Sprint(len(args)+1) +
		` OFFSET $` + fmt.Sprint(len(args)+2)

	args = append(args, limit, offset)
	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var animals []models.Animal
	for rows.Next() {
		a := models.Animal{}
		if err := rows.Scan(&a.ID, &a.Type, &a.Name, &a.Age, &a.Breed, &a.Color,
			&a.IsVaccinated, &a.HasVetPassport, &a.Description,
			&a.ShelterID, &a.Status, &a.ShareCount, &a.CreatedAt); err != nil {
			return nil, 0, err
		}
		animals = append(animals, a)
	}
	return animals, total, nil
}

func (r *AnimalRepository) GetByID(ctx context.Context, id int) (*models.Animal, error) {
	query := `SELECT id, animal_type, name, age, breed, color, 
	   is_vaccinated, has_vet_passport, description,
	   shelter_id, status, share_count, created_at
	   FROM animals WHERE id = $1`
	row := r.pool.QueryRow(ctx, query, id)

	a := models.Animal{}
	err := row.Scan(&a.ID, &a.Type, &a.Name, &a.Age, &a.Breed, &a.Color, &a.IsVaccinated, &a.HasVetPassport, &a.Description, &a.ShelterID, &a.Status, &a.ShareCount, &a.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, models.ErrNotFound
		}
		return nil, err
	}
	return &a, nil
}

func (r *AnimalRepository) Create(ctx context.Context, input models.CreateAnimalInput) (*models.Animal, error) {
	query := `INSERT INTO animals (animal_type, name, age, breed, color, is_vaccinated, has_vet_passport, description, shelter_id)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
			RETURNING id, animal_type, name, age, breed, color, is_vaccinated, has_vet_passport, description, shelter_id, status, share_count, created_at`
	row := r.pool.QueryRow(ctx, query, input.Type, input.Name, input.Age, input.Breed, input.Color, input.IsVaccinated, input.HasVetPassport, input.Description, input.ShelterID)

	a := models.Animal{}
	err := row.Scan(&a.ID, &a.Type, &a.Name, &a.Age, &a.Breed, &a.Color, &a.IsVaccinated, &a.HasVetPassport, &a.Description, &a.ShelterID, &a.Status, &a.ShareCount, &a.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *AnimalRepository) GetShelterIDByAnimalID(ctx context.Context, animalID int) (*int, error) {
	query := `SELECT shelter_id FROM animals WHERE id = $1`
	var id *int
	err := r.pool.QueryRow(ctx, query, animalID).Scan(&id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, models.ErrNotFound
		}
		return nil, err
	}
	return id, nil
}

func (r *AnimalRepository) Update(ctx context.Context, id int, input models.UpdateAnimalInput) (*models.Animal, error) {
	setClauses := []string{}
	args := []any{}

	if input.Name != nil {
		args = append(args, *input.Name)
		setClauses = append(setClauses, fmt.Sprintf("name = $%d", len(args)))
	}
	if input.Age != nil {
		args = append(args, *input.Age)
		setClauses = append(setClauses, fmt.Sprintf("age = $%d", len(args)))
	}
	if input.Breed != nil {
		args = append(args, *input.Breed)
		setClauses = append(setClauses, fmt.Sprintf("breed = $%d", len(args)))
	}
	if input.Color != nil {
		args = append(args, *input.Color)
		setClauses = append(setClauses, fmt.Sprintf("color = $%d", len(args)))
	}
	if input.IsVaccinated != nil {
		args = append(args, *input.IsVaccinated)
		setClauses = append(setClauses, fmt.Sprintf("is_vaccinated = $%d", len(args)))
	}
	if input.HasVetPassport != nil {
		args = append(args, *input.HasVetPassport)
		setClauses = append(setClauses, fmt.Sprintf("has_vet_passport = $%d", len(args)))
	}
	if input.Description != nil {
		args = append(args, *input.Description)
		setClauses = append(setClauses, fmt.Sprintf("description = $%d", len(args)))
	}
	if input.Status != nil {
		args = append(args, *input.Status)
		setClauses = append(setClauses, fmt.Sprintf("status = $%d", len(args)))
	}
	if len(setClauses) == 0 {
		return r.GetByID(ctx, id)
	}

	args = append(args, id)
	query := fmt.Sprintf(`UPDATE animals SET %s WHERE id = $%d
		RETURNING id, animal_type, name, age, breed, color,
		is_vaccinated, has_vet_passport, description,
		shelter_id, status, share_count, created_at`,
		strings.Join(setClauses, ", "), len(args),
	)

	a := models.Animal{}
	err := r.pool.QueryRow(ctx, query, args...).Scan(
		&a.ID, &a.Type, &a.Name, &a.Age, &a.Breed, &a.Color,
		&a.IsVaccinated, &a.HasVetPassport, &a.Description,
		&a.ShelterID, &a.Status, &a.ShareCount, &a.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, models.ErrNotFound
		}
		return nil, err
	}
	return &a, nil
}

func (r *AnimalRepository) Delete(ctx context.Context, id int) error {
	query := `UPDATE animals SET status = 'deleted' WHERE id = $1 AND status != 'deleted'`
	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return nil
	}
	if result.RowsAffected() == 0 {
		return models.ErrNotFound
	}
	return nil
}

func (r *AnimalRepository) IncrementShareCount(ctx context.Context, id int) (*models.Animal, error) {
	query := `
		UPDATE animals
		SET share_count = share_count + 1
		WHERE id = $1 AND status != 'deleted'
		RETURNING id, animal_type, name, age, breed, color,
			is_vaccinated, has_vet_passport, description,
			shelter_id, status, share_count, created_at
	`
	a := models.Animal{}
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&a.ID, &a.Type, &a.Name, &a.Age, &a.Breed, &a.Color,
		&a.IsVaccinated, &a.HasVetPassport, &a.Description,
		&a.ShelterID, &a.Status, &a.ShareCount, &a.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, models.ErrNotFound
		}
		return nil, err
	}
	return &a, nil
}
