package models

type Config struct {
	ClientID                string `json:"clientID"`
	ClientSecret            string `json:"clientSecret"`
	AccessToken             string `json:"accessToken"`
	Scope                   string `json:"scope,omitempty"`
	RedirectUrl             string `json:"redirectURL"`
	TenantID                string `json:"tenantID"`
	RefreshToken            string `json:"refreshToken"`
	Expiry                  int64  `json:"expiry,omitempty"`
	TokenType               string `json:"tokenType,omitempty"`
	UploadFolderPath        string `json:"uploadFolderPath,omitempty"`
	RefreshAPIEndPoint      string `json:"refreshAPIEndpoint,omitempty"`
	UploadAPIEndPoint       string `json:"uploadAPIEndpoint,omitempty"`
	ShareAPIEndPoint        string `json:"shareAPIEndPoint,omitempty"`
	CreateFolderAPIEndPoint string `json:"createFolderAPIEndPoint,omitempty"`
	GetItemsPathEndPoint    string `json:"getItemsPathEndPoint,omitempty"`
	GetItemAPIEndPoint      string `json:"getItemAPIEndPoint,omitempty"`
	DeleteItemAPIEndPoint   string `json:"deleteItemAPIEndPoint,omitempty"`
	MoveItemAPIEndPoint     string `json:"moveItemAPIEndPoint,omitempty"`
	CopyItemAPIEndPoint     string `json:"copyItemAPIEndPoint,omitempty"`
	DownloadItemAPIEndPoint string `json:"downloadItemAPIEndPoint,omitempty"`

	FragSize             int    `json:"fragSize,omitempty"`
	ShareBodyJSON        string `json:"shareBodyJSON,omitempty"`
	UploadBodyJSON       string `json:"uploadBodyJSON,omitempty"`
	CreateFolderBodyJSON string `json:"createFolderBodyJSON,omitempty"`
	MoveItemBodyJSON     string `json:"moveItemBodyJSON,omitempty"`
	CopyItemBodyJSON     string `json:"copyItemBodyJSON,omitempty"`
}
