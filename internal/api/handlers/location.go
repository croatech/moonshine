package handlers

import (
	"log"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"

	"moonshine/internal/api/middleware"
	"moonshine/internal/api/services"
	"moonshine/internal/repository"
)

type LocationHandler struct {
	db              *sqlx.DB
	locationService *services.LocationService
}

func NewLocationHandler(db *sqlx.DB) *LocationHandler {
	locationRepo := repository.NewLocationRepository(db)
	userRepo := repository.NewUserRepository(db)
	locationService := services.NewLocationService(db, locationRepo, userRepo)

	return &LocationHandler{
		db:              db,
		locationService: locationService,
	}
}

// MoveToLocation handles user movement between locations
// POST /api/locations/:slug/move
func (h *LocationHandler) MoveToLocation(c echo.Context) error {
	locationSlug := c.Param("slug")
	if locationSlug == "" {
		log.Printf("[LocationHandler] Bad Request: empty location slug")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "location slug is required"})
	}

	userID, err := middleware.GetUserIDFromContext(c.Request().Context())
	if err != nil {
		log.Printf("[LocationHandler] Unauthorized: user ID not found in context: %v", err)
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}

	err = h.locationService.MoveToLocation(c.Request().Context(), userID, locationSlug)
	if err != nil {
		switch err {
		case services.ErrSameLocation:
			log.Printf("[LocationHandler] Bad Request: user %s already at location %s", userID, locationSlug)
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "already at this location"})
		case services.ErrLocationNotConnected:
			log.Printf("[LocationHandler] Bad Request: locations not connected for user %s to %s", userID, locationSlug)
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "locations not connected"})
		case repository.ErrLocationNotFound:
			log.Printf("[LocationHandler] Not Found: location %s not found", locationSlug)
			return c.JSON(http.StatusNotFound, map[string]string{"error": "location not found"})
		case repository.ErrUserNotFound:
			log.Printf("[LocationHandler] Not Found: user %s not found", userID)
			return c.JSON(http.StatusNotFound, map[string]string{"error": "user not found"})
		default:
			log.Printf("[LocationHandler] Internal Server Error: failed to move to location: %+v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
		}
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "moved to location successfully"})
}

