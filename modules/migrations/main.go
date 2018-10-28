package migrations

import (
	"sunlight/models"
	"sunlight/modules/database"

	"github.com/jinzhu/gorm"
)

func Run() error {
	conn := database.Connection()

	migrate(conn)

	defer conn.Close()

	return nil
}

func migrate(db *gorm.DB) {
	db.AutoMigrate(&models.Avatar{})
	db.AutoMigrate(&models.Bot{})
	db.AutoMigrate(&models.EquipmentCategory{})
	db.AutoMigrate(&models.EquipmentItem{})
	db.AutoMigrate(&models.Event{})
	db.AutoMigrate(&models.Fight{})
	db.AutoMigrate(&models.LocationBot{})
	db.AutoMigrate(&models.LocationLocation{})
	db.AutoMigrate(&models.LocationResource{})
	db.AutoMigrate(&models.Location{})
	db.AutoMigrate(&models.Message{})
	db.AutoMigrate(&models.Movement{})
	db.AutoMigrate(&models.Resource{})
	db.AutoMigrate(&models.Round{})
	db.AutoMigrate(&models.Stuff{})
	db.AutoMigrate(&models.ToolCategory{})
	db.AutoMigrate(&models.ToolItem{})
	db.AutoMigrate(&models.User{})
}
