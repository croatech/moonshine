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

type LocationCellsResponse struct {
	Cells []locationCell `json:"cells"`
}

type MoveToCellResponse struct {
	Message    string `json:"message"`
	PathLength int    `json:"path_length"`
}

type locationCell struct {
	ID       string `json:"id"`
	Slug     string `json:"slug"`
	Name     string `json:"name"`
	Image    string `json:"image"`
	Inactive bool   `json:"inactive"`
}


func NewLocationHandler(db *sqlx.DB) *LocationHandler {
	locationRepo := repository.NewLocationRepository(db)
	userRepo := repository.NewUserRepository(db)
	movingWorker := worker.NewCellsMovingWorker(locationRepo, userRepo, 5*time.Second)
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
		return ErrBadRequest(c, "location slug is required")
	}

	userID, err := middleware.GetUserIDFromContext(c.Request().Context())
	if err != nil {
		return ErrUnauthorized(c)
	}

	err = h.locationService.MoveToLocation(c.Request().Context(), userID, locationSlug)
	if err != nil {
		switch err {
		case services.ErrLocationNotConnected:
			return ErrBadRequest(c, "locations not connected")
		case repository.ErrLocationNotFound:
			return ErrNotFound(c, "location not found")
		default:
			return ErrInternalServerError(c)
		}
	}

	return c.JSON(http.StatusOK, nil)
}

func (h *LocationHandler) MoveToCell(c echo.Context) error {
	cellSlug := c.Param("cell_slug")
	if cellSlug == "" {
		return ErrBadRequest(c, "cell slug is required")
	}

	userID, err := middleware.GetUserIDFromContext(c.Request().Context())
	if err != nil {
		return ErrUnauthorized(c)
	}

	user, err := h.userRepo.FindByID(userID)
	if err != nil {
		return ErrNotFound(c, "user not found")
	}

	currentLocation, err := h.locationRepo.FindByID(user.LocationID)
	if err != nil {
		return ErrNotFound(c, "location not found")
	}

	if currentLocation.Slug == cellSlug {
		return c.JSON(http.StatusOK, nil)
	}

	path, err := h.locationService.FindShortestPath(currentLocation.Slug, cellSlug)
	if err != nil {
		switch err {
		case services.ErrLocationNotConnected:
			return ErrBadRequest(c, "locations not connected")
		case repository.ErrLocationNotFound:
			return ErrNotFound(c, "location not found")
		default:
			return ErrInternalServerError(c)
		}
	}

	if err := h.locationService.StartCellMovement(userID, path); err != nil {
		return ErrBadRequest(c, "")
	}

	return c.JSON(http.StatusOK, &MoveToCellResponse{
		Message:    "movement started",
		PathLength: len(path),
	})
}

func (h *LocationHandler) GetLocationCells(c echo.Context) error {
	locationSlug := c.Param("slug")
	if locationSlug == "" {
		return ErrBadRequest(c, "location slug is required")
	}

	locationRepo := repository.NewLocationRepository(h.db)
	location, err := locationRepo.FindBySlug(locationSlug)
	if err != nil {
		if errors.Is(err, repository.ErrLocationNotFound) {
			return ErrNotFound(c, "location not found")
		}
		return ErrInternalServerError(c)
	}

	cells, err := locationRepo.FindCellsByLocationID(location.ID)
	if err != nil {
		return ErrInternalServerError(c)
	}

	cellsList := make([]locationCell, len(cells))
	for i, cell := range cells {
		cellsList[i] = locationCell{
			ID:       cell.ID.String(),
			Slug:     cell.Slug,
			Name:     cell.Name,
			Image:    cell.Image,
			Inactive: cell.Inactive,
		}
	}

	return c.JSON(http.StatusOK, &LocationCellsResponse{
		Cells: cellsList,
	})
}
