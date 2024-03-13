package models

type AddS3Object struct {
	BucketName string   `json:"bucketName"`
	Objects    []string `json:"files"`
}
