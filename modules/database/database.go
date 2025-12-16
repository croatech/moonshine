package database

import (
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/lib/pq"
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
	adapter := os.Getenv("DATABASE_ADAPTER")
	if adapter == "" {
		adapter = "postgres"
	}

	conf := Config{
		Adapter:  adapter,
		Host:     os.Getenv("DATABASE_HOST"),
		Port:     os.Getenv("DATABASE_PORT"),
		User:     os.Getenv("DATABASE_USER"),
		Database: os.Getenv("DATABASE_NAME"),
		Password: os.Getenv("DATABASE_PASSWORD"),
		SSLMode:  os.Getenv("DATABASE_SSL_MODE"),
	}

	if conf.Host == "" {
		conf.Host = "localhost"
	}
	if conf.Port == "" {
		conf.Port = "5433"
	}
	if conf.User == "" {
		conf.User = "postgres"
	}
	if conf.Database == "" {
		conf.Database = "moonshine"
	}
	if conf.Password == "" {
		conf.Password = "postgres"
	}
	if conf.SSLMode == "" {
		conf.SSLMode = "disable"
	}

	dsn := "host=" + conf.Host + " port=" + conf.Port + " user=" + conf.User + " dbname=" + conf.Database + " password=" + conf.Password + " sslmode=" + conf.SSLMode

	conn, err := gorm.Open(conf.Adapter, dsn)

	if err != nil {
		panic(err)
	}

	return conn
}
