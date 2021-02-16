package models

type UploadBlockResponse struct {
	ExpirationDateTime string   `json:"expirationDateTime,omitempty"`
	NextExpectedRanges []string `json:"nextExpectedRanges,omitempty"`
}
