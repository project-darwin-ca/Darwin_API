package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DataSource struct {
	gorm.Model
	UserId    uuid.UUID
	Provider  string
	Source    string
	SecretKey string
	SecretID  string
}

type DataMetadata struct {
	gorm.Model
	UserID     uuid.UUID
	BucketName string
	ObjectKey  string
}

type RegisterCloudProvider struct {
	Provider        string `json:"cloud_provider"`
	AccessKey       string `json:"access_key"`
	SecretAccessKey string `json:"secret_access_key"`
}
