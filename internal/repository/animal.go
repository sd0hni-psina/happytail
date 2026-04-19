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

	if filter.Type != nil {
		conditions = append(conditions, fmt.Sprintf("animal_type = $%d", len(args)+1))
		args = append(args, *filter.Type)
	}

	if filter.Breed != nil {
		conditions = append(conditions, fmt.Sprintf("breed = $%d", len(args)+1))
		args = append(args, *filter.Breed)
	}

	if filter.Color != nil {
		conditions = append(conditions, fmt.Sprintf("color = $%d", len(args)+1))
		args = append(args, *filter.Color)
	}

	if filter.IsVaccinated != nil {
		conditions = append(conditions, fmt.Sprintf("is_vaccinated = $%d", len(args)+1))
		args = append(args, *filter.IsVaccinated)
	}

	if filter.HasVetPassport != nil {
		conditions = append(conditions, fmt.Sprintf("has_vet_passport = $%d", len(args)+1))
		args = append(args, *filter.HasVetPassport)
	}

	if filter.Status != nil {
		conditions = append(conditions, fmt.Sprintf("status = $%d", len(args)+1))
		args = append(args, *filter.Status)
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	var total int
	err := r.pool.QueryRow(ctx, "SELECT COUNT(*) FROM animals "+whereClause, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query := `SELECT id, animal_type, name, age, breed, color, 
       is_vaccinated, has_vet_passport, description, 
       shelter_id, status, share_count, created_at
	   FROM animals ` + whereClause + ` ORDER BY created_at DESC LIMIT $` + fmt.Sprint(len(args)+1) + ` OFFSET $` + fmt.Sprint(len(args)+2)

	args = append(args, limit, offset)
	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}

	defer rows.Close()
	var animals []models.Animal
	for rows.Next() {
		a := models.Animal{}
		err := rows.Scan(&a.ID, &a.Type, &a.Name, &a.Age, &a.Breed, &a.Color, &a.IsVaccinated, &a.HasVetPassport, &a.Description, &a.ShelterID, &a.Status, &a.ShareCount, &a.CreatedAt)
		if err != nil {
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
