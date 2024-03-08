package db

import (
	"github.com/xuxant/Darwin_API/pkg/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

func GetConnection() *gorm.DB {
	cfg := config.GetConfig()
	db, err := gorm.Open(postgres.Open(cfg.PostgresConnectionString), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
	}
	return db
}
