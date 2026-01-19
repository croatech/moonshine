package handlers

import (
	"moonshine/internal/api/middleware"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"

	"moonshine/internal/api/dto"
	"moonshine/internal/api/services"
	"moonshine/internal/domain"
	"moonshine/internal/repository"
)

type BotHandler struct {
	botService *services.BotService
	avatarRepo *repository.AvatarRepository
}

type BotResponse struct {
	Bots []*dto.Bot `json:"bots"`
}

func NewBotHandler(db *sqlx.DB) *BotHandler {
	botService := services.NewBotService(db)

	avatarRepo := repository.NewAvatarRepository(db)

	return &BotHandler{
		botService: botService,
		avatarRepo: avatarRepo,
	}
}

func (h *BotHandler) GetBots(c echo.Context) error {
	locationSlug := c.Param("location_slug")
	if locationSlug == "" {
		return ErrBadRequest(c, "location slug is required")
	}

	bots, err := h.botService.GetBotsByLocationSlug(locationSlug)
	if err != nil {
		return ErrInternalServerError(c)
	}

	return c.JSON(http.StatusOK, &BotResponse{
		Bots: dto.BotsFromDomain(bots),
	})
}

func (h *BotHandler) Attack(c echo.Context) error {
	botSlug := c.Param("slug")
	if botSlug == "" {
		return ErrBadRequest(c, "bot slug is required")
	}

	userID, err := middleware.GetUserIDFromContext(c.Request().Context())
	if err != nil {
		return ErrUnauthorized(c)
	}

	result, err := h.botService.Attack(c.Request().Context(), botSlug, userID)
	if err != nil {
		if err == repository.ErrBotNotFound {
			return ErrNotFound(c, "bot not found")
		}
		if err == repository.ErrUserNotFound {
			return ErrNotFound(c, "user not found")
		}
		return ErrBadRequest(c, err.Error())
	}

	var avatar *domain.Avatar
	if result.User.AvatarID != nil {
		avatar, _ = h.avatarRepo.FindByID(*result.User.AvatarID)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"user": dto.UserFromDomain(result.User, avatar, nil, nil),
		"bot":  dto.BotFromDomain(result.Bot),
	})
}
