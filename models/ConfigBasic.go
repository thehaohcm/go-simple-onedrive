package models

import "time"

type ConfigBasic struct {
	ClientID         string
	ClientSecret     string
	Scope            string
	RedirectUrl      string
	TenantID         string
	RefreshToken     string
	Expiry           int64
	TokenType        string
	UploadFolderPath string
	ShareBodyJSON    string
	UploadBodyJSON   string

	ExpiredTime time.Time
}
