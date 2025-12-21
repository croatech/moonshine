package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"

	"moonshine/internal/domain"
	"moonshine/internal/repository"
	"moonshine/internal/util"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println(".env not loaded, relying on environment")
	}

	db, err := repository.New()
	if err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}
	defer db.Close()

	log.Println("Starting seed process...")

	seedAvatars(db.DB())
	seedEquipmentCategories(db.DB())
	if err := seedLocations(db.DB()); err != nil {
		log.Printf("Failed to seed locations: %v", err)
	}
	seedUsers(db.DB())

	log.Println("Seed process completed!")
}

func seedAvatars(db *sqlx.DB) {
	log.Println("Seeding avatars...")

	avatarsDir := "frontend/assets/images/players/avatars"
	if _, err := os.Stat(avatarsDir); os.IsNotExist(err) {
		log.Printf("Avatars directory not found: %s, skipping avatars", avatarsDir)
		return
	}

	files, err := filepath.Glob(filepath.Join(avatarsDir, "*.png"))
	if err != nil {
		log.Printf("Failed to read avatars directory: %v, skipping avatars", err)
		return
	}

	if len(files) == 0 {
		log.Println("No PNG avatar files found, skipping avatars")
		return
	}

	count := 0

	for i, file := range files {
		filename := filepath.Base(file)
		imagePath := filepath.Join("players/avatars", filename)

		var existingAvatar domain.Avatar
		err := db.Get(&existingAvatar, "SELECT id, image, private, created_at, updated_at FROM avatars WHERE image = $1", imagePath)
		if err == nil {
			log.Printf("Avatar %s already exists, skipping", imagePath)
			continue
		}

		avatarID := uuid.New()
		now := time.Now()
		query := `INSERT INTO avatars (id, image, private, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)`
		if _, err := db.Exec(query, avatarID, imagePath, false, now, now); err != nil {
			log.Printf("Failed to create avatar %s: %v", imagePath, err)
			continue
		}

		count++
		log.Printf("Created avatar %d: %s", i+1, imagePath)
	}

	log.Printf("Successfully created %d avatars", count)
}

func seedUsers(db *sqlx.DB) {
	log.Println("Seeding users...")

	userRepo := repository.NewUserRepository(db)

	existingUser, err := userRepo.FindByUsername("admin")
	if err == nil && existingUser != nil {
		log.Println("User 'admin' already exists, skipping")
		return
	}

	hashedPassword, err := util.HashPassword("password")
	if err != nil {
		log.Fatalf("Failed to hash password: %v", err)
	}
	var firstAvatar domain.Avatar
	if err := db.Get(&firstAvatar, "SELECT id, image, private, created_at, updated_at FROM avatars LIMIT 1"); err != nil {
		log.Printf("No avatars found, creating user without avatar")
	}

	var moonshineLocation domain.Location
	if err := db.Get(&moonshineLocation,
		"SELECT id, name, slug, cell, inactive, image, image_bg, created_at, updated_at FROM locations WHERE slug = $1",
		"moonshine"); err != nil {
		log.Fatalf("Moonshine location not found, please seed locations first: %v", err)
	}

	user := &domain.User{
		Username:   "admin",
		Email:      "admin@gmail.com",
		Password:   hashedPassword,
		Hp:         20,
		Level:      1,
		Gold:       100,
		Exp:        0,
		FreeStats:  15,
		LocationID: moonshineLocation.ID,
	}

	if firstAvatar.ID != uuid.Nil {
		avatarID := firstAvatar.ID
		user.AvatarID = &avatarID
		log.Printf("Assigned avatar ID %s to user", firstAvatar.ID.String())
	}

	log.Printf("Assigned location ID %s (Moonshine) to user", moonshineLocation.ID.String())

	if err := userRepo.Create(user); err != nil {
		log.Fatalf("Failed to create user: %v", err)
	}

	log.Printf("Successfully created user: %s (%s)", user.Username, user.Email)
}

