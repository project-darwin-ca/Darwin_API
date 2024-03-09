package db

import (
	lg "github.com/rs/zerolog/log"
	"github.com/xuxant/Darwin_API/pkg/config"
	"github.com/xuxant/Darwin_API/pkg/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

func GetConnection(cfg config.ServiceConfig) *gorm.DB {
	db, err := gorm.Open(postgres.Open(cfg.PostgresConnectionString), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
	}
	return db
}

func MigrateDatabase(db *gorm.DB) error {
	lg.Info().Msg("Initialing migration...")
	lg.Info().Msg("Migrating metadata")
	_ = db.AutoMigrate(&models.CloudCredentials{}, &models.DataMetadata{})
	lg.Info().Msg("Migrating data source")

	return nil
}
