package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/sd0hni-psina/happytail/internal/config"
	"github.com/sd0hni-psina/happytail/internal/handler"
)

func main() {
	cfg := config.Config{
		PostgresUser:     os.Getenv("POSTGRES_USER"),
		PostgresPassword: os.Getenv("POSTGRES_PASSWORD"),
		PostgresDB:       os.Getenv("POSTGRES_DB"),
		AppPort:          os.Getenv("APP_PORT"),
		PostgresHost:     os.Getenv("POSTGRES_HOST"),
		PostgresPort:     os.Getenv("POSTGRES_PORT"),
	}
	if cfg.PostgresUser == "" {
		panic("POSTGRES_USER is not set")
	}
	if cfg.PostgresPassword == "" {
		panic("POSTGRES_PASSWORD is not set")
	}
	if cfg.PostgresDB == "" {
		panic("POSTGRES_DB is not set")
	}
	if cfg.AppPort == "" {
		panic("APP_PORT is not set")
	}

	http.HandleFunc("/health", handler.HealthHandler)
	pool, err := cfg.ConnectDB()
	if err != nil {
		panic(fmt.Sprintf("Unable to connect to database: %v", err))
	}
	defer pool.Close()
	_ = pool

	fmt.Println("db connected!")

	fmt.Printf("Starting server on port %s\n", cfg.AppPort)
	err = http.ListenAndServe(":"+cfg.AppPort, nil)
	if err != nil {
		panic(err)
	}
}
