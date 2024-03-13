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

func (s *Service) FetchCloudCredentials(user string) (models.CloudCredentials, error) {
	creds := models.CloudCredentials{
		UserId: user,
	}
	result := s.Client.First(&creds)
	if result.Error != nil {
		return creds, result.Error
	}

	return creds, nil
}

func (s *Service) AddS3Objects(d []models.DataMetadata) error {
	result := s.Client.Create(&d)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (s *Service) FetchMountedFiles(user string) ([]models.DataMetadata, error) {
	var files []models.DataMetadata

	result := s.Client.Where("user_id = ?", user).Find(&files)
	if result.Error != nil {
		return files, result.Error
	}
	return files, nil
}

func (s *Service) UpdateCredentials(user string, accessKey string, secretKey string) error {
	var creds models.CloudCredentials
	result := s.Client.First(&creds, "user_id = ?", user)
	if result.Error != nil {
		return result.Error
	}
	creds.AccessKey = accessKey
	creds.SecretAccessKey = secretKey
	result = s.Client.Save(&creds)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
