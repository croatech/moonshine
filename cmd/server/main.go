package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"moonshine/internal/api"
	"moonshine/internal/config"
	"moonshine/internal/repository"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println(".env not loaded, relying on environment")
	}

	cfg := config.Load()

	db, err := repository.New()
	if err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}
	defer db.Close()

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	api.SetupRoutes(e, db.DB(), cfg.IsProduction())

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	go func() {
		log.Printf("http server starting on %s", cfg.HTTPAddr)
		if err := e.Start(cfg.HTTPAddr); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server failed: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("shutting down gracefully...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("server shutdown failed: %v", err)
	}
}
