package config

import (
	"fmt"

	"github.com/spf13/viper"
)

func Load() {
	viper.SetConfigFile("./config/config.yml")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
}
