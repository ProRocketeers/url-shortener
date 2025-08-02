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

type environment string

const (
	DevelopmentEnvironment environment = "development"
	TestEnvironment        environment = "test"
	ProductionEnvironment  environment = "production"
)

var environmentMap = map[string]environment{
	"development": DevelopmentEnvironment,
	"test":        TestEnvironment,
	"production":  ProductionEnvironment,
}

type serverConfigMeta struct {
	Version    string
	CommitHash string
	BuildTime  string
}

type databaseConfig struct {
	Host     string
	User     string
	Password string
	Database string
	Port     int
}

type domainConfig struct {
	ExpiredLinkCleanupInterval time.Duration
	BaseUrl                    string // including the protocol and possibly port
}

type Config struct {
	Metadata    serverConfigMeta
	Environment environment
	Port        int
	Database    databaseConfig
	Domain      domainConfig
}

func InitialLoggerConfig() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	zerolog.TimeFieldFormat = time.RFC3339
	zerolog.TimestampFunc = func() time.Time {
		return time.Now().UTC()
	}
}

func ParseServerConfig(version, commitHash, buildTime string) (Config, error) {
	if err := godotenv.Load(".env"); err != nil {
		if os.IsNotExist(err) {
			log.Info().Msg(".env file doesn't exist, skipping")
		} else {
			return Config{}, fmt.Errorf("loading .env file: %v", err)
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
		return Config{}, fmt.Errorf("invalid environment %s", environmentStr)
	}

	cfg := Config{
		Metadata: serverConfigMeta{
			version,
			commitHash,
			buildTime,
		},
		Port:        viper.GetInt("PORT"),
		Environment: environment,
		Database: databaseConfig{
			Host:     viper.GetString("DB_HOST"),
			User:     viper.GetString("DB_USER"),
			Password: viper.GetString("DB_PASSWORD"),
			Database: viper.GetString("DB_DATABASE"),
			Port:     viper.GetInt("DB_PORT"),
		},
		Domain: domainConfig{
			BaseUrl:                    viper.GetString("BASE_URL"),
			ExpiredLinkCleanupInterval: viper.GetDuration("EXPIRED_LINK_CLEANUP_INTERVAL"),
		},
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
