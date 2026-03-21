//go:build gorm

package main

import (
	"fmt"
	"os"

	"github.com/ProRocketeers/url-shortener/domain/model"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gen"
	"gorm.io/gorm"
)

func main() {
	g := gen.NewGenerator(gen.Config{
		OutPath:       "./domain/query",
		Mode:          gen.WithDefaultQuery,
		FieldNullable: true,
	})

	cfg, err := parseDbConfig()
	if err != nil {
		fmt.Printf("error parsing DB config: %v", err)
		os.Exit(1)
	}
	connectionString := fmt.Sprintf(
		"host=%v user=%v password=%v dbname=%v port=%v sslmode=disable TimeZone=UTC",
		cfg.Host,
		cfg.User,
		cfg.Password,
		cfg.Database,
		cfg.Port,
	)

	db, err := gorm.Open(postgres.Open(connectionString), &gorm.Config{})

	if err != nil {
		fmt.Printf("error opening connection to db: %v", err)
		os.Exit(1)
	}

	g.UseDB(db)
	g.ApplyBasic(model.ShortLink{}, model.RequestInfo{})
	g.Execute()
}

type databaseConfig struct {
	Host     string
	User     string
	Password string
	Database string
	Port     int
}

func parseDbConfig() (databaseConfig, error) {
	if err := godotenv.Load(".env"); err != nil {
		if os.IsNotExist(err) {
			fmt.Println(".env file doesn't exist, skipping")
		} else {
			return databaseConfig{}, fmt.Errorf("loading .env file error: %v", err)
		}
	}

	// configures Viper to automatically try to fetch keys from env vars
	viper.AutomaticEnv()

	return databaseConfig{
		Host:     viper.GetString("DB_HOST"),
		User:     viper.GetString("DB_USER"),
		Password: viper.GetString("DB_PASSWORD"),
		Database: viper.GetString("DB_NAME"),
		Port:     viper.GetInt("DB_PORT"),
	}, nil
}
