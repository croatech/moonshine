package repository

import (
	"fmt"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	_ "github.com/lib/pq"
)

var DB *gorm.DB

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
	SSLMode  string
}

func Init() error {
	config := getConfig()
	
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.User, config.Password, config.Database, config.SSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	sqlDB.SetMaxIdleConns(2)
	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)

	DB = db
	return nil
}

func Close() error {
	if DB == nil {
		return nil
	}
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func GetDB() *gorm.DB {
	return DB
}

func getConfig() Config {
	config := Config{
		Host:     getEnv("DATABASE_HOST", "localhost"),
		Port:     getEnv("DATABASE_PORT", "5433"),
		User:     getEnv("DATABASE_USER", "postgres"),
		Password: getEnv("DATABASE_PASSWORD", "postgres"),
		Database: getEnv("DATABASE_NAME", "moonshine"),
		SSLMode:  getEnv("DATABASE_SSL_MODE", "disable"),
	}
	return config
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

