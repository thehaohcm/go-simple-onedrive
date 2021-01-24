package models

type UploadSessionResponse struct {
	UploadUrl          string `json:"uploadUrl"`
	ExpirationDateTime string `json:"expirationDateTime"`
}
