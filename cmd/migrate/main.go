package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println(".env not loaded, relying on environment")
	}

	var (
		flags   = flag.NewFlagSet("goose", flag.ExitOnError)
		dir     = flags.String("dir", "migrations", "directory with migration files")
		command = flags.String("command", "up", "goose command: up, down, status, create")
	)

	flags.Parse(os.Args[1:])
	args := flags.Args()

	host := getEnv("DATABASE_HOST", "localhost")
	port := getEnv("DATABASE_PORT", "5433")
	user := getEnv("DATABASE_USER", "postgres")
	password := getEnv("DATABASE_PASSWORD", "postgres")
	dbname := getEnv("DATABASE_NAME", "moonshine")
	sslmode := getEnv("DATABASE_SSL_MODE", "disable")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}

	// If command is create, handle separately
	if *command == "create" {
		if len(args) == 0 {
			log.Fatal("migration name is required for create command")
		}
		name := args[0]
		if err := goose.Create(db, *dir, name, "sql"); err != nil {
			log.Fatalf("failed to create migration: %v", err)
		}
		fmt.Printf("Created migration: %s\n", name)
		return
	}

	if err := goose.Run(*command, db, *dir, args...); err != nil {
		log.Fatalf("goose run: %v", err)
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
