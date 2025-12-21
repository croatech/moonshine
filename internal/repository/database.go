package repository

import (
	"fmt"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Database struct {
	db *sqlx.DB
}

func New() (*Database, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		getEnv("DATABASE_HOST", "localhost"),
		getEnv("DATABASE_PORT", "5433"),
		getEnv("DATABASE_USER", "postgres"),
		getEnv("DATABASE_PASSWORD", "postgres"),
		getEnv("DATABASE_NAME", "moonshine"),
		getEnv("DATABASE_SSL_MODE", "disable"),
	)

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	db.SetMaxIdleConns(2)
	db.SetMaxOpenConns(10)
	db.SetConnMaxLifetime(30 * time.Minute)

	return &Database{db: db}, nil
}

func (d *Database) Close() error {
	return d.db.Close()
}

func (d *Database) DB() *sqlx.DB {
	return d.db
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
