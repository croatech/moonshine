package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"gorm.io/gorm"

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
	if err := seedLocations(db.DB()); err != nil {
		log.Printf("Failed to seed locations: %v", err)
	}
	seedUsers(db.DB())

	log.Println("Seed process completed!")
}

func seedAvatars(db *gorm.DB) {
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
		result := db.Where("image = ?", imagePath).First(&existingAvatar)
		if result.Error == nil {
			log.Printf("Avatar %s already exists, skipping", imagePath)
			continue
		}

		avatar := &domain.Avatar{
			Image:   imagePath,
			Private: false,
		}

		if err := db.Create(avatar).Error; err != nil {
			log.Printf("Failed to create avatar %s: %v", imagePath, err)
			continue
		}

		count++
		log.Printf("Created avatar %d: %s", i+1, imagePath)
	}

	log.Printf("Successfully created %d avatars", count)
}

func seedUsers(db *gorm.DB) {
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
	if err := db.First(&firstAvatar).Error; err != nil {
		log.Printf("No avatars found, creating user without avatar")
	}

	var moonshineLocation domain.Location
	if err := db.Where("slug = ?", "moonshine").First(&moonshineLocation).Error; err != nil {
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

func seedLocations(db *gorm.DB) error {
	log.Println("Seeding locations...")

	var moonshineLocation domain.Location
	if err := db.Where("slug = ?", "moonshine").First(&moonshineLocation).Error; err == nil {
		log.Println("Locations already exist, updating cell values...")
		if err := db.Exec("UPDATE locations SET cell = false WHERE slug IN ('moonshine', 'craft_shop', 'shop_of_artifacts', 'weapon_shop')").Error; err != nil {
			log.Printf("Failed to update city and shops cell values: %v", err)
		}
		if err := db.Exec("UPDATE locations SET cell = true WHERE slug LIKE '%cell'").Error; err != nil {
			log.Printf("Failed to update cells cell values: %v", err)
		}
		return nil
	}

	moonshineID := uuid.New()
	moonshineLocation = domain.Location{
		Model: domain.Model{
			ID: moonshineID,
		},
		Name:     "Moonshine",
		Slug:     "moonshine",
		Cell:     false,
		Inactive: false,
		Image:    "cities/moonshine/icon.jpg",
		ImageBg:  "cities/moonshine/bg.jpg",
	}

	if err := db.Create(&moonshineLocation).Error; err != nil {
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
			shopNameParts[i] = strings.Title(part)
		}
		shopName := strings.Join(shopNameParts, " ")
		shopLocation := domain.Location{
			Model: domain.Model{
				ID: shopID,
			},
			Name:     shopName,
			Slug:     shop.slug,
			Cell:     false,
			Inactive: false,
			Image:    fmt.Sprintf("cities/moonshine/%s/icon.png", shop.slug),
			ImageBg:  fmt.Sprintf("cities/moonshine/%s/bg.jpg", shop.slug),
		}

		if err := db.Create(&shopLocation).Error; err != nil {
			return fmt.Errorf("failed to create shop location %s: %w", shop.slug, err)
		}

		shopLocations[shop.slug] = shopID

		locationLocation := domain.LocationLocation{
			Model: domain.Model{
				ID: uuid.New(),
			},
			LocationID:     moonshineID,
			NearLocationID: shopID,
		}

		if err := db.Create(&locationLocation).Error; err != nil {
			return fmt.Errorf("failed to create location connection for %s: %w", shop.slug, err)
		}

		locationLocationReverse := domain.LocationLocation{
			Model: domain.Model{
				ID: uuid.New(),
			},
			LocationID:     shopID,
			NearLocationID: moonshineID,
		}

		if err := db.Create(&locationLocationReverse).Error; err != nil {
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
		cellLocation := domain.Location{
			Model: domain.Model{
				ID: cellID,
			},
			Name:     "",
			Slug:     cellSlug,
			Cell:     true,
			Inactive: false,
			Image:    fmt.Sprintf("wayward_pines/cells/%s.png", cellSlug),
			ImageBg:  "",
		}

		if err := db.Create(&cellLocation).Error; err != nil {
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

			var existingConnection domain.LocationLocation
			err := db.Where("location_id = ? AND near_location_id = ?", cellID, neighborID).
				First(&existingConnection).Error

			if err != nil {
				locationLocation := domain.LocationLocation{
					Model: domain.Model{
						ID: uuid.New(),
					},
					LocationID:     cellID,
					NearLocationID: neighborID,
				}

				if err := db.Create(&locationLocation).Error; err != nil {
					return fmt.Errorf("failed to create cell connection %d -> %d: %w", cellNum, neighborNum, err)
				}
			}

			var existingReverseConnection domain.LocationLocation
			err = db.Where("location_id = ? AND near_location_id = ?", neighborID, cellID).
				First(&existingReverseConnection).Error

			if err != nil {
				locationLocationReverse := domain.LocationLocation{
					Model: domain.Model{
						ID: uuid.New(),
					},
					LocationID:     neighborID,
					NearLocationID: cellID,
				}

				if err := db.Create(&locationLocationReverse).Error; err != nil {
					return fmt.Errorf("failed to create reverse cell connection %d -> %d: %w", neighborNum, cellNum, err)
				}
			}
		}
	}

	log.Println("Created all cell connections")

	return nil
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
