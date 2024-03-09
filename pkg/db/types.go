package db

import "gorm.io/gorm"

type DBConnection struct {
	Db *gorm.DB
}
