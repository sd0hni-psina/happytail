package service

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sd0hni-psina/happytail/internal/models"
	"golang.org/x/crypto/bcrypt"
)

func generateToken(user *models.User, secret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (s *UserService) Login(ctx context.Context, email, password string) (string, error) {
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return "", fmt.Errorf("invalid credentials")
	}
	if user == nil {
		return "", fmt.Errorf("invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", fmt.Errorf("invalid credentials")
	}

	token, err := generateToken(user, s.jwtSecret)
	if err != nil {
		return "", fmt.Errorf("failed to generate token")
	}
	return token, nil
}
