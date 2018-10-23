package database

import (
	"sunlight/config"
	"sunlight/models"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/spf13/viper"
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

var connection *gorm.DB

func Prepare() error {
	conn := connect()

	migrate(conn)

	defer conn.Close()

	return nil
}

func Connection() *gorm.DB {
	conn := connect()
	conn.LogMode(true)

	return conn
}

func connect() *gorm.DB {
	config.Load()

	conf := Config{
		Adapter:  viper.GetString("database.adapter"),
		Host:     viper.GetString("database.host"),
		Port:     viper.GetString("database.port"),
		User:     viper.GetString("database.user"),
		Database: viper.GetString("database.database"),
		Password: viper.GetString("database.password"),
		SSLMode:  viper.GetString("database.ssl_mode"),
	}

	conn, err := gorm.Open(conf.Adapter, "host="+conf.Host+" port="+conf.Port+" user="+conf.User+" dbname="+conf.Database+" password="+conf.Password+" sslmode="+conf.SSLMode+"")

	if err != nil {
		panic(err)
	}

	return conn
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
