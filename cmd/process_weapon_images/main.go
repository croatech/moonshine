package main

import (
	"fmt"
	"image"
	"image/gif"
	"image/png"
	"io/fs"
	"log"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

var axeRussianNames = map[int]string{
	1:  "Железный топор",
	2:  "Боевой топор",
	3:  "Военный топор",
	4:  "Большой топор",
	5:  "Топор палача",
	6:  "Бородатый топор",
	7:  "Кровожадный топор",
	8:  "Топор бритвенного края",
	9:  "Клинок титана",
	10: "Драконобой",
	11: "Сокрушитель черепов",
	12: "Топор берсерка",
	13: "Топор обморожения",
	14: "Раскол инферно",
	15: "Удар грома",
	16: "Пожиратель душ",
}

var axeEnglishSlugs = map[int]string{
	1:  "iron-axe",
	2:  "battle-axe",
	3:  "war-axe",
	4:  "great-axe",
	5:  "executioner-axe",
	6:  "bearded-axe",
	7:  "bloodthirster-axe",
	8:  "razor-edge-axe",
	9:  "titan-blade",
	10: "dragon-slayer",
	11: "skull-crusher",
	12: "berserker-axe",
	13: "frostbite-axe",
	14: "inferno-cleave",
	15: "thunder-strike",
	16: "soul-eater",
}

var maceRussianNames = map[int]string{
	1:  "Железная дубина",
	2:  "Боевая булава",
	3:  "Военный молот",
	4:  "Цеп",
	5:  "Утренняя звезда",
	6:  "Молот",
	7:  "Громовая булава",
	8:  "Костолом",
	9:  "Молот черепов",
	10: "Камнедробитель",
	11: "Сотрясатель земли",
	12: "Могучий молот",
	13: "Сокрушающий удар",
	14: "Драконья булава",
	15: "Удар титана",
	16: "Небесный молот",
}

var maceEnglishSlugs = map[int]string{
	1:  "iron-club",
	2:  "battle-mace",
	3:  "war-hammer",
	4:  "flail",
	5:  "morning-star",
	6:  "maul",
	7:  "thunder-mace",
	8:  "bone-crusher",
	9:  "skull-hammer",
	10: "stone-breaker",
	11: "earth-shaker",
	12: "mighty-maul",
	13: "crushing-blow",
	14: "dragon-mace",
	15: "titan-strike",
	16: "celestial-hammer",
}

var knifeRussianNames = map[int]string{
	1:  "Ржавый кинжал",
	2:  "Охотничий нож",
	3:  "Стилет",
	4:  "Боевой кинжал",
	5:  "Теневой клинок",
	6:  "Ядовитый кинжал",
	7:  "Рапира",
	8:  "Нож ассасина",
	9:  "Призрачный кинжал",
	10: "Кровавый шип",
	11: "Темное лезвие",
	12: "Багровый клык",
	13: "Удар тени",
	14: "Поцелуй смерти",
	15: "Жнец душ",
	16: "Клинок пустоты",
}

var knifeEnglishSlugs = map[int]string{
	1:  "rusty-dagger",
	2:  "hunting-knife",
	3:  "stiletto",
	4:  "combat-dagger",
	5:  "shadow-blade",
	6:  "venom-dagger",
	7:  "rapier",
	8:  "assassin-knife",
	9:  "phantom-dagger",
	10: "blood-thorn",
	11: "dark-edge",
	12: "crimson-fang",
	13: "shadow-strike",
	14: "death-kiss",
	15: "soul-reaver",
	16: "void-cutter",
}

// Calculate balanced weapon stats (max level 16, max attack ~200)
func calculateWeaponStats(level int, variant int) (attack, defense, hp, price int) {
	if level < 1 {
		level = 1
	}
	if level > 16 {
		level = 16
	}

	// Weapons have NO defense
	defense = 0

	// Base attack progression: exponential from ~5-10 at level 1 to ~180-220 at level 16
	attackBase := 5.0 * math.Pow(1.265, float64(level-1))
	attackMax := 10.0 * math.Pow(1.265, float64(level-1))

	// HP progression: from ~5-10 at level 1 to 80-100 at level 16
	hpBase := 5.0 * math.Pow(1.21, float64(level-1))
	hpMax := 10.0 * math.Pow(1.21, float64(level-1))

	switch variant {
	case 0: // Balanced
		attack = int(attackBase + (attackMax-attackBase)*0.5)
		hp = int(hpBase + (hpMax-hpBase)*0.5)
	case 1: // High attack, low HP
		attack = int(attackMax * 0.95)
		hp = int(hpBase * 0.8)
	case 2: // Lower attack, higher HP
		attack = int(attackBase * 1.15)
		hp = int(hpMax * 0.95)
	}

	// Adjust level 16 to exact targets (~200 attack)
	if level == 16 {
		switch variant {
		case 0: // Balanced
			attack = 200
			hp = 90
		case 1: // High attack, low HP
			attack = 220
			hp = 80
		case 2: // Lower attack, higher HP
			attack = 180
			hp = 100
		}
	}

	// Price calculation: start from 3-5 gold at level 1
	// Use level-based pricing with small variation based on variant
	basePrice := float64(level) * 2.5

	switch variant {
	case 0: // Balanced
		price = int(basePrice)
	case 1: // High attack - slightly more expensive
		price = int(basePrice * 1.1)
	case 2: // High HP - slightly cheaper
		price = int(basePrice * 0.9)
	}

	// Ensure minimum price of 3 gold
	if price < 3 {
		price = 3
	}

	return attack, defense, hp, price
}

type fileWithTime struct {
	path    string
	modTime time.Time
}

func convertGifToPng(gifPath string) error {
	// Open the GIF file
	gifFile, err := os.Open(gifPath)
	if err != nil {
		return fmt.Errorf("failed to open GIF: %w", err)
	}
	defer gifFile.Close()

	// Decode the GIF
	gifImg, err := gif.Decode(gifFile)
	if err != nil {
		return fmt.Errorf("failed to decode GIF: %w", err)
	}

	// Create PNG file with same name but .png extension
	pngPath := strings.TrimSuffix(gifPath, filepath.Ext(gifPath)) + ".png"
	pngFile, err := os.Create(pngPath)
	if err != nil {
		return fmt.Errorf("failed to create PNG: %w", err)
	}
	defer pngFile.Close()

	// Convert to RGBA if needed
	bounds := gifImg.Bounds()
	rgba := image.NewRGBA(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			rgba.Set(x, y, gifImg.At(x, y))
		}
	}

	// Encode as PNG
	if err := png.Encode(pngFile, rgba); err != nil {
		return fmt.Errorf("failed to encode PNG: %w", err)
	}

	// Remove the original GIF
	if err := os.Remove(gifPath); err != nil {
		return fmt.Errorf("failed to remove GIF: %w", err)
	}

	log.Printf("Converted and removed: %s -> %s", gifPath, pngPath)
	return nil
}

