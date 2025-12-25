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

	// Очистка всех таблиц перед заполнением
	if err := truncateTables(db.DB()); err != nil {
		log.Fatalf("Failed to truncate tables: %v", err)
	}

	seedAvatars(db.DB())
	seedEquipmentCategories(db.DB())
	if err := seedEquipmentItems(db.DB()); err != nil {
		log.Printf("Failed to seed equipment items: %v", err)
	}
	if err := seedLocations(db.DB()); err != nil {
		log.Printf("Failed to seed locations: %v", err)
	}
	seedUsers(db.DB())

	log.Println("Seed process completed!")
}

func truncateTables(db *sqlx.DB) error {
	log.Println("Truncating all seed tables...")

	// Порядок важен - сначала зависимые таблицы, потом основные
	// Примечание: не очищаем таблицу users, чтобы не удалять реальных пользователей
	tables := []string{
		"user_equipment_items",
		"bot_equipment_items",
		"location_locations",
		"equipment_items",
		"equipment_categories",
		"locations",
		"avatars",
	}

	for _, table := range tables {
		query := fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE", table)
		if _, err := db.Exec(query); err != nil {
			return fmt.Errorf("failed to truncate %s: %w", table, err)
		}
		log.Printf("Truncated table: %s", table)
	}

	log.Println("All tables truncated successfully")
	return nil
}

