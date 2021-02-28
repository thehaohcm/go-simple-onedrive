package servicee

import (
	"errors"
	"net/http"
	"strings"

	"github.com/thehaohcm/go-simple-onedrive/config"
	"github.com/thehaohcm/go-simple-onedrive/models"
)

func GetDownloadLinkItem(itemInfo *models.ItemInfo) (string, error) {
	config.RefreshTokenFunc()

	url := strings.Replace(config.DownloadItemAPIEndPoint, "{ITEM_ID}", itemInfo.ID, 1)
	downloadItemRequest, _ := http.NewRequest("GET", url, nil)
	downloadItemRequest.Header.Add("Content-Type", "application/json")
	downloadItemRequest.Header.Add("Authorization", config.TokenType+" "+config.SavedToken.AccessToken)

	client := &http.Client{}
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return errors.New("Redirect")
	}
	resp, err := client.Do(downloadItemRequest)
	defer resp.Body.Close()
	var downloadLink string
	if err != nil {
		if resp.StatusCode == http.StatusFound { //status code 302
			downloadLinkUrl, err := resp.Location()
			if err != nil {
				return "", err
			}
			downloadLink = downloadLinkUrl.String()
		} else {
			return "", err
		}
	}

	return downloadLink, nil
}

func GetDownloadLinkItemByPath(itemPath string) (string, error) {
	item, err := GetItemByPath(itemPath)
	if err != nil {
		return "", err
	}

	if item != nil {
		return GetDownloadLinkItem(item)
	}
	return "", nil
}
