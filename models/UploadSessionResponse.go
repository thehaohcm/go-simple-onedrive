package models

type UploadSessionResponse struct {
	UploadUrl          string `json:"uploadUrl,omitempty"`
	ExpirationDateTime string `json:"expirationDateTime,omitempty"`
}
