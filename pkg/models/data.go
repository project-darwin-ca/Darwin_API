package models

import "github.com/google/uuid"

type CloudCredentials struct {
	UserId          string         `gorm:"primaryKey;default:null" json:"user_id,omitempty"`
	Provider        string         `json:"cloud_provider"`
	AccessKey       string         `json:"access_key"`
	SecretAccessKey string         `json:"secret_access_key"`
	DataMetadata    []DataMetadata `gorm:"foreignKey:UserId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

type DataMetadata struct {
	DataId     uuid.UUID `gron:"primaryKey"`
	BucketName string
	ObjectKey  string
	FileName   string
	UserId     string
}
