package services

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"moonshine/internal/domain"
	"moonshine/internal/repository"
)

var (
	ErrLocationNotConnected = errors.New("locations are not connected")
)

type LocationService struct {
	db           *sqlx.DB
	locationRepo *repository.LocationRepository
	userRepo     *repository.UserRepository
}

func NewLocationService(
	db *sqlx.DB,
	locationRepo *repository.LocationRepository,
	userRepo *repository.UserRepository,
) *LocationService {
	return &LocationService{
		db:           db,
		locationRepo: locationRepo,
		userRepo:     userRepo,
	}
}

func (s *LocationService) MoveToLocation(ctx context.Context, userID uuid.UUID, targetLocationSlug string) error {
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return repository.ErrUserNotFound
	}

	currentLocation, err := s.locationRepo.FindByID(user.LocationID)
	if err != nil {
		return repository.ErrLocationNotFound
	}

	var targetLocation *domain.Location
	if targetLocationSlug == "wayward_pines" && currentLocation.Slug == "moonshine" {
		defaultOutDoorLocation, err := s.locationRepo.DefaultOutdoorLocation()
		if err != nil {
			return repository.ErrLocationNotFound
		}
		targetLocation = defaultOutDoorLocation
	} else {
		targetLocation, err = s.locationRepo.FindBySlug(targetLocationSlug)
		if err != nil {
			return repository.ErrLocationNotFound
		}
	}

	if user.LocationID == targetLocation.ID {
		return nil
	}

	var connectionCount int
	checkConnectionQuery := `
		SELECT COUNT(*) 
		FROM location_locations 
		WHERE (location_id = $1 AND near_location_id = $2)
		   OR (location_id = $2 AND near_location_id = $1)
	`
	err = tx.Get(&connectionCount, checkConnectionQuery, user.LocationID, targetLocation.ID)
	if err != nil {
		return err
	}

	if connectionCount == 0 {
		return ErrLocationNotConnected
	}

	updateLocationQuery := `UPDATE users SET location_id = $1 WHERE id = $2`
	_, err = tx.Exec(updateLocationQuery, targetLocation.ID, userID)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
