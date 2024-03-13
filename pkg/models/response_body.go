package models

import "github.com/google/uuid"

type CloudRegister struct {
	Provider string `json:"provider"`
	Identity string `json:"identity,omitempty"`
}

type ResponseBody struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
}

type Buckets struct {
	Buckets []string `json:"buckets"`
}

type Objects struct {
	BucketName string   `json:"bucket"`
	Objects    []string `json:"objects"`
}

type MountedFiles struct {
	DataId   uuid.UUID `json:"data_id"`
	FileName string    `json:"file_name"`
}
