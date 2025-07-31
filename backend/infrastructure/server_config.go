package infrastructure

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type Environment string

const (
	DevelopmentEnvironment Environment = "development"
	TestEnvironment        Environment = "test"
	ProductionEnvironment  Environment = "production"
)

var environmentMap = map[string]Environment{
	"development": DevelopmentEnvironment,
	"test":        TestEnvironment,
	"production":  ProductionEnvironment,
}

type ServerConfig struct {
	Version     string
	Environment Environment
	Port        int
}

func ParseServerConfig(version string) (ServerConfig, error) {
	if err := godotenv.Load(".env"); err != nil {
		if os.IsNotExist(err) {
			log.Info().Msg(".env file doesn't exist, skipping")
		} else {
			return ServerConfig{}, fmt.Errorf("loading .env file: %v", err)
		}
	}

	viper.AutomaticEnv()

	viper.SetDefault("LOG_LEVEL", "info")
	viper.SetDefault("PORT", 3000)

	levelStr := viper.GetString("LOG_LEVEL")
	logLevel, err := zerolog.ParseLevel(levelStr)
	if err != nil {
		log.Warn().Str("level", levelStr).Msg("invalid logger level, using default level 'info'")
		logLevel = zerolog.InfoLevel
	}

	zerolog.SetGlobalLevel(logLevel)

	environmentStr := viper.GetString("ENVIRONMENT")
	environment, ok := environmentMap[environmentStr]
	if !ok {
		return ServerConfig{}, fmt.Errorf("invalid environment %s", environmentStr)
	}

	cfg := ServerConfig{
		Version:     version,
		Port:        viper.GetInt("PORT"),
		Environment: environment,
	}

	if cfg.Environment == DevelopmentEnvironment {
		log.Logger = zerolog.New(zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
			NoColor:    false,
		}).With().Timestamp().Logger()
	}

	return cfg, nil
}
