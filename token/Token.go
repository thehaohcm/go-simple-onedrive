package token

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/thehaohcm/go-simple-onedrive/config"
	"github.com/thehaohcm/go-simple-onedrive/models"
)

func init() {
	RefreshToken()
}

func RefreshToken() {
	url := strings.Replace(config.RefreshAPIEndPoint, "{TENANT_ID}", config.TenantID, 1)

	payload := strings.NewReader("grant_type=refresh_token" +
		"&client_id=" + config.ClientID +
		"&client_secret=" + config.ClientSecret +
		"&scope=" + config.Scope +
		"&redirect_uri=" + config.RedirectUrl +
		"&refresh_token=" + config.RefreshToken)

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

	if config.SavedToken != nil && jsonResult.AccessToken != config.SavedToken.AccessToken {
		saveToken(&jsonResult)
		fmt.Println("saved a new token")
	} else {
		fmt.Println("nothing changed")
	}

}

func saveToken(tokenJSON *models.RefreshTokenResponse) {
	config.SavedToken.AccessToken = tokenJSON.AccessToken
	config.SavedToken.RefreshToken = tokenJSON.RefreshToken
	config.SavedToken.TokenType = tokenJSON.TokenType
	config.SavedToken.Expiry = time.Now().Add(3599 * time.Second)

	//assign the new refreshTokenStartTime
	config.ExpiredTime = time.Now().Add(3000 * time.Second)
}
