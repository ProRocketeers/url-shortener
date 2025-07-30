package infrastructure

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type ServerConfig struct {
	LogLevel string
}

func ParseServerConfig() (ServerConfig, error) {
	if err := godotenv.Load(".env"); err != nil {
		if os.IsNotExist(err) {
			fmt.Printf(".env file doesn't exist, skipping")
		} else {
			return ServerConfig{}, fmt.Errorf("loading .env file: %v", err)
		}
	}

	viper.AutomaticEnv()

	return ServerConfig{
		LogLevel: viper.GetString("LOG_LEVEL"),
	}, nil
}
