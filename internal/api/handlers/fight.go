package handlers

import (
	"moonshine/internal/api/dto"
	"net/http"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"

	"moonshine/internal/api/middleware"
	"moonshine/internal/api/services"
	"moonshine/internal/domain"
	"moonshine/internal/repository"
)

type FightHandler struct {
	fightService *services.FightService
	locationRepo *repository.LocationRepository
}

func NewFightHandler(db *sqlx.DB) *FightHandler {
	fightService := services.NewFightService(db)
	locationRepo := repository.NewLocationRepository(db)

	return &FightHandler{
		fightService: fightService,
		locationRepo: locationRepo,
	}
}

type GetCurrentFightResponse struct {
	User dto.User `json:"user"`
	Bot  dto.Bot  `json:"bot"`
}

// GetCurrentFight godoc
// @Summary Get current fight
// @Description Get information about current active fight
// @Tags fights
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} GetCurrentFightResponse
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/fights/current [get]
func (h *FightHandler) GetCurrentFight(c echo.Context) error {
	userID, err := middleware.GetUserIDFromContext(c.Request().Context())
	if err != nil {
		return ErrUnauthorized(c)
	}

	result, err := h.fightService.GetCurrentFight(c.Request().Context(), userID)
	if err != nil {
		if err == services.ErrNoActiveFight {
			return ErrNotFound(c, "no active fight")
		}
		if err == services.ErrUserNotFound {
			return ErrNotFound(c, "user not found")
		}
		if err == services.ErrBotNotFound {
			return ErrNotFound(c, "bot not found")
		}
		return ErrInternalServerError(c)
	}

	var location *domain.Location
	if result.User.LocationID != uuid.Nil {
		location, _ = h.locationRepo.FindByID(result.User.LocationID)
	}

	return c.JSON(http.StatusOK, &GetCurrentFightResponse{
		User: *dto.UserFromDomain(result.User, location, true),
		Bot:  *dto.BotFromDomain(result.Bot),
	})
}

type HitRequest struct {
	Attack   string `json:"attack" validate:"required"`
	Defense  string `json:"defense" validate:"required"`
}

// Hit godoc
// @Summary Hit in fight
// @Description Perform a hit in current fight
// @Tags fights
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body HitRequest true "Hit request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/fights/current/hit [post]
func (h *FightHandler) Hit(c echo.Context) error {
	userID, err := middleware.GetUserIDFromContext(c.Request().Context())
	if err != nil {
		return ErrUnauthorized(c)
	}

	var req HitRequest
	if err := c.Bind(&req); err != nil {
		return ErrBadRequest(c, "invalid request")
	}

	if err := c.Validate(&req); err != nil {
		return ErrBadRequest(c, err.Error())
	}

	err = h.fightService.Hit(c.Request().Context(), userID, req.Attack, req.Defense)
	if err != nil {
		if err == services.ErrNoActiveFight {
			return ErrNotFound(c, "no active fight")
		}
		if err == services.ErrInvalidBodyPart {
			return ErrBadRequest(c, "invalid body part")
		}
		return ErrInternalServerError(c)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
	})
}
