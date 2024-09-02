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
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Println("No .env file found, using environment variables instead")
		} else {
			log.Fatalf("Error reading .env file: %s", err)
		}
	}
}
