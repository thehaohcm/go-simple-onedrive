package models

import(
	"time"
)

type Config struct{
	ClientID           string
	ClientSecret       string
	Scope              string
	RedirectUrl        string
	TenantID           string
	RefreshToken       string
	Expiry             int64
	TokenType          string
	UploadFolderPath   string
	RefreshAPIEndPoint string
	UploadAPIEndPoint  string
	ShareAPIEndPoint   string
	FragSize           int
	ShareBodyJSON      string
	UploadBodyJSON     string

	ExpiredTime time.Time
}