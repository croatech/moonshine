package worker

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"

	"moonshine/internal/api/services"
	"moonshine/internal/api/ws"
	"moonshine/internal/repository"
)

type HpWorker struct {
	healthRegenerationService *services.HealthRegenerationService
	userRepo                  *repository.UserRepository
	hub                       *ws.Hub
	ticker                    *time.Ticker
}

func NewHpWorker(db *sqlx.DB, interval time.Duration) *HpWorker {
	userRepo := repository.NewUserRepository(db)
	healthRegenerationService := services.NewHealthRegenerationService(userRepo)

	return &HpWorker{
		healthRegenerationService: healthRegenerationService,
		userRepo:                  userRepo,
		hub:                       ws.GetHub(),
		ticker:                    time.NewTicker(interval),
	}
}

func (w *HpWorker) StartWorker(ctx context.Context) {
	defer w.ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-w.ticker.C:
			w.regenerateHp()
		}
	}
}

func (w *HpWorker) regenerateHp() {
	_, err := w.healthRegenerationService.RegenerateAllUsers(1.0)
	if err != nil {
		fmt.Printf("[HpWorker] Error regenerating: %v\n", err)
		return
	}

	connectedUserIDs := w.hub.GetConnectedUserIDs()
	if len(connectedUserIDs) == 0 {
		return
	}

	updates, err := w.userRepo.GetHPForUsers(connectedUserIDs)
	if err != nil {
		fmt.Printf("[HpWorker] Error getting HP: %v\n", err)
		return
	}

	for _, update := range updates {
		_ = w.hub.SendHPUpdate(update.UserID, update.CurrentHp, update.Hp)
	}
}
