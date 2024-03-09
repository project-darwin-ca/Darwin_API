package models

import (
	"gorm.io/gorm"
)

type CloudCredentials struct {
	UserId          string `gorm:"primaryKey" json:"user_id,omitempty"`
	Provider        string `json:"cloud_provider"`
	AccessKey       string `json:"access_key"`
	SecretAccessKey string `json:"secret_access_key"`
	DataMetadata    []DataMetadata
}

type DataMetadata struct {
	gorm.Model
	BucketName string `json:"bucket_name"`
	ObjectKey  string `json:"object_key"`
	UserId     string
}
