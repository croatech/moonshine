package database

import (
	"sunlight/config"

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
