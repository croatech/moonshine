package handlers

import (
	"moonshine/internal/domain"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"

	"moonshine/internal/api/middleware"
	"moonshine/internal/repository"
)

type BotHandler struct {
	db *sqlx.DB
	// botService *services.BotService
	locationRepo *repository.LocationRepository
	userRepo     *repository.UserRepository
	botRepo      *repository.BotRepository
}

type BotResponse struct {
	Bots []*domain.Bot `json:"bots"`
}

type bot struct {
	ID    string
	name  string
	level uint8
}

func NewBotHandler(db *sqlx.DB) *BotHandler {
	locationRepo := repository.NewLocationRepository(db)
	userRepo := repository.NewUserRepository(db)
	botRepo := repository.NewBotRepository(db)

	return &BotHandler{
		db: db,
		// botService: botService,
		locationRepo: locationRepo,
		userRepo:     userRepo,
		botRepo:      botRepo,
	}
}

func (h *BotHandler) GetBots(c echo.Context) error {
	userID, err := middleware.GetUserIDFromContext(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusUnauthorized, "")
	}
	user, _ := h.userRepo.FindByID(userID)
	bots, err := h.botRepo.FindBotsByLocationID(user.LocationID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, &BotResponse{
		Bots: bots,
	})
}
