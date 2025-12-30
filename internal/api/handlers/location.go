package handlers

import (
	"errors"
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
		ID     string `json:"id"`
		Slug   string `json:"slug"`
		Name   string `json:"name"`
		Image  string `json:"image"`
		Inactive bool `json:"inactive"`
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


