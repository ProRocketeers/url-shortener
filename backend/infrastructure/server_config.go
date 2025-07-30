package infrastructure

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type ServerConfig struct {
	Version  string
	LogLevel string
	Port     int
}

func parseServerConfig(version string) (ServerConfig, error) {
	if err := godotenv.Load(".env"); err != nil {
		if os.IsNotExist(err) {
			fmt.Printf(".env file doesn't exist, skipping")
		} else {
			return ServerConfig{}, fmt.Errorf("loading .env file: %v", err)
		}
	}

	viper.AutomaticEnv()

	viper.SetDefault("LOG_LEVEL", "debug")
	viper.SetDefault("PORT", 3000)

	return ServerConfig{
		Version:  version,
		LogLevel: viper.GetString("LOG_LEVEL"),
		Port:     viper.GetInt("PORT"),
	}, nil
}
