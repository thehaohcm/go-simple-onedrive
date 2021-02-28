package service

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

func LoadConfigFromFile() (*Service, error) {
	service := new(Service)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("../")

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	loadConfigVariables(service)
	createAuths(service)
	return service, nil
}

func loadConfigVariables(service *Service) {
	service.ClientID = viper.GetString("CLIENT_ID")
	service.AccessToken = viper.GetString("ACCESS_TOKEN")
	service.AccessToken = viper.GetString("ACCESS_TOKEN")
	service.RedirectUrl = viper.GetString("REDIRECT_URL")
	service.TenantID = viper.GetString("TENANT_ID")
	service.RefreshToken = viper.GetString("REFRESH_TOKEN")
	service.Expiry = viper.GetInt64("EXPIRY")
	service.TokenType = viper.GetString("TOKEN_TYPE")
	service.UploadFolderPath = viper.GetString("UPLOAD_FOLDER_PATH")
	service.RefreshAPIEndPoint = viper.GetString("REFESH_API_ENDPOINT")
	service.UploadAPIEndPoint = viper.GetString("UPLOAD_API_ENDPOINT")
	service.ShareAPIEndPoint = viper.GetString("SHARE_API_ENDPOINT")
	service.CreateFolderAPIEndPoint = viper.GetString("CREATE_FOLDER_API_ENDPOINT")
	service.GetItemsPathEndPoint = viper.GetString("GET_ITEMS_PATH_ENDPOINT")
	service.GetItemAPIEndPoint = viper.GetString("GET_ITEM_API_ENDPOINT")
	service.DeleteItemAPIEndPoint = viper.GetString("DELETE_ITEM_API_ENDPOINT")
	service.MoveItemAPIEndPoint = viper.GetString("MOVE_ITEM_API_ENDPOINT")
	service.CopyItemAPIEndPoint = viper.GetString("COPY_ITEM_API_ENDPOINT")
	service.DownloadItemAPIEndPoint = viper.GetString("DOWNLOAD_ITEM_API_ENDPOINT")

	service.FragSize = viper.GetInt("FRAG_SIZE")
	service.ShareBodyJSON = viper.GetString("SHARE_BODY_JSON")
	service.UploadBodyJSON = viper.GetString("UPLOAD_BODY_JSON")
	service.CreateFolderBodyJSON = viper.GetString("CREATE_FOLDER_BODY_JSON")
	service.MoveItemBodyJSON = viper.GetString("MOVE_ITEM_BODY_JSON")
	service.CopyItemBodyJSON = viper.GetString("COPY_ITEM_BODY_JSON")

	service.ExpiredTime = time.Now().Add(3000 * time.Second)
}

func createAuths(service * Service) {

	service.SavedToken = &oauth2.Token{
		AccessToken:  service.AccessToken,
		RefreshToken: service.RefreshToken,
		Expiry:       time.Now().Add(time.Duration(viper.GetInt64("EXPIRY")) * time.Second),
		TokenType:    service.TokenType,
	}

	service.OauthConf = &oauth2.Config{
		ClientID:     service.ClientID,
		ClientSecret: service.ClientSecret,
		RedirectURL:  service.RedirectUrl,
		// Scopes:       []string{"Files.ReadWrite.All", "Sites.ReadWrite.All", "openid", "User.ReadBasic.All", "User.ReadWrite", "profile", "email"},
		Scopes:   strings.Fields(service.Scope),
		Endpoint: microsoft.AzureADEndpoint(service.TenantID),
	}

	//refresh Token
	RefreshTokenFunc(service)
}

