package infrastructure

import (
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	gormlogger "gorm.io/gorm/logger"
)

type environment string

const (
	DevelopmentEnvironment environment = "development"
	TestEnvironment        environment = "test"
	ProductionEnvironment  environment = "production"
)

var environmentMap = map[string]environment{
	string(DevelopmentEnvironment): DevelopmentEnvironment,
	string(TestEnvironment):        TestEnvironment,
	string(ProductionEnvironment):  ProductionEnvironment,
}

var databaseLogLevelMap = map[string]gormlogger.LogLevel{
	"silent": gormlogger.Silent,
	"error":  gormlogger.Error,
	"warn":   gormlogger.Warn,
	"info":   gormlogger.Info,
}

type versionMetadata struct {
	Version    string
	CommitHash string
	BuildTime  string
}

type databaseConfig struct {
	Host               string
	User               string
	Password           string
	Database           string
	Port               int
	LogLevel           gormlogger.LogLevel
	SlowQueryThreshold time.Duration
}

type domainConfig struct {
	ExpiredLinkCleanupInterval time.Duration
	BaseUrl                    url.URL
}

type Config struct {
	Metadata    versionMetadata
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

	// configures Viper to automatically try to fetch keys from env vars
	viper.AutomaticEnv()

	setEnvDefaults()

	zerolog.SetGlobalLevel(parseLogLevel())

	environmentStr := strings.ToLower(strings.TrimSpace(viper.GetString("ENVIRONMENT")))
	environment, ok := environmentMap[environmentStr]
	if !ok {
		return Config{}, fmt.Errorf("invalid environment %s", environmentStr)
	}

	baseUrlStr := viper.GetString("BASE_URL")
	if baseUrlStr == "" {
		return Config{}, fmt.Errorf("base URL not specified")
	}
	baseUrl, err := url.Parse(baseUrlStr)
	if err != nil {
		return Config{}, fmt.Errorf("parsing base URL: %v", err)
	}

	cfg := Config{
		Metadata: versionMetadata{
			version,
			commitHash,
			buildTime,
		},
		Port:        viper.GetInt("PORT"),
		Environment: environment,
		Database: databaseConfig{
			Host:               viper.GetString("DB_HOST"),
			User:               viper.GetString("DB_USER"),
			Password:           viper.GetString("DB_PASSWORD"),
			Database:           viper.GetString("DB_NAME"),
			Port:               viper.GetInt("DB_PORT"),
			LogLevel:           parseDbLogLevel(),
			SlowQueryThreshold: viper.GetDuration("DB_SLOW_QUERY_THRESHOLD"),
		},
		Domain: domainConfig{
			BaseUrl:                    *baseUrl,
			ExpiredLinkCleanupInterval: cleanupInterval,
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

func setEnvDefaults() {
	viper.SetDefault("ENVIRONMENT", string(ProductionEnvironment))
	viper.SetDefault("LOG_LEVEL", "info")
	viper.SetDefault("PORT", 8080)
	viper.SetDefault("DB_LOG_LEVEL", "info")
	viper.SetDefault("DB_SLOW_QUERY_THRESHOLD", "250ms")
	viper.SetDefault("EXPIRED_LINK_CLEANUP_INTERVAL", "2m")
}

func parseLogLevel() zerolog.Level {
	levelStr := viper.GetString("LOG_LEVEL")
	logLevel, err := zerolog.ParseLevel(levelStr)
	if err != nil {
		log.Warn().Str("level", levelStr).Msg("invalid logger level, using default level 'info'")
		logLevel = zerolog.InfoLevel
	}
	return logLevel
}

func parseDbLogLevel() gormlogger.LogLevel {
	dbLogLevelStr := viper.GetString("DB_LOG_LEVEL")
	dbLogLevel, ok := databaseLogLevelMap[dbLogLevelStr]
	if !ok {
		log.Warn().Str("level", dbLogLevelStr).Msg("invalid database logger level, using default level 'info'")
		dbLogLevel = gormlogger.Info
	}
	return dbLogLevel
}
