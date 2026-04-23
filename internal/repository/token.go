package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sd0hni-psina/happytail/internal/models"
)

type RefreshTokenRepository struct {
	pool *pgxpool.Pool
}

func NewRefreshTokenRepository(pool *pgxpool.Pool) *RefreshTokenRepository {
	return &RefreshTokenRepository{pool: pool}
}

func (r *RefreshTokenRepository) Create(ctx context.Context, userID int, token string, expiresAt time.Time) error {
	query := `INSERT INTO refresh_tokens (user_id, token, expires_at) VALUES ($1, $2, $3)`
	_, err := r.pool.Exec(ctx, query, userID, token, expiresAt)
	return err
}

func (r *RefreshTokenRepository) GetByToken(ctx context.Context, token string) (*models.RefreshToken, error) {
	query := `SELECT id, user_id, token, expires_at, revoked, created_at FROM refresh_tokens WHERE token = $1`

	rt := &models.RefreshToken{}
	err := r.pool.QueryRow(ctx, query, token).Scan(&rt.ID, &rt.UserID, &rt.Token, &rt.ExpiresAt, &rt.Revoked, &rt.CreatedAt)
	if err != nil {
		return nil, err
	}
	return rt, nil
}

func (r *RefreshTokenRepository) Revoke(ctx context.Context, token string) error {
	query := `UPDATE refresh_tokens SET revoked = TRUE WHERE token = $1`
	_, err := r.pool.Exec(ctx, query, token)
	return err
}

func (r *RefreshTokenRepository) RevokeAllForUser(ctx context.Context, userID int) error {
	query := `UPDATE refresh_tokens SET revoked = TRUE WHERE user_id = $1`
	_, err := r.pool.Exec(ctx, query, userID)
	return err
}