func LoadConfigFromJson(config *models.Config) (*Service, error) {
	service := new(Service)
	service.ClientID = config.ClientID
	service.ClientSecret = config.ClientSecret
	service.AccessToken = config.AccessToken
	service.Scope = config.Scope
	if service.Scope == "" {
		service.Scope = "Files.ReadWrite.All openid User.ReadBasic.All User.ReadWrite profile email"
	}
	service.RedirectUrl = config.RedirectUrl
	service.TenantID = config.TenantID
	service.RefreshToken = config.RefreshToken
	service.Expiry = config.Expiry
	if service.Expiry == 0 {
		service.Expiry = 3599
	}
	service.TokenType = config.TokenType
	if len(service.TokenType) == 0 {
		service.TokenType = "Bearer"
	}
	service.UploadFolderPath = config.UploadFolderPath
	if len(service.UploadFolderPath) == 0 {
		service.UploadFolderPath = "/"
	}
	service.RefreshAPIEndPoint = config.RefreshAPIEndPoint
	if len(service.RefreshAPIEndPoint) == 0 {
		service.RefreshAPIEndPoint = "https://login.microsoftonline.com/{TENANT_ID}/oauth2/v2.0/token"
	}
	service.UploadAPIEndPoint = config.UploadAPIEndPoint
	if len(service.UploadAPIEndPoint) == 0 {
		service.UploadAPIEndPoint = "https://graph.microsoft.com/v1.0/me/drive/root:{UPLOAD_FOLDER_PATH}/{FILE_NAME}:/createUploadSession"
	}
	service.ShareAPIEndPoint = config.ShareAPIEndPoint
	if len(service.ShareAPIEndPoint) == 0 {
		service.ShareAPIEndPoint = "https://graph.microsoft.com/v1.0/me/drive/items/{UPLOADED_FILE_ID}/createLink"
	}
	service.CreateFolderAPIEndPoint = config.CreateFolderAPIEndPoint
	if len(service.CreateFolderAPIEndPoint) == 0 {
		service.CreateFolderAPIEndPoint = "https://graph.microsoft.com/v1.0/me/drive/items/{PARENT_FOLDER_ID}/children"
	}
	service.GetItemsPathEndPoint = config.GetItemsPathEndPoint
	if len(service.GetItemsPathEndPoint) == 0 {
		service.GetItemsPathEndPoint = "https://graph.microsoft.com/v1.0/me/drive/root:{PATH}:/children"
	}
	service.GetItemAPIEndPoint = config.GetItemAPIEndPoint
	if len(service.GetItemAPIEndPoint) == 0 {
		service.GetItemAPIEndPoint = "https://graph.microsoft.com/v1.0/me/drive/root:/{PATH}"
	}
	service.DeleteItemAPIEndPoint = config.DeleteItemAPIEndPoint
	if len(service.DeleteItemAPIEndPoint) == 0 {
		service.DeleteItemAPIEndPoint = "https://graph.microsoft.com/v1.0/me/drive/items/{ITEM_ID}"
	}
	service.MoveItemAPIEndPoint = config.MoveItemAPIEndPoint
	if len(service.MoveItemAPIEndPoint) == 0 {
		service.MoveItemAPIEndPoint = "https://graph.microsoft.com/v1.0/me/drive/items/{ITEM_ID}"
	}
	service.CopyItemAPIEndPoint = config.CopyItemAPIEndPoint
	if len(service.CopyItemAPIEndPoint) == 0 {
		service.CopyItemAPIEndPoint = "https://graph.microsoft.com/v1.0/me/drive/items/{ITEM_ID}/copy"
	}
	service.DownloadItemAPIEndPoint = config.DownloadItemAPIEndPoint
	if len(service.DownloadItemAPIEndPoint) == 0 {
		service.DownloadItemAPIEndPoint = "https://graph.microsoft.com/v1.0/me/drive/items/{ITEM_ID}/content"
	}
	service.FragSize = config.FragSize
	if service.FragSize == 0 {
		service.FragSize = 62259200
	}
	service.ShareBodyJSON = config.ShareBodyJSON
	if len(service.ShareBodyJSON) == 0 {
		service.ShareBodyJSON = "{ \"type\": \"view\", \"scope\": \"anonymous\" }"
	}
	service.UploadBodyJSON = config.UploadBodyJSON
	if len(service.UploadBodyJSON) == 0 {
		service.UploadBodyJSON = "{\"item\":{\"@microsoft.graph.conflictBehavior\":\"rename\",\"name\":\"{FILE_NAME}\"}}"
	}
	service.CreateFolderBodyJSON = config.CreateFolderBodyJSON
	if len(service.CreateFolderBodyJSON) == 0 {
		service.CreateFolderBodyJSON = "{\"name\": \"{FOLDER_NAME}\", \"folder\": {}, \"@microsoft.graph.conflictBehavior\": \"rename\" }"
	}
	service.MoveItemBodyJSON = config.MoveItemBodyJSON
	if len(service.MoveItemBodyJSON) == 0 {
		service.MoveItemBodyJSON = "{ \"parentReference\": { \"id\": \"{NEW_PARENT_FOLDER_ID}\" }, \"name\": \"{NEW_ITEM_NAME}\" }"
	}
	service.CopyItemBodyJSON = config.CopyItemBodyJSON
	if len(service.CopyItemBodyJSON) == 0 {
		service.CopyItemBodyJSON = "{ \"name\": \"{NEW_FILE_NAME}\" }"
	}

	createAuths(service)
	return service, nil
}

func RefreshTokenFunc(service * Service) {
	if time.Now().After(service.ExpiredTime) {
		url := strings.Replace(service.RefreshAPIEndPoint, "{TENANT_ID}", service.TenantID, 1)

		payload := strings.NewReader("grant_type=refresh_token" +
			"&client_id=" + service.ClientID +
			"&client_secret=" + service.ClientSecret +
			"&scope=" + service.Scope +
			"&redirect_uri=" + service.RedirectUrl +
			"&refresh_token=" + service.RefreshToken)

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

		saveToken(service, &jsonResult)
		saveToken(service, &jsonResult)
		fmt.Println("saved a new token")
		service.ExpiredTime = time.Now().Add(3000 * time.Second)
	} else {
		fmt.Println("nothing changed")
	}

}

func saveToken(service * Service, tokenJSON *models.RefreshTokenResponse) {
	// fmt.Println("AccessToken: " + tokenJSON.AccessToken)
	service.SavedToken.AccessToken = tokenJSON.AccessToken
	service.SavedToken.RefreshToken = tokenJSON.RefreshToken
	service.SavedToken.TokenType = tokenJSON.TokenType
	service.SavedToken.Expiry = time.Now().Add(3599 * time.Second)

	//assign the new refreshTokenStartTime
}
