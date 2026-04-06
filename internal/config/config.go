package config

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
	PostgresUser     string
	PostgresPassword string
	PostgresDB       string
	AppPort          string
	PostgresHost     string
	PostgresPort     string
}

func (c *Config) DatabaseURL() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", c.PostgresUser, c.PostgresPassword, c.PostgresHost, c.PostgresPort, c.PostgresDB)
}

func (c *Config) ConnectDB() (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(context.Background(), c.DatabaseURL())
	if err != nil {
		return nil, err
	}
	return pool, nil
}
