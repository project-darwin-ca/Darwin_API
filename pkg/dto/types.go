package dto

import "gorm.io/gorm"

type Service struct {
	Client *gorm.DB
}

func NewDtoService(client *gorm.DB) *Service {
	return &Service{
		client,
	}
}
