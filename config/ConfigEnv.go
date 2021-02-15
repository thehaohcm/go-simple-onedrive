package config

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/spf13/viper"
	"github.com/thehaohcm/go-simple-onedrive/models"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/microsoft"
)

var (
	ClientID             string
	ClientSecret         string
	AccessToken          string
	Scope                string
	RedirectUrl          string
	TenantID             string
	RefreshToken         string
	Expiry               int64
	TokenType            string
	UploadFolderPath     string
	RefreshAPIEndPoint   string
	UploadAPIEndPoint    string
	ShareAPIEndPoint     string
	FragSize             int
	ShareBodyJSON        string
	UploadBodyJSON       string
	GetItemsPathEndPoint string

	ExpiredTime time.Time

	//Token
	SavedToken *oauth2.Token

	//Auths
	OauthConf *oauth2.Config
)

func LoadConfigFromFile() {
	viper.SetConfigName("config")
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
	AccessToken = viper.GetString("ACCESS_TOKEN")
	AccessToken = viper.GetString("ACCESS_TOKEN")
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
	GetItemsPathEndPoint = viper.GetString("GET_ITEMS_PATH_ENDPOINT")

	ExpiredTime = time.Now().Add(3000 * time.Second)
}

func createAuths() {

	SavedToken = &oauth2.Token{
		AccessToken:  AccessToken,
		RefreshToken: RefreshToken,
		Expiry:       time.Now().Add(time.Duration(viper.GetInt64("EXPIRY")) * time.Second),
		TokenType:    TokenType,
	}

	OauthConf = &oauth2.Config{
		ClientID:     ClientID,
		ClientSecret: ClientSecret,
		RedirectURL:  RedirectUrl,
		// Scopes:       []string{"Files.ReadWrite.All", "Sites.ReadWrite.All", "openid", "User.ReadBasic.All", "User.ReadWrite", "profile", "email"},
		Scopes:   strings.Fields(Scope),
		Endpoint: microsoft.AzureADEndpoint(TenantID),
	}

	//refresh Token
	RefreshTokenFunc()
}

func LoadConfigFromJson(config *models.Config) {
	ClientID = config.ClientID
	ClientSecret = config.ClientSecret
	AccessToken = config.AccessToken
	Scope = config.Scope
	if Scope == "" {
		Scope = "Files.ReadWrite.All openid User.ReadBasic.All User.ReadWrite profile email"
	}
	RedirectUrl = config.RedirectUrl
	TenantID = config.TenantID
	RefreshToken = config.RefreshToken
	Expiry = config.Expiry
	if Expiry == 0 {
		Expiry = 3599
	}
	TokenType = config.TokenType
	if len(TokenType) == 0 {
		TokenType = "Bearer"
	}
	UploadFolderPath = config.UploadFolderPath
	if len(UploadFolderPath) == 0 {
		UploadFolderPath = "/"
	}
	RefreshAPIEndPoint = config.RefreshAPIEndPoint
	if len(RefreshAPIEndPoint) == 0 {
		RefreshAPIEndPoint = "https://login.microsoftonline.com/{TENANT_ID}/oauth2/v2.0/token"
	}
	UploadAPIEndPoint = config.UploadAPIEndPoint
	if len(UploadAPIEndPoint) == 0 {
		UploadAPIEndPoint = "https://graph.microsoft.com/v1.0/me/drive/root:{UPLOAD_FOLDER_PATH}/{FILE_NAME}:/createUploadSession"
	}
	ShareAPIEndPoint = config.ShareAPIEndPoint
	if len(ShareAPIEndPoint) == 0 {
		ShareAPIEndPoint = "https://graph.microsoft.com/v1.0/me/drive/items/{UPLOADED_FILE_ID}/createLink"
	}
	FragSize = config.FragSize
	if FragSize == 0 {
		FragSize = 62259200
	}
	ShareBodyJSON = config.ShareBodyJSON
	if len(ShareBodyJSON) == 0 {
		ShareBodyJSON = "{ \"type\": \"view\", \"scope\": \"anonymous\" }"
	}
	UploadBodyJSON = config.UploadBodyJSON
	if len(UploadBodyJSON) == 0 {
		UploadBodyJSON = "{\"item\":{\"@microsoft.graph.conflictBehavior\":\"rename\",\"name\":\"{FILE_NAME}\"}}"
	}
	GetItemsPathEndPoint = config.GetItemsPathEndPoint
	if len(GetItemsPathEndPoint) == 0 {
		GetItemsPathEndPoint = "https://graph.microsoft.com/v1.0/me/drive/root:{PATH}:/children"
	}
	createAuths()
}

func RefreshTokenFunc() {
	if time.Now().After(ExpiredTime) {
		url := strings.Replace(RefreshAPIEndPoint, "{TENANT_ID}", TenantID, 1)

		payload := strings.NewReader("grant_type=refresh_token" +
			"&client_id=" + ClientID +
			"&client_secret=" + ClientSecret +
			"&scope=" + Scope +
			"&redirect_uri=" + RedirectUrl +
			"&refresh_token=" + RefreshToken)

		req, _ := http.NewRequest("POST", url, payload)

		req.Header.Add("content-type", "application/x-www-form-urlencoded")

		res, _ := http.DefaultClient.Do(req)

		var jsonResult models.RefreshTokenResponse
		err := json.NewDecoder(res.Body).Decode(&jsonResult)
		if err != nil {
			fmt.Println("Error: " + err.Error())
			return
		}
		defer res.Body.Close()

		if SavedToken != nil && jsonResult.AccessToken != SavedToken.AccessToken {
			saveToken(&jsonResult)
			fmt.Println("saved a new token")
			ExpiredTime = time.Now().Add(3000 * time.Second)
		} else {
			fmt.Println("nothing changed")
		}
	}
}

func saveToken(tokenJSON *models.RefreshTokenResponse) {
	// fmt.Println("AccessToken: " + tokenJSON.AccessToken)
	SavedToken.AccessToken = tokenJSON.AccessToken
	SavedToken.RefreshToken = tokenJSON.RefreshToken
	SavedToken.TokenType = tokenJSON.TokenType
	SavedToken.Expiry = time.Now().Add(3599 * time.Second)

	//assign the new refreshTokenStartTime
}
