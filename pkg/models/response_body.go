package models

type CloudRegister struct {
	Provider string `json:"provider"`
	Identity string `json:"identity,omitempty"`
}

type ResponseBody struct {
	Success bool          `json:"success"`
	Data    CloudRegister `json:"data"`
}
