package handlers

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"

	"moonshine/internal/api/middleware"
	"moonshine/internal/api/services"
	"moonshine/internal/repository"
	"moonshine/internal/worker"
)

type LocationHandler struct {
	db              *sqlx.DB
	locationService *services.LocationService
	locationRepo    *repository.LocationRepository
	userRepo        *repository.UserRepository
}

func NewLocationHandler(db *sqlx.DB) *LocationHandler {
	locationRepo := repository.NewLocationRepository(db)
	userRepo := repository.NewUserRepository(db)
	movementRepo := repository.NewMovementRepository(db)
	movingWorker := worker.NewCellsMovingWorker(locationRepo, userRepo, movementRepo, 5*time.Second)
	locationService, err := services.NewLocationService(db, locationRepo, userRepo, movingWorker)
	if err != nil {
		log.Fatalf("Failed to create LocationService: %v", err)
	}

	return &LocationHandler{
		db:              db,
		locationService: locationService,
		locationRepo:    locationRepo,
		userRepo:        userRepo,
	}
}

func (h *LocationHandler) MoveToLocation(c echo.Context) error {
	locationSlug := c.Param("slug")
	if locationSlug == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "location slug is required"})
	}

	userID, err := middleware.GetUserIDFromContext(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}

	err = h.locationService.MoveToLocation(c.Request().Context(), userID, locationSlug)
	if err != nil {
		switch err {
		case services.ErrLocationNotConnected:
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "locations not connected"})
		case repository.ErrLocationNotFound:
			return c.JSON(http.StatusNotFound, map[string]string{"error": "location not found"})
		default:
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
		}
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "moved to location successfully"})
}

func (h *LocationHandler) MoveToCell(c echo.Context) error {
	cellSlug := c.Param("cell_slug")
	if cellSlug == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "cell slug is required"})
	}

	userID, err := middleware.GetUserIDFromContext(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}

	user, err := h.userRepo.FindByID(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to get user"})
	}

	currentLocation, err := h.locationRepo.FindByID(user.LocationID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to get current location"})
	}

	if currentLocation.Slug == cellSlug {
		return c.JSON(http.StatusOK, map[string]string{"message": "already at destination"})
	}

	path, err := h.locationService.FindShortestPath(currentLocation.Slug, cellSlug)
	if err != nil {
		switch err {
		case services.ErrLocationNotConnected:
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "locations not connected"})
		case repository.ErrLocationNotFound:
			return c.JSON(http.StatusNotFound, map[string]string{"error": "location not found"})
		default:
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
		}
	}

	if err := h.locationService.StartCellMovement(userID, path); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":     "movement started",
		"path_length": len(path),
	})
}

func (h *LocationHandler) GetLocationCells(c echo.Context) error {
	locationSlug := c.Param("slug")
	if locationSlug == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "location slug is required"})
	}

	locationRepo := repository.NewLocationRepository(h.db)
	location, err := locationRepo.FindBySlug(locationSlug)
	if err != nil {
		if errors.Is(err, repository.ErrLocationNotFound) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "location not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
	}

	cells, err := locationRepo.FindCellsByLocationID(location.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
	}

	type LocationCell struct {
		ID       string `json:"id"`
		Slug     string `json:"slug"`
		Name     string `json:"name"`
		Image    string `json:"image"`
		Inactive bool   `json:"inactive"`
	}

	cellsResponse := make([]LocationCell, len(cells))
	for i, cell := range cells {
		cellsResponse[i] = LocationCell{
			ID:       cell.ID.String(),
			Slug:     cell.Slug,
			Name:     cell.Name,
			Image:    cell.Image,
			Inactive: cell.Inactive,
		}
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"cells": cellsResponse,
	})
}
