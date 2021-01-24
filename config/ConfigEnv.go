package config

import (
	"time"

	"github.com/spf13/viper"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/microsoft"
)

var (
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

	//Token
	SavedToken *oauth2.Token

	//Auths
	OauthConf      *oauth2.Config
	OneDriveClient oauth2.TokenSource
)

func init() {
	// viper.SetConfigName("config")
	viper.SetConfigName("config-prod")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("../")

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	loadConfigVariables()
	createAuths()
}

func loadConfigVariables() {
	ClientID = viper.GetString("CLIENT_ID")
	ClientSecret = viper.GetString("CLIENT_SECRET")
	Scope = viper.GetString("SCOPE")
	RedirectUrl = viper.GetString("REDIRECT_URL")
	TenantID = viper.GetString("TENANT_ID")
	RefreshToken = viper.GetString("REFRESH_TOKEN")
	Expiry = viper.GetInt64("EXPIRY")
	TokenType = viper.GetString("TOKEN_TYPE")
	UploadFolderPath = viper.GetString("UPLOAD_FOLDER_PATH")
	RefreshAPIEndPoint = viper.GetString("REFESH_API_ENDPOINT")
	UploadAPIEndPoint = viper.GetString("UPLOAD_API_ENDPOINT")
	ShareAPIEndPoint = viper.GetString("SHARE_API_ENDPOINT")
	FragSize = viper.GetInt("FRAG_SIZE")
	ShareBodyJSON = viper.GetString("SHARE_BODY_JSON")
	UploadBodyJSON = viper.GetString("UPLOAD_BODY_JSON")

	ExpiredTime = time.Now().Add(3000 * time.Second)
}

