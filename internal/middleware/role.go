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

func checkShelterAccess(ctx context.Context, repo RoleChecker, userID int, shelterID *int, w http.ResponseWriter) bool {
	hasShelterRole, err := repo.HasRole(ctx, userID, models.RoleShelterAdmin, shelterID)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return false
	}

	hasAdminRole, err := repo.HasRole(ctx, userID, models.RoleAdmin, nil)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return false
	}
	if !hasShelterRole && !hasAdminRole {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return false
	}
	return true
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

			if !checkShelterAccess(r.Context(), repo, userID, &shelterID, w) {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func RequireShelterAdminForAnimal(repo RoleChecker, animalGetter AnimalShelterGetter) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID, ok := GetUserID(r.Context())
			if !ok {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			animalIDStr := r.PathValue("id")
			animalID, err := strconv.Atoi(animalIDStr)
			if err != nil {
				http.Error(w, "Invalid animal ID", http.StatusBadRequest)
				return
			}

			shelterID, err := animalGetter.GetShelterIDByAnimalID(r.Context(), animalID)
			if err != nil {
				http.Error(w, "Animal not found", http.StatusNotFound)
				return
			}

			if shelterID == nil {
				hasAdmin, err := repo.HasRole(r.Context(), userID, models.RoleAdmin, nil)
				if err != nil {
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
					return
				}
				if !hasAdmin {
					http.Error(w, "Forbidden", http.StatusForbidden)
					return
				}
				next.ServeHTTP(w, r)
				return
			}
			if !checkShelterAccess(r.Context(), repo, userID, shelterID, w) {
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
