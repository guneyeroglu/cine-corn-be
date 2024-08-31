package config

import (
	"github.com/spf13/viper"
)

func Init() {
	viper.SetConfigFile(".env")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		viper.AutomaticEnv()
	}

}
