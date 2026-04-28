package middleware

import (
	"context"
	"net/http"
	"strconv"

	"github.com/sd0hni-psina/happytail/internal/models"
)
// НИГДЕ НЕ ИСПОЛЬЗУЕТСЯ, НУЖНО РЕАЛИЗОВАТЬ 
type AnimalShelterGetter interface {
	GetShelterIDByAnimalID(ctx context.Context, animalID int) (*int, error)
}

type RoleChecker interface {
	HasRole(ctx context.Context, userID int, role models.RoleType, shelterID *int) (bool, error)
}

func RequireRole(role models.RoleType, repo RoleChecker) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID, ok := GetUserID(r.Context())
			if !ok {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			hasRole, err := repo.HasRole(r.Context(), userID, role, nil)
			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			if !hasRole {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func RequireShelterAdmin(repo RoleChecker) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID, ok := GetUserID(r.Context())
			if !ok {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			shelterIDStr := r.PathValue("id")
			shelterID, err := strconv.Atoi(shelterIDStr)
			if err != nil {
				http.Error(w, "Invalid shelter ID", http.StatusBadRequest)
				return
			}

			// проверяем и shelter_admin для конкретного приюта
			hasShelterRole, err := repo.HasRole(r.Context(), userID, models.RoleShelterAdmin, &shelterID)
			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			// admin может всё
			hasAdminRole, err := repo.HasRole(r.Context(), userID, models.RoleAdmin, nil)
			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			if !hasShelterRole && !hasAdminRole {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
