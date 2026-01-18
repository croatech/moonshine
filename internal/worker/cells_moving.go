package worker

import (
	"context"
	"database/sql"
	"sync"
	"time"

	"moonshine/internal/domain"
	"moonshine/internal/repository"

	"github.com/google/uuid"
)

type CellsMovingWorker struct {
	locationRepo *repository.LocationRepository
	userRepo     *repository.UserRepository
	movementRepo *repository.MovementRepository
	interval     time.Duration
	mu           sync.Mutex
	activeUsers  map[uuid.UUID]context.CancelFunc
}

func NewCellsMovingWorker(
	locationRepo *repository.LocationRepository,
	userRepo *repository.UserRepository,
	movementRepo *repository.MovementRepository,
	interval time.Duration,
) *CellsMovingWorker {
	return &CellsMovingWorker{
		locationRepo: locationRepo,
		userRepo:     userRepo,
		movementRepo: movementRepo,
		interval:     interval,
		activeUsers:  make(map[uuid.UUID]context.CancelFunc),
	}
}

func (w *CellsMovingWorker) StartMovement(userID uuid.UUID, cellSlugs []string) error {
	activeMovement, err := w.movementRepo.FindActiveByUserID(userID)
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	if activeMovement != nil {
		return nil
	}

	movement := &domain.Movement{
		UserID: userID,
		Status: domain.MovementStatusActive,
	}
	if err := w.movementRepo.Create(movement); err != nil {
		return err
	}

	w.mu.Lock()
	if cancel, exists := w.activeUsers[userID]; exists {
		cancel()
	}

	ctx, cancel := context.WithCancel(context.Background())
	w.activeUsers[userID] = cancel
	w.mu.Unlock()

	go func() {
		defer func() {
			w.mu.Lock()
			delete(w.activeUsers, userID)
			w.mu.Unlock()

			if err := w.movementRepo.UpdateStatus(movement.ID, domain.MovementStatusFinished); err != nil {
			}
		}()

		ticker := time.NewTicker(w.interval)
		defer ticker.Stop()

		user, err := w.userRepo.FindByID(userID)
		if err != nil {
			return
		}

		prevLocationID := user.LocationID

		for _, cellSlug := range cellSlugs {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				location, _ := w.locationRepo.FindBySlug(cellSlug)
				if location == nil {
					continue
				}

				cell := &domain.MovementCell{
					MovementID: movement.ID,
					FromCellID: prevLocationID,
					ToCellID:   location.ID,
				}
				if err := w.movementRepo.CreateMovementCell(cell); err != nil {
				}

				err := w.userRepo.UpdateLocationID(userID, location.ID)
				if err != nil {
					return
				}

				prevLocationID = location.ID
			}
		}
	}()

	return nil
}