func processWeaponCategory(categoryPath string, categoryName string) error {
	log.Printf("Processing category: %s", categoryName)

	// Select appropriate name maps based on category
	var russianNames map[int]string
	var englishSlugs map[int]string
	var categoryOffset int

	if categoryName == "axes" {
		russianNames = axeRussianNames
		englishSlugs = axeEnglishSlugs
		categoryOffset = 0
	} else if categoryName == "maces" {
		russianNames = maceRussianNames
		englishSlugs = maceEnglishSlugs
		categoryOffset = 1 // Different variant offset for maces
	} else if categoryName == "knifes" {
		russianNames = knifeRussianNames
		englishSlugs = knifeEnglishSlugs
		categoryOffset = 2 // Different variant offset for knifes
	} else {
		return fmt.Errorf("unknown category: %s", categoryName)
	}

	// Collect all image files with their modification times
	var files []fileWithTime

	err := filepath.WalkDir(categoryPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		ext := strings.ToLower(filepath.Ext(path))
		if ext == ".gif" || ext == ".png" || ext == ".jpg" || ext == ".jpeg" {
			info, err := d.Info()
			if err != nil {
				return err
			}

			// Convert GIF to PNG first
			if ext == ".gif" {
				if err := convertGifToPng(path); err != nil {
					log.Printf("Warning: failed to convert %s: %v", path, err)
					return nil
				}
				// Update path to PNG
				path = strings.TrimSuffix(path, filepath.Ext(path)) + ".png"
				info, err = os.Stat(path)
				if err != nil {
					return err
				}
			}

			files = append(files, fileWithTime{
				path:    path,
				modTime: info.ModTime(),
			})
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to walk directory: %w", err)
	}

	// Sort by modification time (oldest first)
	sort.Slice(files, func(i, j int) bool {
		return files[i].modTime.Before(files[j].modTime)
	})

	log.Printf("Found %d files in %s", len(files), categoryName)

	// Collect all files with their assigned levels first
	type fileWithLevel struct {
		file  fileWithTime
		level int
	}
	var allFilesWithLevels []fileWithLevel

	maxLevels := 16 // Changed from 20 to 16
	for i, file := range files {
		// Distribute files evenly across 16 levels
		level := (i % maxLevels) + 1
		allFilesWithLevels = append(allFilesWithLevels, fileWithLevel{
			file:  file,
			level: level,
		})
	}

	// Process files
	for i, fwl := range allFilesWithLevels {
		level := fwl.level
		file := fwl.file

		// Use category offset + index to determine variant
		// This ensures different categories have different variants for the same level
		variant := (i + categoryOffset) % 3

		attack, defense, hp, price := calculateWeaponStats(level, variant)

		// Get Russian name for this level
		russianName := russianNames[level]
		if russianName == "" {
			russianName = fmt.Sprintf("Оружие_%d", level)
		}

		// Get English slug for filename
		englishSlug := englishSlugs[level]
		if englishSlug == "" {
			englishSlug = fmt.Sprintf("weapon-%d", level)
		}

		// Format: order-slug-price-attack-level-hp-defense.png
		newName := fmt.Sprintf("%d-%s-%d-%d-%d-%d-%d.png",
			level, englishSlug, price, attack, level, hp, defense)

		newPath := filepath.Join(filepath.Dir(file.path), newName)

		if err := os.Rename(file.path, newPath); err != nil {
			log.Printf("Error renaming %s to %s: %v", file.path, newPath, err)
			continue
		}

		variantName := []string{"Balanced", "High Attack", "High HP"}[variant]
		log.Printf("Renamed: %s -> %s (Level %d: %s, %s, Attack: %d, HP: %d, Price: %d)",
			filepath.Base(file.path), newName, level, russianName, variantName, attack, hp, price)
	}

	return nil
}

func main() {
	rand.Seed(time.Now().UnixNano())

	baseDir := "frontend/assets/images/equipment_items/weapon"

	categories := []struct {
		path string
		name string
	}{
		{filepath.Join(baseDir, "axes"), "axes"},
		{filepath.Join(baseDir, "maces"), "maces"},
		{filepath.Join(baseDir, "knifes"), "knifes"},
	}

	for _, cat := range categories {
		if _, err := os.Stat(cat.path); os.IsNotExist(err) {
			log.Printf("Directory does not exist: %s", cat.path)
			continue
		}

		if err := processWeaponCategory(cat.path, cat.name); err != nil {
			log.Printf("Error processing category %s: %v", cat.name, err)
		}
	}

	log.Println("Weapon image processing completed!")
}
