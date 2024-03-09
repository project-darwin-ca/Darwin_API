package db

import (
	"fmt"
	lg "github.com/rs/zerolog/log"
	"github.com/xuxant/Darwin_API/pkg/config"
	"github.com/xuxant/Darwin_API/pkg/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

var Db DBConnection

func GetConnection(cfg config.ServiceConfig) DBConnection {
	db, err := gorm.Open(postgres.Open(cfg.PostgresConnectionString), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
	}
	return DBConnection{
		Db: db,
	}
}

func (db *DBConnection) MigrateDatabase() error {
	lg.Info().Msg("Initialing migration...")
	lg.Info().Msg("Migrating metadata")
	err := db.Db.AutoMigrate(&models.DataMetadata{})
	if err != nil {
		return fmt.Errorf("migrating table datametadata: %w", err)
	}
	lg.Info().Msg("Migrating data source")
	err = db.Db.AutoMigrate(&models.DataSource{})
	if err != nil {
		return fmt.Errorf("migrating table datasource: %w", err)
	}

	return nil
}