func createAuths() {

	SavedToken = &oauth2.Token{
		AccessToken: "eyJ0eXAiOiJKV1QiLCJub25jZSI6IlFEbFQ5RzBnall0LTBjVy1rQ0NsLXhBcDNGbEtXdXdMclEtNzZxOHQxbUEiLCJhbGciOiJSUzI1NiIsIng1dCI6IjVPZjlQNUY5Z0NDd0NtRjJCT0hIeEREUS1EayIsImtpZCI6IjVPZjlQNUY5Z0NDd0NtRjJCT0hIeEREUS1EayJ9.eyJhdWQiOiIwMDAwMDAwMy0wMDAwLTAwMDAtYzAwMC0wMDAwMDAwMDAwMDAiLCJpc3MiOiJodHRwczovL3N0cy53aW5kb3dzLm5ldC9mYTM4OTRiYS0wZmZjLTQ4NDctYWU5NC01MTAxNDZkN2NhOTcvIiwiaWF0IjoxNjA5NjYyODA5LCJuYmYiOjE2MDk2NjI4MDksImV4cCI6MTYwOTY2NjcwOSwiYWNjdCI6MCwiYWNyIjoiMSIsImFjcnMiOlsidXJuOnVzZXI6cmVnaXN0ZXJzZWN1cml0eWluZm8iLCJ1cm46bWljcm9zb2Z0OnJlcTEiLCJ1cm46bWljcm9zb2Z0OnJlcTIiLCJ1cm46bWljcm9zb2Z0OnJlcTMiLCJjMSIsImMyIiwiYzMiLCJjNCIsImM1IiwiYzYiLCJjNyIsImM4IiwiYzkiLCJjMTAiLCJjMTEiLCJjMTIiLCJjMTMiLCJjMTQiLCJjMTUiLCJjMTYiLCJjMTciLCJjMTgiLCJjMTkiLCJjMjAiLCJjMjEiLCJjMjIiLCJjMjMiLCJjMjQiLCJjMjUiXSwiYWlvIjoiRTJKZ1lGajl6N2J5aXRvRXJ5aC96bThuclA4NXg0cnNZRGhjN0hsSjdON0MwenVxV2dzQiIsImFtciI6WyJwd2QiXSwiYXBwX2Rpc3BsYXluYW1lIjoiZGlydHlmaWxtIGFwcCIsImFwcGlkIjoiZmJlMWZmZDEtOTNjYS00YWFmLWExMjEtNjU2ODQ5YjJjZmQzIiwiYXBwaWRhY3IiOiIxIiwiaWR0eXAiOiJ1c2VyIiwiaXBhZGRyIjoiNDIuMTE5LjYxLjE4MiIsIm5hbWUiOiJOZ3V5ZW4gSGFvIiwib2lkIjoiMjczN2NlNzEtZGRmYy00MWQyLTkwMjgtMDE4Y2U1NGExOGY0IiwicGxhdGYiOiI1IiwicHVpZCI6IjEwMDMyMDAxMDgxMTAzNkEiLCJyaCI6IjAuQUFBQXVwUTQtdndQUjBpdWxGRUJSdGZLbDlIXzRmdktrNjlLb1NGbGFFbXl6OU5KQUJFLiIsInNjcCI6IkZpbGVzLlJlYWRXcml0ZS5BbGwgb3BlbmlkIFVzZXIuUmVhZEJhc2ljLkFsbCBVc2VyLlJlYWRXcml0ZSBwcm9maWxlIGVtYWlsIiwic2lnbmluX3N0YXRlIjpbImttc2kiXSwic3ViIjoiRzkwcG5uUmJTM2U3MUV5MnpSUGctaldkV0lXbWxkelI4eV9mWm9MRkNVZyIsInRlbmFudF9yZWdpb25fc2NvcGUiOiJBUyIsInRpZCI6ImZhMzg5NGJhLTBmZmMtNDg0Ny1hZTk0LTUxMDE0NmQ3Y2E5NyIsInVuaXF1ZV9uYW1lIjoiZGF2aWRudGgxMjE3MUBtb2Qub2JhZ2cuY29tIiwidXBuIjoiZGF2aWRudGgxMjE3MUBtb2Qub2JhZ2cuY29tIiwidXRpIjoiLXVCMUxYcmh6RTJPOWV5NDZadUpBQSIsInZlciI6IjEuMCIsIndpZHMiOlsiYjc5ZmJmNGQtM2VmOS00Njg5LTgxNDMtNzZiMTk0ZTg1NTA5Il0sInhtc19zdCI6eyJzdWIiOiJiaXJUQWVfMGt1TllNYU03RU5Ib2JBdi1MVkJYUmtPUElYeE9KdlJWMGEwIn0sInhtc190Y2R0IjoxNTcxNDIxMDAzfQ.HNqqNFDzzTnkan8dSbD1Zw0O4DlmjPS7P5AD98Y6wtBRSNTGqSr1qGO-f6D9HyQPRnOkxoKJbJVTi1PesoG-BxZF1lpXd9s_Yc0pN5ckqPZnIXBsk4psg8zWrThuKTXQEcHdsDiq1MheZpxuT9763D3EmFIghTcbiyPtf3eTK3_nuTMxIKRlps4IGUUxb85pjqYeQ4f6_ikAt5hbL8l4SOMisTPtnYdyNnKW_poWMWEP158L90zDtXk6Tu5DC1P49zotmug6x8SYDkAhgeFGvXFJ_uCpBw-gLNvWs-FdGLLbXQILM_qNlRQpmeq0eofvUAVqgnA-UsjvPDXoqTxe5w",
		// AccessToken:  "",
		RefreshToken: RefreshToken,
		Expiry:       time.Now().Add(time.Duration(viper.GetInt64("EXPIRY")) * time.Second),
		TokenType:    TokenType,
	}

	OauthConf = &oauth2.Config{
		ClientID:     ClientID,
		ClientSecret: ClientSecret,
		RedirectURL:  RedirectUrl,
		Scopes:       []string{"Files.ReadWrite.All", "Sites.ReadWrite.All", "openid", "User.ReadBasic.All", "User.ReadWrite", "profile", "email"},
		Endpoint:     microsoft.AzureADEndpoint(TenantID),
	}

	OneDriveClient = oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: SavedToken.AccessToken},
	)
}
