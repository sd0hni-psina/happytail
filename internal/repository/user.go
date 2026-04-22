package repository

import (
	"context"
	"errors"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sd0hni-psina/happytail/internal/models"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{pool: pool}
}

func (r *UserRepository) GetAll(ctx context.Context) ([]models.UserPublic, error) {
	query := `SELECT id, full_name, email, phone_number, city, points, password_hash, created_at FROM users`
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.UserPublic
	for rows.Next() {
		u := models.UserPublic{}
		err := rows.Scan(&u.ID, &u.FullName, &u.Email, &u.PhoneNumber, &u.City, &u.Points, &u.CreatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id int) (*models.UserPublic, error) {
	query := `SELECT id, full_name, email, phone_number, city, points, password_hash, created_at 
FROM users WHERE id = $1`
	row := r.pool.QueryRow(ctx, query, id)

	u := models.UserPublic{}
	err := row.Scan(&u.ID, &u.FullName, &u.Email, &u.PhoneNumber, &u.City, &u.Points, &u.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, models.ErrNotFound
		}
		slog.Error("failed to get user by id", "error", err, "user_id", id)
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) Create(ctx context.Context, input models.CreateUserInput) (*models.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	input.Password = string(hash)
	query := `INSERT INTO users (full_name, email, phone_number, password_hash) VALUES ($1, $2, $3, $4) 
	RETURNING id, full_name, email, phone_number, city, points, password_hash, created_at`

	row := r.pool.QueryRow(ctx, query, input.FullName, input.Email, input.PhoneNumber, input.Password)

	u := models.User{}
	err = row.Scan(&u.ID, &u.FullName, &u.Email, &u.PhoneNumber, &u.City, &u.Points, &u.PasswordHash, &u.CreatedAt)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, models.ErrConflict
		}
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `SELECT id, full_name, email, phone_number, city, points, password_hash, created_at 
	FROM users WHERE email = $1`
	row := r.pool.QueryRow(ctx, query, email)

	u := models.User{}
	err := row.Scan(&u.ID, &u.FullName, &u.Email, &u.PhoneNumber, &u.City, &u.Points, &u.PasswordHash, &u.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}
