package handlers

import (
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"

	"moonshine/internal/api/dto"
	"moonshine/internal/api/services"
	"moonshine/internal/repository"
)

type AvatarHandler struct {
	avatarService *services.AvatarService
}

func NewAvatarHandler(db *sqlx.DB) *AvatarHandler {
	avatarRepo := repository.NewAvatarRepository(db)
	avatarService := services.NewAvatarService(avatarRepo)

	return &AvatarHandler{
		avatarService: avatarService,
	}
}

func (h *AvatarHandler) GetAllAvatars(c echo.Context) error {
	avatars, err := h.avatarService.GetAllAvatars(c.Request().Context())
	if err != nil {
		return ErrInternalServerError(c)
	}

	return c.JSON(http.StatusOK, dto.AvatarsFromDomain(avatars))
}







