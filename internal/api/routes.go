package api

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/jmoiron/sqlx"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"

	"moonshine/internal/api/handlers"
	jwtMiddleware "moonshine/internal/api/middleware"
)

func SetupRoutes(e *echo.Echo, db *sqlx.DB, isProduction bool) {
	e.GET("/health", healthCheck)

	// Add middleware to disable cache for static files in development (must be before Static)
	if !isProduction {
		e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
			return func(c echo.Context) error {
				if strings.HasPrefix(c.Request().URL.Path, "/assets") {
					c.Response().Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
					c.Response().Header().Set("Pragma", "no-cache")
					c.Response().Header().Set("Expires", "0")
				}
				return next(c)
			}
		})
	}

	// Serve static assets
	// Try multiple possible paths relative to common project root locations
	var assetsPath string
	possiblePaths := []string{
		"frontend/assets",    // From project root
		"../frontend/assets", // From cmd/server
		filepath.Join(filepath.Dir(os.Args[0]), "../frontend/assets"), // From compiled binary
	}

	for _, path := range possiblePaths {
		absPath, err := filepath.Abs(path)
		if err == nil {
			if _, err := os.Stat(filepath.Join(absPath, "images")); err == nil {
				assetsPath = absPath
				break
			}
		}
	}

	if assetsPath == "" {
		// Fallback: try to find project root by looking for go.mod
		wd, _ := os.Getwd()
		for {
			if _, err := os.Stat(filepath.Join(wd, "go.mod")); err == nil {
				// Found project root
				assetsPath, _ = filepath.Abs(filepath.Join(wd, "frontend/assets"))
				if _, err := os.Stat(assetsPath); err == nil {
					break
				}
			}
			parent := filepath.Dir(wd)
			if parent == wd {
				break
			}
			wd = parent
		}
	}

	if assetsPath != "" {
		e.Static("/assets", assetsPath)
		log.Printf("[Static Assets] Serving /assets from: %s", assetsPath)
	} else {
		e.Static("/assets", "frontend/assets")
		log.Printf("[Static Assets] Using default relative path: frontend/assets")
	}

	// Validator for request validation
	e.Validator = NewValidator()

	// Auth handlers (public routes)
	authHandler := handlers.NewAuthHandler(db)
	authGroup := e.Group("/api/auth")
	authGroup.POST("/signup", authHandler.SignUp)
	authGroup.POST("/signin", authHandler.SignIn)

	// Protected routes with JWT authentication
	jwtConfig := echojwt.Config{
		SigningKey: []byte(os.Getenv("JWT_KEY")),
		ContextKey: "user",
		ErrorHandler: func(c echo.Context, err error) error {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		},
	}

	// Protected routes with JWT authentication
	apiGroup := e.Group("/api")
	apiGroup.Use(echojwt.WithConfig(jwtConfig))
	apiGroup.Use(jwtMiddleware.ExtractUserIDFromJWT())

	// User handlers
	userHandler := handlers.NewUserHandler(db)
	apiGroup.GET("/user/me", userHandler.GetCurrentUser)
	apiGroup.PUT("/user/me", userHandler.UpdateCurrentUser)
	apiGroup.GET("/users/me/inventory", userHandler.GetUserInventory)
	apiGroup.GET("/users/me/equipped", userHandler.GetUserEquippedItems)

	// Avatar handlers
	avatarHandler := handlers.NewAvatarHandler(db)
	apiGroup.GET("/avatars", avatarHandler.GetAllAvatars)

	// Location handlers
	locationHandler := handlers.NewLocationHandler(db)
	apiGroup.POST("/locations/:slug/move", locationHandler.MoveToLocation)

	// Equipment item handlers (protected route)
	equipmentItemHandler := handlers.NewEquipmentItemHandler(db)
	apiGroup.GET("/equipment_items", equipmentItemHandler.GetEquipmentItems)
	apiGroup.POST("/equipment_items/take_off/:slot", equipmentItemHandler.TakeOffEquipmentItem)
	apiGroup.POST("/equipment_items/:slug/buy", equipmentItemHandler.BuyEquipmentItem)
	apiGroup.POST("/equipment_items/:slug/sell", equipmentItemHandler.SellEquipmentItem)
	apiGroup.POST("/equipment_items/:slug/take_on", equipmentItemHandler.TakeOnEquipmentItem)
}

func healthCheck(c echo.Context) error {
	return c.String(http.StatusOK, "ok")
}
