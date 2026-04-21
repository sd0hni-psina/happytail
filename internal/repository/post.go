package repository

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sd0hni-psina/happytail/internal/models"
)

type PostRepository struct {
	pool *pgxpool.Pool
}

func NewPostRepository(pool *pgxpool.Pool) *PostRepository {
	return &PostRepository{pool: pool}
}

func (r *PostRepository) GetAll(ctx context.Context) ([]models.Post, error) {
	query := `
        SELECT
            id,
            user_id,
            animal_id,
            listing_type,
            price_amount,
            price_currency,
            reason,
            photo_urls,
            contact_info,
            status,
            created_at,
            updated_at
        FROM posts
        ORDER BY created_at DESC
    `
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []models.Post

	for rows.Next() {
		var p models.Post
		var amount *int64
		var currency *string
		var reason *string

		err := rows.Scan(&p.ID, &p.UserID, &p.AnimalID, &p.ListingType, &amount, &currency, &reason, &p.PhotoURLs, &p.ContactInfo, &p.Status, &p.CreatedAt, &p.UpdatedAt)
		if err != nil {
			return nil, err
		}
		if amount != nil && currency != nil {
			p.Price = &models.Money{
				Amount:   *amount,
				Currency: *currency,
			}
		}
		p.Reason = reason

		posts = append(posts, p)
	}
	return posts, rows.Err()
}

func (r *PostRepository) GetByID(ctx context.Context, id int) (*models.Post, error) {
	query := `
        SELECT
            id,
            user_id,
            animal_id,
            listing_type,
            price_amount,
            price_currency,
            reason,
            photo_urls,
            contact_info,
            status,
            created_at,
            updated_at
        FROM posts
        WHERE id = $1
    `
	var p models.Post

	var amount *int64
	var currency *string
	var reason *string

	err := r.pool.QueryRow(ctx, query, id).Scan(
		&p.ID, &p.UserID, &p.AnimalID,
		&p.ListingType, &amount, &currency,
		&reason, &p.PhotoURLs, &p.ContactInfo,
		&p.Status, &p.CreatedAt, &p.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, models.ErrNotFound
		}
		return nil, err
	}

	if amount != nil && currency != nil {
		p.Price = &models.Money{
			Amount:   *amount,
			Currency: *currency,
		}
	}

	p.Reason = reason

	return &p, nil
}

func (r *PostRepository) Create(ctx context.Context, input models.CreatePostInput) (*models.Post, error) {
	query := `
        INSERT INTO posts (
            user_id,
            animal_id,
            listing_type,
            price_amount,
            price_currency,
            reason,
            photo_urls,
            contact_info,
            status,
            created_at,
            updated_at
        )
        VALUES ($1,$2,$3,$4,$5,$6,$7,$8,'active',NOW(),NOW())
        RETURNING id, created_at, updated_at
    `

	var (
		id        int
		createdAt time.Time
		updatedAt time.Time
	)

	var amount *int64
	var currency *string

	if input.Price != nil {
		amount = &input.Price.Amount
		currency = &input.Price.Currency
	}

	err := r.pool.QueryRow(ctx, query,
		input.UserID,
		input.AnimalID,
		input.ListingType,
		amount,
		currency,
		input.Reason,
		input.PhotoURLs,
		input.ContactInfo,
	).Scan(&id, &createdAt, &updatedAt)

	if err != nil {
		return nil, err
	}

	post := &models.Post{
		ID:          id,
		UserID:      input.UserID,
		AnimalID:    input.AnimalID,
		ListingType: input.ListingType,
		Price:       input.Price,
		Reason:      input.Reason,
		PhotoURLs:   input.PhotoURLs,
		ContactInfo: input.ContactInfo,
		Status:      "active",
	}

	return post, nil
}