func seedLocations(db *sqlx.DB) error {
	log.Println("Seeding locations...")

	var moonshineLocation domain.Location
	if err := db.Get(&moonshineLocation,
		"SELECT id, name, slug, cell, inactive, image, image_bg, created_at, updated_at FROM locations WHERE slug = $1",
		"moonshine"); err == nil {
		log.Println("Locations already exist, updating cell values...")
		if _, err := db.Exec("UPDATE locations SET cell = false WHERE slug IN ('moonshine', 'craft_shop', 'shop_of_artifacts', 'weapon_shop')"); err != nil {
			log.Printf("Failed to update city and shops cell values: %v", err)
		}
		if _, err := db.Exec("UPDATE locations SET cell = true WHERE slug LIKE '%cell'"); err != nil {
			log.Printf("Failed to update cells cell values: %v", err)
		}
		return nil
	}

	moonshineID := uuid.New()
	now := time.Now()
	moonshineLocation = domain.Location{
		Model: domain.Model{
			ID:        moonshineID,
			CreatedAt: now,
			UpdatedAt: now,
		},
		Name:     "Moonshine",
		Slug:     "moonshine",
		Cell:     false,
		Inactive: false,
		Image:    "cities/moonshine/icon.jpg",
		ImageBg:  "cities/moonshine/bg.jpg",
	}

	query := `INSERT INTO locations (id, name, slug, cell, inactive, image, image_bg, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	if _, err := db.Exec(query, moonshineLocation.ID, moonshineLocation.Name, moonshineLocation.Slug,
		moonshineLocation.Cell, moonshineLocation.Inactive, moonshineLocation.Image, moonshineLocation.ImageBg,
		moonshineLocation.CreatedAt, moonshineLocation.UpdatedAt); err != nil {
		return fmt.Errorf("failed to create Moonshine location: %w", err)
	}

	log.Printf("Created Moonshine location: %s", moonshineID)

	shops := []struct {
		name string
		slug string
	}{
		{"craft_shop", "craft_shop"},
		{"shop_of_artifacts", "shop_of_artifacts"},
		{"weapon_shop", "weapon_shop"},
	}

	shopLocations := make(map[string]uuid.UUID)

	for _, shop := range shops {
		shopID := uuid.New()
		shopNameParts := strings.Split(shop.name, "_")
		for i, part := range shopNameParts {
			if len(part) > 0 {
				shopNameParts[i] = strings.ToUpper(part[:1]) + strings.ToLower(part[1:])
			}
		}
		shopName := strings.Join(shopNameParts, " ")
		shopNow := time.Now()
		shopLocation := domain.Location{
			Model: domain.Model{
				ID:        shopID,
				CreatedAt: shopNow,
				UpdatedAt: shopNow,
			},
			Name:     shopName,
			Slug:     shop.slug,
			Cell:     false,
			Inactive: false,
			Image:    fmt.Sprintf("cities/moonshine/%s/icon.png", shop.slug),
			ImageBg:  fmt.Sprintf("cities/moonshine/%s/bg.jpg", shop.slug),
		}

		shopQuery := `INSERT INTO locations (id, name, slug, cell, inactive, image, image_bg, created_at, updated_at) 
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
		if _, err := db.Exec(shopQuery, shopLocation.ID, shopLocation.Name, shopLocation.Slug,
			shopLocation.Cell, shopLocation.Inactive, shopLocation.Image, shopLocation.ImageBg,
			shopLocation.CreatedAt, shopLocation.UpdatedAt); err != nil {
			return fmt.Errorf("failed to create shop location %s: %w", shop.slug, err)
		}

		shopLocations[shop.slug] = shopID

		locLocID := uuid.New()
		locLocNow := time.Now()
		locationLocationQuery := `INSERT INTO location_locations (id, location_id, near_location_id, created_at, updated_at) 
			VALUES ($1, $2, $3, $4, $5)`
		if _, err := db.Exec(locationLocationQuery, locLocID, moonshineID, shopID, locLocNow, locLocNow); err != nil {
			return fmt.Errorf("failed to create location connection for %s: %w", shop.slug, err)
		}

		locLocReverseID := uuid.New()
		locLocReverseNow := time.Now()
		if _, err := db.Exec(locationLocationQuery, locLocReverseID, shopID, moonshineID, locLocReverseNow, locLocReverseNow); err != nil {
			return fmt.Errorf("failed to create reverse location connection for %s: %w", shop.slug, err)
		}

		log.Printf("Created shop location: %s (%s)", shop.slug, shopID)
	}

	cellsDir := "frontend/assets/images/locations/wayward_pines/cells"
	files, err := filepath.Glob(filepath.Join(cellsDir, "*.png"))
	if err != nil {
		return fmt.Errorf("failed to read cells directory: %w", err)
	}

	if len(files) != 64 {
		return fmt.Errorf("expected 64 cell files, found %d", len(files))
	}

	sort.Slice(files, func(i, j int) bool {
		numI := extractCellNumber(files[i])
		numJ := extractCellNumber(files[j])
		return numI < numJ
	})

	cellLocations := make(map[int]uuid.UUID)

	for _, file := range files {
		cellNum := extractCellNumber(file)
		if cellNum == 0 {
			continue
		}

		cellID := uuid.New()
		cellSlug := fmt.Sprintf("%dcell", cellNum)
		cellNow := time.Now()
		cellLocation := domain.Location{
			Model: domain.Model{
				ID:        cellID,
				CreatedAt: cellNow,
				UpdatedAt: cellNow,
			},
			Name:     "",
			Slug:     cellSlug,
			Cell:     true,
			Inactive: false,
			Image:    fmt.Sprintf("wayward_pines/cells/%s.png", cellSlug),
			ImageBg:  "",
		}

		cellQuery := `INSERT INTO locations (id, name, slug, cell, inactive, image, image_bg, created_at, updated_at) 
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
		if _, err := db.Exec(cellQuery, cellLocation.ID, cellLocation.Name, cellLocation.Slug,
			cellLocation.Cell, cellLocation.Inactive, cellLocation.Image, cellLocation.ImageBg,
			cellLocation.CreatedAt, cellLocation.UpdatedAt); err != nil {
			return fmt.Errorf("failed to create cell location %d: %w", cellNum, err)
		}

		cellLocations[cellNum] = cellID
	}

	log.Printf("Created %d cell locations", len(cellLocations))

	for cellNum := 1; cellNum <= 64; cellNum++ {
		cellID := cellLocations[cellNum]
		row := (cellNum - 1) / 8
		col := (cellNum - 1) % 8

		neighbors := []int{}

		if col > 0 {
			neighbors = append(neighbors, cellNum-1)
		}
		if col < 7 {
			neighbors = append(neighbors, cellNum+1)
		}
		if row > 0 {
			neighbors = append(neighbors, cellNum-8)
		}
		if row < 7 {
			neighbors = append(neighbors, cellNum+8)
		}

		for _, neighborNum := range neighbors {
			neighborID := cellLocations[neighborNum]

			var existingConnectionID uuid.UUID
			err := db.Get(&existingConnectionID,
				"SELECT id FROM location_locations WHERE location_id = $1 AND near_location_id = $2",
				cellID, neighborID)

			if err != nil {
				locLocID := uuid.New()
				locLocNow := time.Now()
				locLocQuery := `INSERT INTO location_locations (id, location_id, near_location_id, created_at, updated_at) 
					VALUES ($1, $2, $3, $4, $5)`
				if _, err := db.Exec(locLocQuery, locLocID, cellID, neighborID, locLocNow, locLocNow); err != nil {
					return fmt.Errorf("failed to create cell connection %d -> %d: %w", cellNum, neighborNum, err)
				}
			}

			var existingReverseConnectionID uuid.UUID
			err = db.Get(&existingReverseConnectionID,
				"SELECT id FROM location_locations WHERE location_id = $1 AND near_location_id = $2",
				neighborID, cellID)

			if err != nil {
				locLocReverseID := uuid.New()
				locLocReverseNow := time.Now()
				locLocQuery := `INSERT INTO location_locations (id, location_id, near_location_id, created_at, updated_at) 
					VALUES ($1, $2, $3, $4, $5)`
				if _, err := db.Exec(locLocQuery, locLocReverseID, neighborID, cellID, locLocReverseNow, locLocReverseNow); err != nil {
					return fmt.Errorf("failed to create reverse cell connection %d -> %d: %w", neighborNum, cellNum, err)
				}
			}
		}
	}

	log.Println("Created all cell connections")

	return nil
}

func seedEquipmentCategories(db *sqlx.DB) {
	log.Println("Seeding equipment categories...")

	categories := []struct {
		name string
		typ  string
	}{
		{"Chest", "chest"},
		{"Belt", "belt"},
		{"Head", "head"},
		{"Neck", "neck"},
		{"Weapon", "weapon"},
		{"Shield", "shield"},
		{"Legs", "legs"},
		{"Feet", "feet"},
		{"Arms", "arms"},
		{"Hands", "hands"},
		{"Ring", "ring"},
	}

	for _, cat := range categories {
		var existingID uuid.UUID
		err := db.Get(&existingID, "SELECT id FROM equipment_categories WHERE type = $1", cat.typ)
		if err == nil {
			log.Printf("Equipment category %s already exists, skipping", cat.name)
			continue
		}

		categoryID := uuid.New()
		now := time.Now()
		query := `INSERT INTO equipment_categories (id, name, type, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)`
		if _, err := db.Exec(query, categoryID, cat.name, cat.typ, now, now); err != nil {
			log.Printf("Failed to create equipment category %s: %v", cat.name, err)
			continue
		}

		log.Printf("Created equipment category: %s (%s)", cat.name, cat.typ)
	}

	log.Println("Equipment categories seeding completed!")
}

func extractCellNumber(filename string) int {
	base := filepath.Base(filename)
	base = strings.TrimSuffix(base, ".png")
	base = strings.TrimSuffix(base, "cell")

	num, err := strconv.Atoi(base)
	if err != nil {
		return 0
	}

	return num
}
