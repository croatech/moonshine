package services

import (
	"testing"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"moonshine/internal/domain"
	"moonshine/internal/repository"
)

func setupLocationTestData(db *sqlx.DB) ([]*domain.Location, error) {
	locations := []*domain.Location{
		{Model: domain.Model{ID: uuid.New()}, Name: "Cell 1", Slug: "1cell", Cell: true, Inactive: false},
		{Model: domain.Model{ID: uuid.New()}, Name: "Cell 2", Slug: "2cell", Cell: true, Inactive: false},
		{Model: domain.Model{ID: uuid.New()}, Name: "Cell 3", Slug: "3cell", Cell: true, Inactive: false},
		{Model: domain.Model{ID: uuid.New()}, Name: "Cell 4", Slug: "4cell", Cell: true, Inactive: false},
		{Model: domain.Model{ID: uuid.New()}, Name: "Cell 5", Slug: "5cell", Cell: true, Inactive: false},
	}

	locationRepo := repository.NewLocationRepository(db)
	for _, loc := range locations {
		if err := locationRepo.Create(loc); err != nil {
			return nil, err
		}
	}

	return locations, nil
}

func createConnections(db *sqlx.DB, connections [][2]uuid.UUID) error {
	for _, conn := range connections {
		query := `INSERT INTO location_locations (id, location_id, near_location_id) VALUES ($1, $2, $3)`
		_, err := db.Exec(query, uuid.New(), conn[0], conn[1])
		if err != nil {
			return err
		}
	}
	return nil
}

