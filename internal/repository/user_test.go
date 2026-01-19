package repository

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"moonshine/internal/domain"
)

func TestUserRepository_Create(t *testing.T) {
	if testDB == nil {
		t.Skip("Test database not initialized")
	}

	repo := NewUserRepository(testDB.DB())
	locationRepo := NewLocationRepository(testDB.DB())
	ts := time.Now().UnixNano()

	location := &domain.Location{
		Name:     fmt.Sprintf("Test Location %d", ts),
		Slug:     fmt.Sprintf("test-location-%d", ts),
		Cell:     false,
		Inactive: false,
	}
	err := locationRepo.Create(location)
	require.NoError(t, err)

	user := &domain.User{
		Username:   fmt.Sprintf("testuser%d", ts),
		Email:      fmt.Sprintf("test%d@example.com", ts),
		Password:   "hashedpassword",
		LocationID: location.ID,
	}

	err = repo.Create(user)
	require.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, user.ID)
}

func TestUserRepository_FindByUsername(t *testing.T) {
	if testDB == nil {
		t.Skip("Test database not initialized")
	}

	repo := NewUserRepository(testDB.DB())
	locationRepo := NewLocationRepository(testDB.DB())
	ts := time.Now().UnixNano()

	location := &domain.Location{
		Name:     fmt.Sprintf("Test Location %d", ts),
		Slug:     fmt.Sprintf("test-location-%d", ts),
		Cell:     false,
		Inactive: false,
	}
	err := locationRepo.Create(location)
	require.NoError(t, err)

	username := fmt.Sprintf("finduser%d", ts)
	user := &domain.User{
		Username:   username,
		Email:      fmt.Sprintf("find%d@example.com", ts),
		Password:   "hashedpassword",
		LocationID: location.ID,
	}
	err = repo.Create(user)
	require.NoError(t, err)

	found, err := repo.FindByUsername(username)
	require.NoError(t, err)
	assert.Equal(t, username, found.Username)
}

func TestUserRepository_FindByID(t *testing.T) {
	if testDB == nil {
		t.Skip("Test database not initialized")
	}

	repo := NewUserRepository(testDB.DB())
	locationRepo := NewLocationRepository(testDB.DB())
	ts := time.Now().UnixNano()

	location := &domain.Location{
		Name:     fmt.Sprintf("Test Location %d", ts),
		Slug:     fmt.Sprintf("test-location-%d", ts),
		Cell:     false,
		Inactive: false,
	}
	err := locationRepo.Create(location)
	require.NoError(t, err)

	user := &domain.User{
		Username:   fmt.Sprintf("iduser%d", ts),
		Email:      fmt.Sprintf("id%d@example.com", ts),
		Password:   "hashedpassword",
		LocationID: location.ID,
	}
	err = repo.Create(user)
	require.NoError(t, err)

	found, err := repo.FindByID(user.ID)
	require.NoError(t, err)
	assert.Equal(t, user.ID, found.ID)
}

func TestUserRepository_FindByUsername_NotFound(t *testing.T) {
	if testDB == nil {
		t.Skip("Test database not initialized")
	}

	repo := NewUserRepository(testDB.DB())

	_, err := repo.FindByUsername("nonexistent")
	assert.ErrorIs(t, err, ErrUserNotFound)
}

func TestUserRepository_FindByID_NotFound(t *testing.T) {
	if testDB == nil {
		t.Skip("Test database not initialized")
	}

	repo := NewUserRepository(testDB.DB())

	_, err := repo.FindByID(uuid.New())
	assert.ErrorIs(t, err, ErrUserNotFound)
}
