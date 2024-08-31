package config

import (
	"log"

	"github.com/spf13/viper"
)

func Init() {
	viper.SetConfigFile(".env")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error reading .env file: %s", err)
	}

}
