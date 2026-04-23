package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	_ "github.com/sd0hni-psina/happytail/docs"
	"github.com/sd0hni-psina/happytail/internal/config"
	"github.com/sd0hni-psina/happytail/internal/handler"
	"github.com/sd0hni-psina/happytail/internal/logger"
	"github.com/sd0hni-psina/happytail/internal/middleware"
	"github.com/sd0hni-psina/happytail/internal/models"
	"github.com/sd0hni-psina/happytail/internal/repository"
	"github.com/sd0hni-psina/happytail/internal/service"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Happytail API
// @version 1.0
// @description Backend API для платформы по усыновлению домашних животных.
// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	mux := http.NewServeMux()
	cfg := config.Config{
		PostgresUser:     os.Getenv("POSTGRES_USER"),
		PostgresPassword: os.Getenv("POSTGRES_PASSWORD"),
		PostgresDB:       os.Getenv("POSTGRES_DB"),
		AppPort:          os.Getenv("APP_PORT"),
		PostgresHost:     os.Getenv("POSTGRES_HOST"),
		PostgresPort:     os.Getenv("POSTGRES_PORT"),
		JWTSecret:        os.Getenv("JWT_SECRET"),
	}
	if err := cfg.Validate(); err != nil {
		panic(err)
	}

	appEnv := os.Getenv("APP_ENV")
	log := logger.New(appEnv)

	mux.HandleFunc("/health", handler.HealthHandler)
	mux.HandleFunc("/swagger/", httpSwagger.WrapHandler)
	mux.Handle("GET /metrics", promhttp.Handler())

	pool, err := cfg.ConnectDB()
	if err != nil {
		panic(fmt.Sprintf("Unable to connect to database: %v", err))
	}
	defer pool.Close()

	authMiddleware := middleware.Auth(cfg.JWTSecret)

	// ANIMALS HANDLERS
	animalRepo := repository.NewAnimalRepository(pool)
	animalSvc := service.NewAnimalService(animalRepo)
	animalHandler := handler.NewAnimalHandler(animalSvc)
	mux.HandleFunc("GET /animals", animalHandler.GetAllAnimals)
	mux.HandleFunc("GET /animals/{id}", animalHandler.GetAnimalByID)
	mux.Handle("POST /animals", authMiddleware(http.HandlerFunc(animalHandler.CreateAnimal)))
	// SHELTERS HANDLERS
	shelterRepo := repository.NewShelterRepository(pool)
	shelterSvc := service.NewShelterService(shelterRepo)
	shelterHandler := handler.NewShelterHandler(shelterSvc)
	mux.HandleFunc("GET /shelters", shelterHandler.GetAllShelters)
	mux.HandleFunc("GET /shelters/{id}", shelterHandler.GetShelterByID)
	mux.Handle("POST /shelters", authMiddleware(http.HandlerFunc(shelterHandler.CreateShelter)))
	// USERS HANDLERS
	userRepo := repository.NewUserRepository(pool)
	tokenRepo := repository.NewRefreshTokenRepository(pool)
	userSvc := service.NewUserService(userRepo, tokenRepo, cfg.JWTSecret)
	userHandler := handler.NewUserHandler(userSvc)
	mux.HandleFunc("GET /users", userHandler.GetAllUsers)
	mux.HandleFunc("GET /users/{id}", userHandler.GetUserByID)
	mux.HandleFunc("POST /users", userHandler.CreateUser)
	mux.HandleFunc("POST /auth/login", userHandler.Login)
	mux.Handle("GET /users/me", authMiddleware(http.HandlerFunc(userHandler.GetMe)))
	mux.HandleFunc("POST /auth/refresh", userHandler.Refresh)
	mux.Handle("POST /auth/logout", authMiddleware(http.HandlerFunc(userHandler.Logout)))
	// ADOPTIONS HANDLERS
	adoptionRepo := repository.NewAdoptionRepository(pool)
	adoptionSvc := service.NewAdoptionService(adoptionRepo)
	adoptionHandler := handler.NewAdoptionHandler(adoptionSvc)
	mux.Handle("POST /adoptions", authMiddleware(http.HandlerFunc(adoptionHandler.CreateAdoption)))
	// POSTS HANDLERS
	postRepo := repository.NewPostRepository(pool)
	postSvc := service.NewPostService(postRepo)
	postHandler := handler.NewPostHandler(postSvc)
	mux.HandleFunc("GET /posts", postHandler.GetAllPost)
	mux.HandleFunc("GET /posts/{id}", postHandler.GetPostByID)
	mux.Handle("POST /posts", authMiddleware(http.HandlerFunc(postHandler.CreatePost)))
	// PHOTOS HANDLERS
	photoRepo := repository.NewAnimalPhotoRepository(pool)
	photoSvc := service.NewAnimalPhotoService(photoRepo)
	photoHandler := handler.NewAnimalPhotoHandler(photoSvc)
	mux.Handle("POST /animals/{id}/photos", authMiddleware(http.HandlerFunc(photoHandler.AddPhoto)))
	mux.Handle("DELETE /animals/{id}/photos/{photo_id}", authMiddleware(http.HandlerFunc(photoHandler.DeletePhoto)))
	mux.Handle("PATCH /animals/{id}/photos/{photo_id}/main", authMiddleware(http.HandlerFunc(photoHandler.MakeMainPhoto)))
	mux.HandleFunc("GET /animals/{id}/photos", photoHandler.GetAllPhotos)
	// ROLE HANDLERS
	roleRepo := repository.NewRoleRepository(pool)
	roleSvc := service.NewRoleService(roleRepo)
	roleHandler := handler.NewRoleHandler(roleSvc)
	mux.Handle("POST /roles", authMiddleware(middleware.RequireRole(models.RoleAdmin, roleRepo)(http.HandlerFunc(roleHandler.AppointRole))))
	mux.Handle("DELETE /roles/{id}", authMiddleware(middleware.RequireRole(models.RoleAdmin, roleRepo)(http.HandlerFunc(roleHandler.RemoveRole))))

	fmt.Println("db connected!")
	fmt.Printf("Starting server on port %s\n", cfg.AppPort)

	rateLimiter := middleware.NewRateLimite(10, 30)
	srv := &http.Server{
		Addr: ":" + cfg.AppPort,
		Handler: middleware.Logger(log)(
			middleware.Recovery(
				middleware.CORS(
					middleware.Metrics(
						rateLimiter.Middleware(mux),
					),
				),
			),
		),
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		err = srv.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}

	}()
	<-ctx.Done()

	ctxShutDown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	srv.Shutdown(ctxShutDown)
}
