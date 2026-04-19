package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sd0hni-psina/happytail/internal/models"
)

type AnimalPhotoRepository struct {
	pool *pgxpool.Pool
}

func NewAnimalPhotoRepository(pool *pgxpool.Pool) *AnimalPhotoRepository {
	return &AnimalPhotoRepository{pool: pool}
}

func (ap *AnimalPhotoRepository) Add(ctx context.Context, input models.AnimalPhotoInput) (*models.AnimalPhoto, error) {
	query := `INSERT INTO animal_photos
				(animal_id, url, is_main)
				VALUES ($1,$2,$3)
				RETURNING id, animal_id, url, is_main, created_at`
	row := ap.pool.QueryRow(ctx, query, input.AnimalID, input.URL, input.IsMain)

	p := models.AnimalPhoto{}
	err := row.Scan(&p.ID, &p.AnimalID, &p.URL, &p.IsMain, &p.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (ap *AnimalPhotoRepository) Delete(ctx context.Context, photoID int) error {
	query := `DELETE FROM animal_photos
				WHERE id = $1`

	_, err := ap.pool.Exec(ctx, query, photoID)
	if err != nil {
		return err
	}
	return nil
}

func (ap *AnimalPhotoRepository) MakeMain(ctx context.Context, animalID, photoID int) error {
	query := `UPDATE animal_photos SET is_main = FALSE WHERE animal_id = $1`
	query2 := `UPDATE animal_photos SET is_main = TRUE WHERE id = $1`
	tx, err := ap.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, query, animalID)
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, query2, photoID)
	if err != nil {
		return err
	}

	if err = tx.Commit(ctx); err != nil {
		return err
	}
	return nil
}

func (ap *AnimalPhotoRepository) GetAll(ctx context.Context, animalID int) ([]models.AnimalPhoto, error) {
	query := `SELECT 
	id, animal_id, url, is_main, created_at
	FROM animal_photos
	WHERE animal_id = $1
	ORDER BY is_main DESC, id ASC`

	rows, err := ap.pool.Query(ctx, query, animalID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var photos []models.AnimalPhoto
	for rows.Next() {
		var p models.AnimalPhoto
		err := rows.Scan(&p.ID, &p.AnimalID, &p.URL, &p.IsMain, &p.CreatedAt)
		if err != nil {
			return nil, err
		}
		photos = append(photos, p)
	}
	return photos, rows.Err()
}
