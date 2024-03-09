package models

type RegisterCloudProviderRequest struct {
	Provider        string `json:"cloud_provider"`
	AccessKey       string `json:"access_key"`
	SecretAccessKey string `json:"secret_access_key"`
}
