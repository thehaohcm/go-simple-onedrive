package upload

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/thehaohcm/go-simple-onedrive/config"
	"github.com/thehaohcm/go-simple-onedrive/models"
)

func GetChildItemsByPath(path string) ([]*models.ItemInfo, error) {
	config.RefreshTokenFunc()

	var url string
	if path == "" || path == "/" {
		//treat the root path as a special url
		url = strings.Replace(config.GetItemsPathEndPoint, ":{PATH}:", "", 1)
	} else {
		url = strings.Replace(config.GetItemsPathEndPoint, "{PATH}", path, 1)
	}
	getItemsRequest, _ := http.NewRequest("GET", url, nil)
	getItemsRequest.Header.Add("Content-Type", "application/json")
	getItemsRequest.Header.Add("Authorization", config.TokenType+" "+config.SavedToken.AccessToken)

	client := &http.Client{}
	resp, err := client.Do(getItemsRequest)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	var ItemsInfo *models.ItemsInfo
	json.Unmarshal(body, &ItemsInfo)

	return ItemsInfo.ItemInfos, nil
}

func GetItemByPath(path string) (*models.ItemInfo, error) {
	config.RefreshTokenFunc()

	url := strings.Replace(config.GetItemAPIEndPoint, "{PATH}", path, 1)
	getItemRequest, _ := http.NewRequest("GET", url, nil)
	getItemRequest.Header.Add("Content-Type", "application/json")
	getItemRequest.Header.Add("Authorization", config.TokenType+" "+config.SavedToken.AccessToken)

	client := &http.Client{}
	resp, err := client.Do(getItemRequest)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	var ItemInfo *models.ItemInfo
	json.Unmarshal(body, &ItemInfo)

	return ItemInfo, nil
}
