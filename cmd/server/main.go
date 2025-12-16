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
	"moonshine/internal/repository"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println(".env not loaded, relying on environment")
	}

	if err := repository.Init(); err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}
	defer repository.Close()

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	api.SetupRoutes(e)

	addr := normalizeAddr(getEnv("HTTP_ADDR", ":8080"))

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	go func() {
		log.Printf("http server starting on %s", addr)
		if err := e.Start(addr); err != nil && err != http.ErrServerClosed {
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

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func normalizeAddr(addr string) string {
	if addr == "" {
		return addr
	}

	if addr[0] == ':' || addr[0] == '[' {
		return addr
	}

	for _, r := range addr {
		if r < '0' || r > '9' {
			return addr
		}
	}

	return ":" + addr
}
