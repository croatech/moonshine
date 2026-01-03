package handlers

import (
	"log"
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
	log.Printf("[AvatarHandler] Fetching all avatars")
	avatars, err := h.avatarService.GetAllAvatars(c.Request().Context())
	if err != nil {
		log.Printf("[AvatarHandler] Error fetching avatars: %+v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
	}

	log.Printf("[AvatarHandler] Found %d avatars", len(avatars))
	return c.JSON(http.StatusOK, dto.AvatarsFromDomain(avatars))
}







