package config

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
	PostgresUser     string
	PostgresPassword string
	PostgresDB       string
	PostgresHost     string
	PostgresPort     string

	AppPort   string
	JWTSecret string

	SMTPHost     string
	SMTPPort     string
	SMTPUsername string
	SMTPPassword string
	SMTPFrom     string

	MinioEndpoint  string
	MinioUser      string
	MinioPassword  string
	MinioBucket    string
	MinioPublicURL string
}

func (c *Config) Validate() error {
	var errs []string

	if c.PostgresUser == "" {
		errs = append(errs, "POSTGRES_USER is required")
	}
	if c.PostgresPassword == "" {
		errs = append(errs, "POSTGRES_PASSWORD is required")
	}
	if c.PostgresDB == "" {
		errs = append(errs, "POSTGRES_DB is required")
	}
	if c.PostgresHost == "" {
		errs = append(errs, "POSTGRES_HOST is required")
	}
	if c.PostgresPort == "" {
		errs = append(errs, "POSTGRES_PORT is required")
	}
	if c.AppPort == "" {
		errs = append(errs, "APP_PORT is required")
	}
	if c.JWTSecret == "" {
		errs = append(errs, "JWT_SECRET is required")
	} else if len(c.JWTSecret) < 32 {
		errs = append(errs, "JWT_SECRET must be at least 32 characters")
	}
	if c.SMTPHost == "" {
		errs = append(errs, "SMTP_HOST is required")
	}
	if c.SMTPPort == "" {
		errs = append(errs, "SMTP_PORT is required")
	} else {
		port, err := strconv.Atoi(c.SMTPPort)
		if err != nil {
			errs = append(errs, "SMTP_PORT must be a number")
		} else if port < 1 || port > 65535 {
			errs = append(errs, "SMTP_PORT must be between 1 and 65535")
		}
	}
	if c.SMTPUsername == "" {
		errs = append(errs, "SMTP_USERNAME is required")
	}
	if c.SMTPPassword == "" {
		errs = append(errs, "SMTP_PASSWORD is required")
	}
	if c.MinioEndpoint == "" {
		errs = append(errs, "MINIO_ENDPOINT is required")
	}
	if c.MinioUser == "" {
		errs = append(errs, "MINIO_USER is required")
	}
	if c.MinioPassword == "" {
		errs = append(errs, "MINIO_PASSWORD is required")
	}
	if c.MinioBucket == "" {
		errs = append(errs, "MINIO_BUCKET is required")
	}

	if len(errs) > 0 {
		return errors.New("config validation failed:\n  - " + strings.Join(errs, "\n  - "))
	}
	return nil
}

func (c *Config) DatabaseURL() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", c.PostgresUser, c.PostgresPassword, c.PostgresHost, c.PostgresPort, c.PostgresDB)
}

func (c *Config) ConnectDB() (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(c.DatabaseURL())
	if err != nil {
		return nil, err
	}

	config.MaxConns = 25
	config.MinConns = 5
	config.MaxConnLifetime = 1 * time.Hour
	config.MaxConnIdleTime = 30 * time.Minute
	config.HealthCheckPeriod = 1 * time.Minute

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, err
	}
	return pool, nil
}
