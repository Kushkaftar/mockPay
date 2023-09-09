package config

import (
	"log"
	"mockPay/internal/pkg/models"

	"github.com/spf13/viper"
)

func MustConfig(fileName, directory string) *models.Config {
	var config models.Config

	if err := initConfig(fileName, directory); err != nil {
		log.Fatalf("crush init config, %v", err)
	}

	err := viper.Unmarshal(&config)
	if err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
	}

	return &config
}

func initConfig(fileName, directory string) error {
	viper.AddConfigPath(directory)
	viper.SetConfigName(fileName)
	viper.SetConfigType("yaml")
	return viper.ReadInConfig()
}
