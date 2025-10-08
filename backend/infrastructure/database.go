package infrastructure

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectToDatabase(config Config) (*gorm.DB, error) {
	connectionString := fmt.Sprintf(
		"host=%v user=%v password=%v dbname=%v port=%v TimeZone=UTC",
		config.Database.Host,
		config.Database.User,
		config.Database.Password,
		config.Database.Database,
		config.Database.Port,
	)
	return gorm.Open(postgres.Open(connectionString), &gorm.Config{
		Logger: NewZerologGormLogger(log.Logger, config.Database.LogLevel, config.Database.SlowQueryThreshold),
	})
}
