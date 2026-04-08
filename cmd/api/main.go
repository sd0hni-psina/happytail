package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/sd0hni-psina/happytail/internal/config"
	"github.com/sd0hni-psina/happytail/internal/handler"
	"github.com/sd0hni-psina/happytail/internal/middleware"
	"github.com/sd0hni-psina/happytail/internal/repository"
	"github.com/sd0hni-psina/happytail/internal/service"
)

func main() {

	mux := http.NewServeMux()

	cfg := config.Config{
		PostgresUser:     os.Getenv("POSTGRES_USER"),
		PostgresPassword: os.Getenv("POSTGRES_PASSWORD"),
		PostgresDB:       os.Getenv("POSTGRES_DB"),
		AppPort:          os.Getenv("APP_PORT"),
		PostgresHost:     os.Getenv("POSTGRES_HOST"),
		PostgresPort:     os.Getenv("POSTGRES_PORT"),
	}

	mux.HandleFunc("/health", handler.HealthHandler)

	pool, err := cfg.ConnectDB()
	if err != nil {
		panic(fmt.Sprintf("Unable to connect to database: %v", err))
	}
	defer pool.Close()

	animalRepo := repository.NewAnimalRepository(pool)
	animalSvc := service.NewAnimalService(animalRepo)
	animalHandler := handler.NewAnimalHandler(animalSvc)
	mux.HandleFunc("GET /animals", animalHandler.GetAllAnimals)

	mux.HandleFunc("GET /animals/{id}", animalHandler.GetAnimalByID)

	mux.HandleFunc("POST /animals", animalHandler.CreateAnimal)

	fmt.Println("db connected!")
	fmt.Printf("Starting server on port %s\n", cfg.AppPort)
	err = http.ListenAndServe(":"+cfg.AppPort, middleware.Logger(middleware.Recovery(mux)))
	if err != nil {
		panic(err)
	}
}
