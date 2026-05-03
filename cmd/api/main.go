package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	_ "github.com/sd0hni-psina/happytail/docs"
	"github.com/sd0hni-psina/happytail/internal/cache"
	"github.com/sd0hni-psina/happytail/internal/config"
	"github.com/sd0hni-psina/happytail/internal/handler"
	"github.com/sd0hni-psina/happytail/internal/logger"
	"github.com/sd0hni-psina/happytail/internal/middleware"
	"github.com/sd0hni-psina/happytail/internal/models"
	"github.com/sd0hni-psina/happytail/internal/notifier"
	"github.com/sd0hni-psina/happytail/internal/repository"
	"github.com/sd0hni-psina/happytail/internal/service"
	"github.com/sd0hni-psina/happytail/internal/storage"
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
		SMTPHost:         os.Getenv("SMTP_HOST"),
		SMTPPort:         os.Getenv("SMTP_PORT"),
		SMTPUsername:     os.Getenv("SMTP_USERNAME"),
		SMTPPassword:     os.Getenv("SMTP_PASSWORD"),
		SMTPFrom:         os.Getenv("SMTP_FROM"),
		MinioEndpoint:    os.Getenv("MINIO_ENDPOINT"),
		MinioUser:        os.Getenv("MINIO_USER"),
		MinioPassword:    os.Getenv("MINIO_PASSWORD"),
		MinioBucket:      os.Getenv("MINIO_BUCKET"),
		MinioPublicURL:   os.Getenv("MINIO_PUBLIC_URL"),
		RedisAddr:        os.Getenv("REDIS_ADDR"),
	}
	if err := cfg.Validate(); err != nil {
		panic(err)
	}
	minioStorage, err := storage.NewMinioStorage(
		cfg.MinioEndpoint,
		cfg.MinioUser,
		cfg.MinioPassword,
		cfg.MinioBucket,
		cfg.MinioPublicURL,
	)
	if err != nil {
		panic(fmt.Sprintf("failed to connect to minio: %v", err))
	}
	emailNotifier, err := notifier.NewEmailNotifier(
		cfg.SMTPHost,
		cfg.SMTPPort,
		cfg.SMTPUsername,
		cfg.SMTPPassword,
		cfg.SMTPFrom,
	)
	if err != nil {
		panic(fmt.Sprintf("failed to create email notifier: %v", err))
	}

	appEnv := os.Getenv("APP_ENV")
	log := logger.New(appEnv)
	slog.SetDefault(log)
	redisCache, err := cache.New(cfg.RedisAddr)
	if err != nil {
		panic(fmt.Sprintf("failed to connect to redis: %v", err))
	}
	fmt.Println("redis connected!")

	mux.HandleFunc("/health", handler.HealthHandler)
	mux.HandleFunc("/swagger/", httpSwagger.WrapHandler)
	mux.Handle("GET /metrics", promhttp.Handler())

	pool, err := cfg.ConnectDB()
	if err != nil {
		panic(fmt.Sprintf("Unable to connect to database: %v", err))
	}
	defer pool.Close()
	defer redisCache.Close()

	authMiddleware := middleware.Auth(cfg.JWTSecret, redisCache)

	// REPOSITORIES
	roleRepo := repository.NewRoleRepository(pool)
	animalRepo := repository.NewAnimalRepository(pool)
	shelterRepo := repository.NewShelterRepository(pool)
	userRepo := repository.NewUserRepository(pool)
	tokenRepo := repository.NewRefreshTokenRepository(pool)
	adoptionRepo := repository.NewAdoptionRepository(pool)
	postRepo := repository.NewPostRepository(pool)
	photoRepo := repository.NewAnimalPhotoRepository(pool)
	// SERVICES
	animalSvc := service.NewAnimalService(animalRepo, redisCache)
	roleSvc := service.NewRoleService(roleRepo)
	userSvc := service.NewUserService(userRepo, tokenRepo, cfg.JWTSecret, redisCache)
	adoptionSvc := service.NewAdoptionService(adoptionRepo, userRepo, animalRepo, emailNotifier, redisCache)
	shelterSvc := service.NewShelterService(shelterRepo, redisCache)
	postSvc := service.NewPostService(postRepo, roleRepo, redisCache)
	photoSvc := service.NewAnimalPhotoService(photoRepo, minioStorage, redisCache)
	// HANDLERS
	roleHandler := handler.NewRoleHandler(roleSvc)
	animalHandler := handler.NewAnimalHandler(animalSvc)
	shelterHandler := handler.NewShelterHandler(shelterSvc, animalSvc)
	userHandler := handler.NewUserHandler(userSvc)
	adoptionHandler := handler.NewAdoptionHandler(adoptionSvc)
	postHandler := handler.NewPostHandler(postSvc)
	photoHandler := handler.NewAnimalPhotoHandler(photoSvc)
	//
	requireShelterAdminForAnimal := middleware.RequireShelterAdminForAnimal(roleRepo, animalRepo)
	requireShelterAdmin := middleware.RequireShelterAdmin(roleRepo)

	// ANIMALS HANDLERS
	mux.HandleFunc("GET /animals", animalHandler.GetAllAnimals)
	mux.HandleFunc("GET /animals/{id}", animalHandler.GetAnimalByID)
	mux.Handle("POST /animals", authMiddleware(http.HandlerFunc(animalHandler.CreateAnimal)))
	mux.Handle("PATCH /animals/{id}", authMiddleware(requireShelterAdminForAnimal(http.HandlerFunc(animalHandler.UpdateAnimal))))
	mux.Handle("DELETE /animals/{id}", authMiddleware(requireShelterAdminForAnimal(http.HandlerFunc(animalHandler.DeleteAnimal))))
	// ROLE HANDLERS
	mux.Handle("POST /roles", authMiddleware(middleware.RequireRole(models.RoleAdmin, roleRepo)(http.HandlerFunc(roleHandler.AppointRole))))
	mux.Handle("DELETE /roles/{id}", authMiddleware(middleware.RequireRole(models.RoleAdmin, roleRepo)(http.HandlerFunc(roleHandler.RemoveRole))))
	// SHELTERS HANDLERS
	mux.HandleFunc("GET /shelters/nearby", shelterHandler.FindNearby)
	mux.HandleFunc("GET /shelters", shelterHandler.GetAllShelters)
	mux.HandleFunc("GET /shelters/{id}", shelterHandler.GetShelterByID)
	mux.Handle("POST /shelters", authMiddleware(http.HandlerFunc(shelterHandler.CreateShelter)))
	mux.Handle("PATCH /shelters/{id}", authMiddleware(requireShelterAdminForAnimal(http.HandlerFunc(shelterHandler.UpdateShelter))))
	mux.HandleFunc("GET /shelters/{id}/animals", shelterHandler.GetShelterAnimals)
	mux.Handle("DELETE /shelters/{id}", authMiddleware(requireShelterAdmin(http.HandlerFunc(shelterHandler.DeleteShelter))))
	// USERS HANDLERS
	mux.HandleFunc("GET /users", userHandler.GetAllUsers)
	mux.HandleFunc("GET /users/{id}", userHandler.GetUserByID)
	mux.HandleFunc("POST /users", userHandler.CreateUser)
	mux.HandleFunc("POST /auth/login", userHandler.Login)
	mux.Handle("GET /users/me", authMiddleware(http.HandlerFunc(userHandler.GetMe)))
	mux.HandleFunc("POST /auth/refresh", userHandler.Refresh)
	mux.Handle("POST /auth/logout", authMiddleware(http.HandlerFunc(userHandler.Logout)))
	// ADOPTIONS HANDLERS
	mux.Handle("POST /adoptions", authMiddleware(http.HandlerFunc(adoptionHandler.CreateAdoption)))
	mux.Handle("GET /users/{id}/adoptions", authMiddleware(http.HandlerFunc(adoptionHandler.GetUserAdoptions)))
	// POSTS HANDLERS
	mux.HandleFunc("GET /posts", postHandler.GetAllPost)
	mux.HandleFunc("GET /posts/{id}", postHandler.GetPostByID)
	mux.Handle("POST /posts", authMiddleware(http.HandlerFunc(postHandler.CreatePost)))
	mux.Handle("PATCH /posts/{id}/status", authMiddleware(http.HandlerFunc(postHandler.UpdatePostStatus)))
	// PHOTOS HANDLERS
	mux.Handle("POST /animals/{id}/photos", authMiddleware(requireShelterAdminForAnimal(http.HandlerFunc(photoHandler.AddPhoto))))
	mux.Handle("DELETE /animals/{id}/photos/{photo_id}", authMiddleware(requireShelterAdminForAnimal(http.HandlerFunc(photoHandler.DeletePhoto))))
	mux.Handle("PATCH /animals/{id}/photos/{photo_id}/main", authMiddleware(requireShelterAdminForAnimal(http.HandlerFunc(photoHandler.MakeMainPhoto))))
	mux.HandleFunc("GET /animals/{id}/photos", photoHandler.GetAllPhotos)

	fmt.Println("db connected!")
	fmt.Printf("Starting server on port %s\n", cfg.AppPort)

	rateLimiter := middleware.NewRateLimite(10, 30)
	srv := &http.Server{
		Addr: ":" + cfg.AppPort,
		Handler: middleware.Logger(log)(
			middleware.Recovery(
				middleware.CORS(
					middleware.Metrics(
						middleware.BodyLimit(1 << 20)(
							rateLimiter.Middleware(mux),
						),
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
