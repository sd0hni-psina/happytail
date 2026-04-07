package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sd0hni-psina/happytail/internal/models"
)

type AnimalRepository struct {
	pool *pgxpool.Pool
}

func NewAnimalRepository(pool *pgxpool.Pool) *AnimalRepository {
	return &AnimalRepository{pool: pool}
}

func (r *AnimalRepository) GetAll(ctx context.Context) ([]models.Animal, error) {
	query := `SELECT id, animal_type, name, age, breed, color, 
       is_vaccinated, has_vet_passport, description, 
       shelter_id, status, share_count, created_at
	   FROM animals`
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var animals []models.Animal
	for rows.Next() {
		a := models.Animal{}
		err := rows.Scan(&a.ID, &a.Type, &a.Name, &a.Age, &a.Breed, &a.Color, &a.IsVaccinated, &a.HasVetPassport, &a.Description, &a.ShelterID, &a.Status, &a.ShareCount, &a.CreatedAt)
		if err != nil {
			return nil, err
		}
		animals = append(animals, a)
	}
	return animals, nil
}
