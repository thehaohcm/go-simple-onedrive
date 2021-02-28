package service

import (
	"time"

	"golang.org/x/oauth2"
)

type Service struct {
	ClientID                string
	ClientSecret            string
	AccessToken             string
	Scope                   string
	RedirectUrl             string
	TenantID                string
	RefreshToken            string
	Expiry                  int64
	TokenType               string
	UploadFolderPath        string
	RefreshAPIEndPoint      string
	UploadAPIEndPoint       string
	ShareAPIEndPoint        string
	CreateFolderAPIEndPoint string
	FragSize                int
	GetItemsPathEndPoint    string
	GetItemAPIEndPoint      string
	DeleteItemAPIEndPoint   string
	MoveItemAPIEndPoint     string
	CopyItemAPIEndPoint     string
	DownloadItemAPIEndPoint string

	ShareBodyJSON        string
	UploadBodyJSON       string
	CreateFolderBodyJSON string
	MoveItemBodyJSON     string
	CopyItemBodyJSON     string

	ExpiredTime time.Time

	//Token
	SavedToken *oauth2.Token

	//Auths
	OauthConf *oauth2.Config
}
