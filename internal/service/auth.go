package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sd0hni-psina/happytail/internal/models"
	"golang.org/x/crypto/bcrypt"
)

func generateAccessToken(user *models.User, secret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(15 * time.Minute).Unix(),
	})
	return token.SignedString([]byte(secret))
}

func generateRefreshToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func (s *UserService) Login(ctx context.Context, email, password string) (*models.AuthResponse, error) {
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil || user == nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, fmt.Errorf("Invalid credentials")
	}

	accessToken, err := generateAccessToken(user, s.jwtSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token")
	}

	refreshToken, err := generateRefreshToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token")
	}

	expiresAt := time.Now().Add(30 * 24 * time.Hour)
	if err = s.tokenRepo.Create(ctx, user.ID, refreshToken, expiresAt); err != nil {
		return nil, fmt.Errorf("failed to save refresh token")
	}
	return &models.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *UserService) Refresh(ctx context.Context, refreshToken string) (*models.AuthResponse, error) {
	rt, err := s.tokenRepo.GetByToken(ctx, refreshToken)
	if err != nil {
		return nil, fmt.Errorf("invalid refresher token")
	}
	if rt.Revoked {
		return nil, fmt.Errorf("refresh token revoked")
	}
	if time.Now().After(rt.ExpiresAt) {
		return nil, fmt.Errorf("refresh token expired")
	}

	user, err := s.repo.GetByID(ctx, rt.UserID)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}
	if err = s.tokenRepo.Revoke(ctx, refreshToken); err != nil {
		return nil, fmt.Errorf("failed to revoke token")
	}

	accessToken, err := generateAccessToken(&models.User{ID: user.ID}, s.jwtSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token")
	}

	newRefreshToken, err := generateRefreshToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token")
	}

	expireAt := time.Now().Add(30 * 24 * time.Hour)
	if err = s.tokenRepo.Create(ctx, user.ID, newRefreshToken, expireAt); err != nil {
		return nil, fmt.Errorf("failed to save refresh token")
	}

	return &models.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

func (s *UserService) Logout(ctx context.Context, accessToken, refreshToken string) error {
	if err := s.tokenRepo.Revoke(ctx, refreshToken); err != nil {
		return fmt.Errorf("failed to revoke refresh token: %w", err)
	}

	token, _, err := new(jwt.Parser).ParseUnverified(accessToken, jwt.MapClaims{})
	if err != nil {
		return nil
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil
	}

	expRaw, ok := claims["exp"]
	if !ok {
		return nil
	}
	expFloat, ok := expRaw.(float64)
	if !ok {
		return nil
	}

	expTime := time.Unix(int64(expFloat), 0)
	ttl := time.Until(expTime)
	if ttl <= 0 {
		return nil
	}

	blacklistKey := "blacklist:access:" + accessToken
	if s.cache != nil {
		if err := s.cache.Set(ctx, blacklistKey, "1", ttl); err != nil {
			return nil
		}
	}

	return nil
}
