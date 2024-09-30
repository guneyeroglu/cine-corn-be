package config

import (
	"log"

	"github.com/spf13/viper"
)

func Init() {
	viper.SetConfigFile(".env")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		log.Println("Configuration Error: .env file not found. Falling back to environment variables.")
		viper.AutomaticEnv()
	}
}
