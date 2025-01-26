package database

import (
	"moonshine/models"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type Config struct {
	Adapter  string
	Host     string
	User     string
	Port     string
	Database string
	Password string
	SSLMode  string
}

func Connection() *gorm.DB {
	conn := connect()
	conn.LogMode(true)

	return conn
}

func connect() *gorm.DB {
	conf := Config{
		Adapter:  os.Getenv("DATABASE_ADAPTER"),
		Host:     os.Getenv("DATABASE_HOST"),
		Port:     os.Getenv("DATABASE_PORT"),
		User:     os.Getenv("DATABASE_USER"),
		Database: os.Getenv("DATABASE_NAME"),
		Password: os.Getenv("DATABASE_PASSWORD"),
		SSLMode:  os.Getenv("DATABASE_SSL_MODE"),
	}

	conn, err := gorm.Open(conf.Adapter, "host="+conf.Host+" port="+conf.Port+" user="+conf.User+" dbname="+conf.Database+" password="+conf.Password+" sslmode="+conf.SSLMode+"")

	if err != nil {
		panic(err)
	}

	return conn
}

func Migrate() {
	db := connect()

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

	defer db.Close()
}

func Drop() {
	db := connect()

	db.DropTableIfExists(&models.Avatar{})
	db.DropTableIfExists(&models.Bot{})
	db.DropTableIfExists(&models.EquipmentCategory{})
	db.DropTableIfExists(&models.EquipmentItem{})
	db.DropTableIfExists(&models.Event{})
	db.DropTableIfExists(&models.Fight{})
	db.DropTableIfExists(&models.LocationBot{})
	db.DropTableIfExists(&models.LocationLocation{})
	db.DropTableIfExists(&models.LocationResource{})
	db.DropTableIfExists(&models.Location{})
	db.DropTableIfExists(&models.Message{})
	db.DropTableIfExists(&models.Movement{})
	db.DropTableIfExists(&models.Resource{})
	db.DropTableIfExists(&models.Round{})
	db.DropTableIfExists(&models.Stuff{})
	db.DropTableIfExists(&models.ToolCategory{})
	db.DropTableIfExists(&models.ToolItem{})
	db.DropTableIfExists(&models.User{})

	defer db.Close()
}

func Clean() {
	db := connect()

	db.Exec("DELETE FROM avatars")
	db.Exec("DELETE FROM bots")
	db.Exec("DELETE FROM equipment_categories")
	db.Exec("DELETE FROM equipment_items")
	db.Exec("DELETE FROM events")
	db.Exec("DELETE FROM fights")
	db.Exec("DELETE FROM location_bots")
	db.Exec("DELETE FROM location_locations")
	db.Exec("DELETE FROM location_resources")
	db.Exec("DELETE FROM locations")
	db.Exec("DELETE FROM messages")
	db.Exec("DELETE FROM movements")
	db.Exec("DELETE FROM resources")
	db.Exec("DELETE FROM rounds")
	db.Exec("DELETE FROM stuffs")
	db.Exec("DELETE FROM tool_categories")
	db.Exec("DELETE FROM tool_items")
	db.Exec("DELETE FROM users")

	defer db.Close()
}
