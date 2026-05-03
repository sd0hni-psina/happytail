package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sd0hni-psina/happytail/internal/cache"
)

type contextKey string

const UserIDKey contextKey = "user_id"

func GetUserID(ctx context.Context) (int, bool) {
	userID, ok := ctx.Value(UserIDKey).(int)
	return userID, ok
}

func Auth(secret string, cache *cache.Cache) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			if tokenString == authHeader {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return []byte(secret), nil
			})

			if err != nil || !token.Valid {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			blacklistKey := "blacklist:access:" + tokenString
			if cache != nil {
				exists, err := cache.Exists(r.Context(), blacklistKey)
				if err == nil && exists {
					http.Error(w, "Unauthorized", http.StatusUnauthorized)
					return
				}
			}
			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			rawID, exists := claims["user_id"]
			if !exists {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			floatID, ok := rawID.(float64)
			if !ok || floatID <= 0 {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			ctx := context.WithValue(r.Context(), UserIDKey, int(floatID))
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// Хранить отозованные access token в Редис.
