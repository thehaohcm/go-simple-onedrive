package service

import (
	"errors"
	"net/http"
	"strings"

	"github.com/thehaohcm/go-simple-onedrive/models"
)

func (service *Service) GetDownloadLinkItem(itemInfo *models.ItemInfo) (string, error) {
	RefreshTokenFunc(service)

	url := strings.Replace(service.DownloadItemAPIEndPoint, "{ITEM_ID}", itemInfo.ID, 1)
	downloadItemRequest, _ := http.NewRequest("GET", url, nil)
	downloadItemRequest.Header.Add("Content-Type", "application/json")
	downloadItemRequest.Header.Add("Authorization", service.TokenType+" "+service.SavedToken.AccessToken)

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

func (service *Service) GetDownloadLinkItemByPath(itemPath string) (string, error) {
	item, err := service.GetItemByPath(itemPath)
	if err != nil {
		return "", err
	}

	if item != nil {
		return service.GetDownloadLinkItem(item)
	}
	return "", nil
}
