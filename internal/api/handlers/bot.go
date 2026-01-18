package handlers

import (
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"

	"moonshine/internal/api/dto"
	"moonshine/internal/repository"
)

type BotHandler struct {
	db           *sqlx.DB
	locationRepo *repository.LocationRepository
	userRepo     *repository.UserRepository
	botRepo      *repository.BotRepository
}

type BotResponse struct {
	Bots []*dto.Bot `json:"bots"`
}

func NewBotHandler(db *sqlx.DB) *BotHandler {
	locationRepo := repository.NewLocationRepository(db)
	userRepo := repository.NewUserRepository(db)
	botRepo := repository.NewBotRepository(db)

	return &BotHandler{
		db:           db,
		locationRepo: locationRepo,
		userRepo:     userRepo,
		botRepo:      botRepo,
	}
}

func (h *BotHandler) GetBots(c echo.Context) error {
	locationSlug := c.Param("location_slug")
	if locationSlug == "" {
		return ErrBadRequest(c, "location slug is required")
	}

	location, err := h.locationRepo.FindBySlug(locationSlug)
	if err != nil {
		return ErrNotFound(c, "location not found")
	}

	bots, err := h.botRepo.FindBotsByLocationID(location.ID)
	if err != nil {
		return ErrInternalServerError(c)
	}

	return c.JSON(http.StatusOK, &BotResponse{
		Bots: dto.BotsFromDomain(bots),
	})
}

func (h *BotHandler) AttackBot(c echo.Context) error {
	botSlug := c.Param("slug")
	if botSlug == "" {
		return ErrBadRequest(c, "bot slug is required")
	}

	_, err := h.botRepo.FindBySlug(botSlug)
	if err != nil {
		if err == repository.ErrBotNotFound {
			return ErrNotFound(c, "bot not found")
		}
		return ErrInternalServerError(c)
	}

	return SuccessResponse(c, "attack endpoint - implementation pending")
}