func TestLocationService_FindShortestPath(t *testing.T) {
	if testDB == nil {
		t.Skip("Test database not initialized")
	}

	db := testDB.DB()
	locationRepo := repository.NewLocationRepository(db)
	userRepo := repository.NewUserRepository(db)
	service := NewLocationService(db, locationRepo, userRepo)

	t.Run("successful path finding - direct connection", func(t *testing.T) {
		db.Exec("TRUNCATE TABLE location_locations CASCADE")
		db.Exec("TRUNCATE TABLE locations CASCADE")

		locations, err := setupLocationTestData(db)
		if err != nil {
			t.Fatalf("Failed to setup test data: %v", err)
		}

		err = createConnections(db, [][2]uuid.UUID{
			{locations[0].ID, locations[1].ID},
		})
		if err != nil {
			t.Fatalf("Failed to create connections: %v", err)
		}

		path, err := service.FindShortestPath(locations[0].Slug, locations[1].Slug)
		if err != nil {
			t.Fatalf("FindShortestPath returned error: %v", err)
		}

		if len(path) != 1 {
			t.Fatalf("Expected path length 1, got %d", len(path))
		}

		if path[0] != locations[1].Slug {
			t.Fatalf("Expected path[0] to be %s, got %s", locations[1].Slug, path[0])
		}
	})

	t.Run("successful path finding - multi-step path", func(t *testing.T) {
		db.Exec("TRUNCATE TABLE location_locations CASCADE")
		db.Exec("TRUNCATE TABLE locations CASCADE")

		locations, err := setupLocationTestData(db)
		if err != nil {
			t.Fatalf("Failed to setup test data: %v", err)
		}

		err = createConnections(db, [][2]uuid.UUID{
			{locations[0].ID, locations[1].ID},
			{locations[1].ID, locations[2].ID},
			{locations[2].ID, locations[3].ID},
		})
		if err != nil {
			t.Fatalf("Failed to create connections: %v", err)
		}

		path, err := service.FindShortestPath(locations[0].Slug, locations[3].Slug)
		if err != nil {
			t.Fatalf("FindShortestPath returned error: %v", err)
		}

		expectedPath := []string{locations[1].Slug, locations[2].Slug, locations[3].Slug}
		if len(path) != len(expectedPath) {
			t.Fatalf("Expected path length %d, got %d", len(expectedPath), len(path))
		}

		for i, expected := range expectedPath {
			if path[i] != expected {
				t.Fatalf("Expected path[%d] to be %s, got %s", i, expected, path[i])
			}
		}
	})

	t.Run("same start and end location", func(t *testing.T) {
		db.Exec("TRUNCATE TABLE location_locations CASCADE")
		db.Exec("TRUNCATE TABLE locations CASCADE")

		locations, err := setupLocationTestData(db)
		if err != nil {
			t.Fatalf("Failed to setup test data: %v", err)
		}

		path, err := service.FindShortestPath(locations[0].Slug, locations[0].Slug)
		if err != nil {
			t.Fatalf("FindShortestPath returned error: %v", err)
		}

		if len(path) != 1 {
			t.Fatalf("Expected path length 1, got %d", len(path))
		}

		if path[0] != locations[0].Slug {
			t.Fatalf("Expected path[0] to be %s, got %s", locations[0].Slug, path[0])
		}
	})

	t.Run("path not found - disconnected locations", func(t *testing.T) {
		db.Exec("TRUNCATE TABLE location_locations CASCADE")
		db.Exec("TRUNCATE TABLE locations CASCADE")

		locations, err := setupLocationTestData(db)
		if err != nil {
			t.Fatalf("Failed to setup test data: %v", err)
		}

		err = createConnections(db, [][2]uuid.UUID{
			{locations[0].ID, locations[1].ID},
		})
		if err != nil {
			t.Fatalf("Failed to create connections: %v", err)
		}

		_, err = service.FindShortestPath(locations[0].Slug, locations[4].Slug)
		if err != ErrLocationNotConnected {
			t.Fatalf("Expected ErrLocationNotConnected, got %v", err)
		}
	})

	t.Run("non-existent start location", func(t *testing.T) {
		db.Exec("TRUNCATE TABLE location_locations CASCADE")
		db.Exec("TRUNCATE TABLE locations CASCADE")

		locations, err := setupLocationTestData(db)
		if err != nil {
			t.Fatalf("Failed to setup test data: %v", err)
		}

		_, err = service.FindShortestPath("non_existent", locations[0].Slug)
		if err == nil {
			t.Fatal("Expected error for non-existent location, got nil")
		}
	})

	t.Run("non-existent end location", func(t *testing.T) {
		db.Exec("TRUNCATE TABLE location_locations CASCADE")
		db.Exec("TRUNCATE TABLE locations CASCADE")

		locations, err := setupLocationTestData(db)
		if err != nil {
			t.Fatalf("Failed to setup test data: %v", err)
		}

		_, err = service.FindShortestPath(locations[0].Slug, "non_existent")
		if err == nil {
			t.Fatal("Expected error for non-existent location, got nil")
		}
	})

	t.Run("path with alternative routes - shortest path", func(t *testing.T) {
		db.Exec("TRUNCATE TABLE location_locations CASCADE")
		db.Exec("TRUNCATE TABLE locations CASCADE")

		locations, err := setupLocationTestData(db)
		if err != nil {
			t.Fatalf("Failed to setup test data: %v", err)
		}

		err = createConnections(db, [][2]uuid.UUID{
			{locations[0].ID, locations[1].ID},
			{locations[1].ID, locations[4].ID},
			{locations[0].ID, locations[2].ID},
			{locations[2].ID, locations[3].ID},
			{locations[3].ID, locations[4].ID},
		})
		if err != nil {
			t.Fatalf("Failed to create connections: %v", err)
		}

		path, err := service.FindShortestPath(locations[0].Slug, locations[4].Slug)
		if err != nil {
			t.Fatalf("FindShortestPath returned error: %v", err)
		}

		expectedShortestPath := []string{locations[1].Slug, locations[4].Slug}
		if len(path) != len(expectedShortestPath) {
			t.Fatalf("Expected path length %d (shortest), got %d", len(expectedShortestPath), len(path))
		}

		for i, expected := range expectedShortestPath {
			if path[i] != expected {
				t.Fatalf("Expected path[%d] to be %s, got %s", i, expected, path[i])
			}
		}
	})
}

