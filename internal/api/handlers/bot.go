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
	botRepo         *repository.BotRepository
}

type LocationCellsResponse struct {
	Cells []locationCell `json:"cells"`
}

type LocationBotsResponse struct {
	Bots []locationBot `json:"bots"`
}

type MoveToCellResponse struct {
	Message    string `json:"message"`
	PathLength int    `json:"path_length"`
}

type locationCell struct {
	ID       string
	Slug     string
	Name     string
	Image    string
	Inactive bool
}

type locationBot struct {
	ID    string
	name  string
	level uint8
}

var (
	errLocationNotFound       = map[string]string{"error": "location not found"}
	errLocationsNotConnected  = map[string]string{"error": "locations not connected"}
	errLocationSlugIsRequired = map[string]string{"error": "location slug is required"}
	errCellSlugIsRequired     = map[string]string{"error": "cell slug is required"}
)

func NewLocationHandler(db *sqlx.DB) *LocationHandler {
	locationRepo := repository.NewLocationRepository(db)
	userRepo := repository.NewUserRepository(db)
	movementRepo := repository.NewMovementRepository(db)
	botRepo := repository.NewBotRepository(db)
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
		botRepo:         botRepo,
	}
}

func (h *LocationHandler) MoveToLocation(c echo.Context) error {
	locationSlug := c.Param("slug")
	if locationSlug == "" {
		return c.JSON(http.StatusBadRequest, errLocationSlugIsRequired)
	}

	userID, err := middleware.GetUserIDFromContext(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusUnauthorized, "")
	}

	err = h.locationService.MoveToLocation(c.Request().Context(), userID, locationSlug)
	if err != nil {
		switch err {
		case services.ErrLocationNotConnected:
			return c.JSON(http.StatusBadRequest, errLocationsNotConnected)
		case repository.ErrLocationNotFound:
			return c.JSON(http.StatusNotFound, errLocationNotFound)
		default:
			return c.JSON(http.StatusInternalServerError, "")
		}
	}

	return c.JSON(http.StatusOK, nil)
}

func (h *LocationHandler) MoveToCell(c echo.Context) error {
	cellSlug := c.Param("cell_slug")
	if cellSlug == "" {
		return c.JSON(http.StatusBadRequest, errCellSlugIsRequired)
	}

	userID, err := middleware.GetUserIDFromContext(c.Request().Context())
	user, err := h.userRepo.FindByID(userID)

	if err != nil {
		return c.JSON(http.StatusUnauthorized, nil)
	}

	currentLocation, err := h.locationRepo.FindByID(user.LocationID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, errLocationNotFound)
	}

	if currentLocation.Slug == cellSlug {
		return c.JSON(http.StatusOK, nil)
	}

	path, err := h.locationService.FindShortestPath(currentLocation.Slug, cellSlug)
	if err != nil {
		switch err {
		case services.ErrLocationNotConnected:
			return c.JSON(http.StatusBadRequest, errLocationsNotConnected)
		case repository.ErrLocationNotFound:
			return c.JSON(http.StatusNotFound, errLocationNotFound)
		default:
			return c.JSON(http.StatusInternalServerError, nil)
		}
	}

	if err := h.locationService.StartCellMovement(userID, path); err != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}

	return c.JSON(http.StatusOK, &MoveToCellResponse{
		Message:    "movement started",
		PathLength: len(path),
	})
}

func (h *LocationHandler) GetLocationCells(c echo.Context) error {
	locationSlug := c.Param("slug")
	if locationSlug == "" {
		return c.JSON(http.StatusBadRequest, errLocationSlugIsRequired)
	}

	locationRepo := repository.NewLocationRepository(h.db)
	location, err := locationRepo.FindBySlug(locationSlug)
	if err != nil {
		if errors.Is(err, repository.ErrLocationNotFound) {
			return c.JSON(http.StatusNotFound, errLocationNotFound)
		}
		return c.JSON(http.StatusInternalServerError, nil)
	}

	cells, err := locationRepo.FindCellsByLocationID(location.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, nil)
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

func (h *LocationHandler) GetCurrentLocation(c echo.Context) error {
	userID, err := middleware.GetUserIDFromContext(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusUnauthorized, "")
	}
	user, _ := h.userRepo.FindByID(userID)
	bots, err := h.botRepo.FindBotsByLocationID(user.LocationID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "")
	}

	botsList := make([]locationBot, len(bots))
	for i, bot := range bots {
		botsList[i] = locationBot{
			ID:    bot.ID.String(),
			level: bot.Level,
			name:  bot.Name,
		}
	}

	return c.JSON(http.StatusOK, &LocationBotsResponse{
		Bots: botsList,
	})
}
