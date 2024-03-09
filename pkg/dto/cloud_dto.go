package dto

import (
	"github.com/xuxant/Darwin_API/pkg/models"
)

func (s *Service) AddCloudCredentials(d models.CloudCredentials) error {
	result := s.Client.Create(&d)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
