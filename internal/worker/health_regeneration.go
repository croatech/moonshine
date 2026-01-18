package worker

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"

	"moonshine/internal/api/services"
	"moonshine/internal/repository"
)

type HpWorker struct {
	healthRegenerationService *services.HealthRegenerationService
	ticker                    *time.Ticker
}

func NewHpWorker(db *sqlx.DB, interval time.Duration) *HpWorker {
	userRepo := repository.NewUserRepository(db)
	healthRegenerationService := services.NewHealthRegenerationService(userRepo)

	return &HpWorker{
		healthRegenerationService: healthRegenerationService,
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
	if err := w.healthRegenerationService.RegenerateAllUsers(1.0); err != nil {
	}
}
