package services

import (
	"context"
	"errors"
	"log"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"moonshine/internal/repository"
)

var (
	ErrLocationNotConnected = errors.New("locations are not connected")
	ErrSameLocation         = errors.New("already at this location")
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
		log.Printf("[LocationService] Failed to begin transaction: %+v", err)
		return err
	}
	defer tx.Rollback()

	// Get user's current location
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		log.Printf("[LocationService] User not found: %+v", err)
		return repository.ErrUserNotFound
	}

	// Get target location by slug
	targetLocation, err := s.locationRepo.FindBySlug(targetLocationSlug)
	if err != nil {
		log.Printf("[LocationService] Target location not found by slug %s: %+v", targetLocationSlug, err)
		return repository.ErrLocationNotFound
	}

	// Check if already at target location
	if user.LocationID == targetLocation.ID {
		log.Printf("[LocationService] User %s already at location %s", userID, targetLocationSlug)
		return ErrSameLocation
	}

	// Check if locations are connected
	var connectionCount int
	checkConnectionQuery := `
		SELECT COUNT(*) 
		FROM location_locations 
		WHERE (location_id = $1 AND connected_location_id = $2)
		   OR (location_id = $2 AND connected_location_id = $1)
	`
	err = tx.Get(&connectionCount, checkConnectionQuery, user.LocationID, targetLocation.ID)
	if err != nil {
		log.Printf("[LocationService] Failed to check location connection: %+v", err)
		return err
	}

	if connectionCount == 0 {
		log.Printf("[LocationService] Locations %s and %s are not connected", user.LocationID, targetLocation.ID)
		return ErrLocationNotConnected
	}

	// Update user's location
	updateLocationQuery := `UPDATE users SET location_id = $1 WHERE id = $2`
	_, err = tx.Exec(updateLocationQuery, targetLocation.ID, userID)
	if err != nil {
		log.Printf("[LocationService] Failed to update user location: %+v", err)
		return err
	}

	if err := tx.Commit(); err != nil {
		log.Printf("[LocationService] Failed to commit transaction: %+v", err)
		return err
	}

	log.Printf("[LocationService] User %s moved to location %s (%s)", userID, targetLocation.Slug, targetLocation.ID)
	return nil
}