func seedAvatars(db *sqlx.DB) {
	log.Println("Seeding avatars...")

	avatarRepo := repository.NewAvatarRepository(db)

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

		_, err := avatarRepo.FindByImage(imagePath)
		if err == nil {
			log.Printf("Avatar %s already exists, skipping", imagePath)
			continue
		}

		avatar := &domain.Avatar{
			Image:   imagePath,
			Private: false,
		}

		if err := avatarRepo.Create(avatar); err != nil {
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

	avatarRepo := repository.NewAvatarRepository(db)
	locationRepo := repository.NewLocationRepository(db)

	// Try to get first avatar
	var firstAvatar *domain.Avatar
	avatarsDir := "frontend/assets/images/players/avatars"
	files, err := filepath.Glob(filepath.Join(avatarsDir, "*.png"))
	if err == nil && len(files) > 0 {
		filename := filepath.Base(files[0])
		imagePath := filepath.Join("players/avatars", filename)
		firstAvatar, err = avatarRepo.FindByImage(imagePath)
		if err != nil {
			log.Printf("No avatars found, creating user without avatar")
		}
	} else {
		log.Printf("No avatars found, creating user without avatar")
	}

	moonshineLocation, err := locationRepo.FindStartLocation()
	if err != nil {
		log.Fatalf("Moonshine location not found, please seed locations first: %v", err)
	}

	user := &domain.User{
		Username:   "admin",
		Email:      "admin@gmail.com",
		Password:   hashedPassword,
		Hp:         20,
		CurrentHp:  20,
		Level:      1,
		Gold:       100,
		Exp:        0,
		FreeStats:  15,
		LocationID: moonshineLocation.ID,
	}

	if firstAvatar != nil && firstAvatar.ID != uuid.Nil {
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

	locationRepo := repository.NewLocationRepository(db)

	moonshineLocation, err := locationRepo.FindStartLocation()
	if err == nil && moonshineLocation != nil {
		log.Println("Locations already exist, updating cell values...")
		if _, err := db.Exec("UPDATE locations SET cell = false WHERE slug IN ('moonshine', 'craft_shop', 'shop_of_artifacts', 'weapon_shop')"); err != nil {
			log.Printf("Failed to update city and shops cell values: %v", err)
		}
		if _, err := db.Exec("UPDATE locations SET cell = true WHERE slug LIKE '%cell'"); err != nil {
			log.Printf("Failed to update cells cell values: %v", err)
		}
		return nil
	}

	moonshineLocation = &domain.Location{
		Name:     "Moonshine",
		Slug:     "moonshine",
		Cell:     false,
		Inactive: false,
		Image:    "cities/moonshine/icon.jpg",
		ImageBg:  "cities/moonshine/bg.jpg",
	}

	if err := locationRepo.Create(moonshineLocation); err != nil {
		return fmt.Errorf("failed to create Moonshine location: %w", err)
	}

	log.Printf("Created Moonshine location: %s", moonshineLocation.ID)

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
		shopNameParts := strings.Split(shop.name, "_")
		for i, part := range shopNameParts {
			if len(part) > 0 {
				shopNameParts[i] = strings.ToUpper(part[:1]) + strings.ToLower(part[1:])
			}
		}
		shopName := strings.Join(shopNameParts, " ")
		shopLocation := &domain.Location{
			Name:     shopName,
			Slug:     shop.slug,
			Cell:     false,
			Inactive: false,
			Image:    fmt.Sprintf("cities/moonshine/%s/icon.png", shop.slug),
			ImageBg:  fmt.Sprintf("cities/moonshine/%s/bg.jpg", shop.slug),
		}

		if err := locationRepo.Create(shopLocation); err != nil {
			return fmt.Errorf("failed to create shop location %s: %w", shop.slug, err)
		}

		shopLocations[shop.slug] = shopLocation.ID

		locLocID := uuid.New()
		locationLocationQuery := `INSERT INTO location_locations (id, location_id, near_location_id) 
			VALUES ($1, $2, $3)`
		if _, err := db.Exec(locationLocationQuery, locLocID, moonshineLocation.ID, shopLocation.ID); err != nil {
			return fmt.Errorf("failed to create location connection for %s: %w", shop.slug, err)
		}

		locLocReverseID := uuid.New()
		if _, err := db.Exec(locationLocationQuery, locLocReverseID, shopLocation.ID, moonshineLocation.ID); err != nil {
			return fmt.Errorf("failed to create reverse location connection for %s: %w", shop.slug, err)
		}

		log.Printf("Created shop location: %s (%s)", shop.slug, shopLocation.ID)
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

		cellSlug := fmt.Sprintf("%dcell", cellNum)
		cellLocation := &domain.Location{
			Name:     "",
			Slug:     cellSlug,
			Cell:     true,
			Inactive: false,
			Image:    fmt.Sprintf("wayward_pines/cells/%s.png", cellSlug),
			ImageBg:  "",
		}

		if err := locationRepo.Create(cellLocation); err != nil {
			return fmt.Errorf("failed to create cell location %d: %w", cellNum, err)
		}

		cellLocations[cellNum] = cellLocation.ID
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
			err := db.QueryRow(
				"SELECT id FROM location_locations WHERE location_id = $1 AND near_location_id = $2",
				cellID, neighborID).Scan(&existingConnectionID)

			if err != nil {
				locLocID := uuid.New()
				locLocQuery := `INSERT INTO location_locations (id, location_id, near_location_id) 
					VALUES ($1, $2, $3)`
				if _, err := db.Exec(locLocQuery, locLocID, cellID, neighborID); err != nil {
					return fmt.Errorf("failed to create cell connection %d -> %d: %w", cellNum, neighborNum, err)
				}
			}

			var existingReverseConnectionID uuid.UUID
			err = db.QueryRow(
				"SELECT id FROM location_locations WHERE location_id = $1 AND near_location_id = $2",
				neighborID, cellID).Scan(&existingReverseConnectionID)

			if err != nil {
				locLocReverseID := uuid.New()
				locLocQuery := `INSERT INTO location_locations (id, location_id, near_location_id) 
					VALUES ($1, $2, $3)`
				if _, err := db.Exec(locLocQuery, locLocReverseID, neighborID, cellID); err != nil {
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
		err := db.QueryRow("SELECT id FROM equipment_categories WHERE type = $1", cat.typ).Scan(&existingID)
		if err == nil {
			log.Printf("Equipment category %s already exists, skipping", cat.name)
			continue
		}

		categoryID := uuid.New()
		query := `INSERT INTO equipment_categories (id, name, type) VALUES ($1, $2, $3)`
		if _, err := db.Exec(query, categoryID, cat.name, cat.typ); err != nil {
			log.Printf("Failed to create equipment category %s: %v", cat.name, err)
			continue
		}

		log.Printf("Created equipment category: %s (%s)", cat.name, cat.typ)
	}

	log.Println("Equipment categories seeding completed!")
}

type equipmentFileInfo struct {
	path          string
	categoryType  string
	name          string
	englishName   string // English name for slug generation
	price         uint
	attack        uint
	requiredLevel uint
	hp            uint
	defense       uint
}

// Russian translations for equipment names: english -> russian
var equipmentRussianNames = map[string]string{
	// Axes
	"Iron_Axe":          "Железный топор",
	"Battle_Axe":        "Боевой топор",
	"War_Axe":           "Топор войны",
	"Great_Axe":         "Великий топор",
	"Executioner_Axe":   "Топор палача",
	"Bearded_Axe":       "Бородатый топор",
	"Bloodthirster_Axe": "Кровожадный топор",
	"Razor_Edge_Axe":    "Бритвенный топор",
	"Titan_Blade":       "Клинок титана",
	"Dragon_Slayer":     "Убийца драконов",
	"Skull_Crusher":     "Сокрушитель черепов",
	"Berserker_Axe":     "Топор берсерка",
	"Frostbite_Axe":     "Морозный топор",
	"Inferno_Cleave":    "Инфернальный топор",
	"Thunder_Strike":    "Громовой удар",
	"Soul_Eater":        "Пожиратель душ",
	"Death_Bringer":     "Вестник смерти",
	"Abyss_Reaper":      "Жнец бездны",
	"Eternal_Wrath":     "Вечная ярость",
	"Apocalypse_Axe":    "Топор апокалипсиса",

	// Knifes
	"Rusty_Dagger":   "Ржавый кинжал",
	"Hunting_Knife":  "Охотничий нож",
	"Stiletto":       "Стилет",
	"Combat_Dagger":  "Боевой кинжал",
	"Shadow_Blade":   "Теневой клинок",
	"Venom_Dagger":   "Ядовитый кинжал",
	"Rapier":         "Рапира",
	"Assassin_Knife": "Нож ассасина",
	"Phantom_Dagger": "Призрачный кинжал",
	"Blood_Thorn":    "Кровавый шип",
	"Dark_Edge":      "Темное лезвие",
	"Crimson_Fang":   "Багровый клык",
	"Shadow_Strike":  "Удар тени",
	"Death_Kiss":     "Поцелуй смерти",
	"Soul_Reaver":    "Жнец душ",
	"Void_Cutter":    "Клинок пустоты",

	// Maces
	"Iron_Club":        "Железная дубина",
	"Battle_Mace":      "Боевая булава",
	"War_Hammer":       "Боевой молот",
	"Flail":            "Цеп",
	"Morning_Star":     "Утренняя звезда",
	"Maul":             "Колотушка",
	"Thunder_Mace":     "Громовая булава",
	"Bone_Crusher":     "Костолом",
	"Skull_Hammer":     "Молот черепов",
	"Stone_Breaker":    "Камнедробитель",
	"Earth_Shaker":     "Сотрясатель земли",
	"Mighty_Maul":      "Могучий молот",
	"Crushing_Blow":    "Сокрушительный удар",
	"Dragon_Mace":      "Драконья булава",
	"Titan_Strike":     "Удар титана",
	"Celestial_Hammer": "Небесный молот",
	"Divine_Judgment":  "Божий суд",
	"Eternal_Crusher":  "Вечный сокрушитель",
	"World_Breaker":    "Разрушитель миров",
	"Oblivion_Mace":    "Булава забвения",

	// Armor (chest)
	"Cloth_Tunic":       "Тканевая туника",
	"Leather_Vest":      "Кожаный жилет",
	"Chainmail":         "Кольчуга",
	"Plate_Armor":       "Латная броня",
	"Scale_Mail":        "Чешуйчатая броня",
	"Brigandine":        "Бригантина",
	"Cuirass":           "Кираса",
	"Hauberk":           "Хауберг",
	"Breastplate":       "Нагрудная пластина",
	"Chestguard":        "Нагрудник",
	"Battle_Armor":      "Боевая броня",
	"War_Plate":         "Доспехи войны",
	"Dragon_Scale":      "Драконья чешуя",
	"Titanium_Plate":    "Титановая броня",
	"Mythril_Mail":      "Мифриловая кольчуга",
	"Divine_Armor":      "Божественные доспехи",
	"Celestial_Plate":   "Небесная броня",
	"Eternal_Guard":     "Вечная защита",
	"Legendary_Cuirass": "Легендарная кираса",
	"Immortal_Plate":    "Бессмертная броня",
	"Godforged_Armor":   "Богокованные доспехи",
}

func seedEquipmentItems(db *sqlx.DB) error {
	log.Println("Seeding equipment items...")

	baseDir := "frontend/assets/images/equipment_items"

	// Маппинг директорий к типам категорий
	categoryMap := map[string]string{
		"chest":  "chest",
		"belt":   "belt",
		"head":   "head",
		"neck":   "neck",
		"weapon": "weapon",
		"shield": "shield",
		"legs":   "legs",
		"feet":   "feet",
		"arms":   "arms",
		"hands":  "hands",
		"ring":   "ring",
	}

	// Маппинг поддиректорий weapon к типам
	weaponSubdirMap := map[string]string{
		"axes":   "weapon",
		"knifes": "weapon",
		"maces":  "weapon",
	}

	var allFiles []equipmentFileInfo

	// Собираем все файлы
	err := filepath.Walk(baseDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		ext := strings.ToLower(filepath.Ext(path))
		if ext != ".png" {
			return nil
		}

		relPath, err := filepath.Rel(baseDir, path)
		if err != nil {
			return err
		}

		parts := strings.Split(relPath, string(filepath.Separator))
		if len(parts) < 1 {
			return nil
		}

		dir := parts[0]
		subdir := ""
		if len(parts) > 2 && dir == "weapon" {
			subdir = parts[1]
		}

		categoryType := categoryMap[dir]
		if categoryType == "" && dir == "weapon" && subdir != "" {
			categoryType = weaponSubdirMap[subdir]
		}

		if categoryType == "" {
			return nil
		}

		fileName := filepath.Base(path)
		fileInfo := equipmentFileInfo{
			path:         path,
			categoryType: categoryType,
		}

		// Парсим имя файла: порядок-название-стоимость-урон-уровень-хп-защита.png
		if !parseEquipmentFileName(fileName, &fileInfo) {
			log.Printf("Failed to parse equipment file name: %s, skipping", fileName)
			return nil
		}

		allFiles = append(allFiles, fileInfo)
		return nil
	})

	if err != nil {
		return fmt.Errorf("walk equipment items directory: %w", err)
	}

	// Получаем категории из БД
	categoryIDs := make(map[string]uuid.UUID)
	for catType := range categoryMap {
		var catID uuid.UUID
		err := db.QueryRow("SELECT id FROM equipment_categories WHERE type = $1", catType).Scan(&catID)
		if err == nil {
			categoryIDs[catType] = catID
		}
	}
	// Получаем weapon категорию
	var weaponCatID uuid.UUID
	err = db.QueryRow("SELECT id FROM equipment_categories WHERE type = 'weapon'").Scan(&weaponCatID)
	if err == nil {
		categoryIDs["weapon"] = weaponCatID
	}

	// Загружаем items в БД
	count := 0
	for _, file := range allFiles {
		catID := categoryIDs[file.categoryType]
		if catID == uuid.Nil {
			log.Printf("Category %s not found, skipping item %s", file.categoryType, filepath.Base(file.path))
			continue
		}

		// Формируем путь для БД (относительно frontend/assets/images)
		dbImagePath := strings.TrimPrefix(file.path, "frontend/assets/images/")
		dbImagePath = strings.ReplaceAll(dbImagePath, "\\", "/")

		// Проверяем, существует ли уже такой item
		var existingID uuid.UUID
		err := db.QueryRow("SELECT id FROM equipment_items WHERE image = $1", dbImagePath).Scan(&existingID)
		if err == nil {
			log.Printf("Equipment item %s already exists, skipping", dbImagePath)
			continue
		}

		// Generate slug from English name (not translated Russian name)
		slug := generateSlug(file.englishName)

		// Создаем equipment item
		itemID := uuid.New()
		query := `INSERT INTO equipment_items 
			(id, name, slug, attack, defense, hp, required_level, price, equipment_category_id, image) 
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

		_, err = db.Exec(query,
			itemID,
			file.name,
			slug,
			file.attack,
			file.defense,
			file.hp,
			file.requiredLevel,
			file.price,
			catID,
			dbImagePath,
		)

		if err != nil {
			log.Printf("Failed to create equipment item %s: %v", file.name, err)
			continue
		}

		count++
		log.Printf("Created equipment item: %s (level %d, attack %d, defense %d, hp %d, price %d)",
			file.name, file.requiredLevel, file.attack, file.defense, file.hp, file.price)
	}

	log.Printf("Equipment items seeding completed! Created %d items", count)
	return nil
}

func parseEquipmentFileName(filename string, info *equipmentFileInfo) bool {
	ext := filepath.Ext(filename)
	base := strings.TrimSuffix(filename, ext)

	// Формат: порядок-название-стоимость-урон-уровень-хп-защита
	// Пример: 1-Dagger-15-1-1-0-0.png
	parts := strings.Split(base, "-")
	if len(parts) < 7 {
		return false
	}

	// Пропускаем порядок (первая часть), не нужен
	// Название - все части между порядком и числами
	// Ищем где начинаются числа (стоимость)
	nameParts := []string{}
	numStartIdx := -1
	for i := 1; i < len(parts); i++ {
		// Проверяем, является ли это числом (стоимость)
		if _, err := strconv.Atoi(parts[i]); err == nil {
			// Проверяем, что следующая часть тоже число (урон)
			if i+1 < len(parts) {
				if _, err2 := strconv.Atoi(parts[i+1]); err2 == nil {
					numStartIdx = i
					break
				}
			}
		}
		nameParts = append(nameParts, parts[i])
	}

	if numStartIdx == -1 || len(nameParts) == 0 {
		return false
	}

	// English name with underscores (for slug generation)
	info.englishName = strings.Join(nameParts, "_")

	// Capitalize each part for lookup in Russian names map
	capitalizedParts := make([]string, len(nameParts))
	for i, part := range nameParts {
		if len(part) > 0 {
			capitalizedParts[i] = strings.ToUpper(part[:1]) + strings.ToLower(part[1:])
		}
	}
	capitalizedEnglishName := strings.Join(capitalizedParts, "_")

	// Try to find Russian translation
	if russianName, exists := equipmentRussianNames[capitalizedEnglishName]; exists {
		info.name = russianName
	} else {
		// Fallback to English name with spaces if no translation found
		info.name = strings.Join(nameParts, " ")
		info.name = strings.ReplaceAll(info.name, "_", " ")
	}

	// Парсим числа: стоимость-урон-уровень-хп-защита
	if numStartIdx+4 < len(parts) {
		if price, err := strconv.ParseUint(parts[numStartIdx], 10, 32); err == nil {
			info.price = uint(price)
		} else {
			return false
		}
		if attack, err := strconv.ParseUint(parts[numStartIdx+1], 10, 32); err == nil {
			info.attack = uint(attack)
		} else {
			return false
		}
		if level, err := strconv.ParseUint(parts[numStartIdx+2], 10, 32); err == nil {
			info.requiredLevel = uint(level)
		} else {
			return false
		}
		if hp, err := strconv.ParseUint(parts[numStartIdx+3], 10, 32); err == nil {
			info.hp = uint(hp)
		} else {
			return false
		}
		if defense, err := strconv.ParseUint(parts[numStartIdx+4], 10, 32); err == nil {
			info.defense = uint(defense)
		} else {
			return false
		}
		return true
	}

	return false
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

// generateSlug creates a URL-friendly slug from a name
func generateSlug(name string) string {
	// Convert to lowercase
	slug := strings.ToLower(name)

	// Replace spaces and underscores with hyphens
	slug = strings.ReplaceAll(slug, " ", "-")
	slug = strings.ReplaceAll(slug, "_", "-")

	// Remove any special characters, keep only alphanumeric and hyphens
	var result strings.Builder
	for _, r := range slug {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			result.WriteRune(r)
		}
	}

	slug = result.String()

	// Remove multiple consecutive hyphens
	for strings.Contains(slug, "--") {
		slug = strings.ReplaceAll(slug, "--", "-")
	}

	// Trim hyphens from start and end
	slug = strings.Trim(slug, "-")

	return slug
}
