package service

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/thehaohcm/go-simple-onedrive/config"
	"github.com/thehaohcm/go-simple-onedrive/models"
)

func DeleteItem(itemInfo *models.ItemInfo) (bool, error) {
	config.RefreshTokenFunc()

	url := strings.Replace(config.DeleteItemAPIEndPoint, "{ITEM_ID}", itemInfo.ID, 1)
	deleteItemRequest, _ := http.NewRequest("DELETE", url, nil)
	deleteItemRequest.Header.Add("Content-Type", "application/json")
	deleteItemRequest.Header.Add("Authorization", config.TokenType+" "+config.SavedToken.AccessToken)

	client := &http.Client{}
	resp, err := client.Do(deleteItemRequest)
	defer resp.Body.Close()
	if err != nil {
		panic(err)
	}

	if resp.StatusCode == 204 {
		return true, nil
	}

	body, _ := ioutil.ReadAll(resp.Body)
	var errorResponse *models.ErrorResponse
	json.Unmarshal(body, &errorResponse)
	err := errors.New(errorResponse.Message)

	return false, err
}

func DeleteItemByPath(itemPath string) bool {
	item, err := GetItemByPath(itemPath)
	if err != nil {
		panic(err)
	}

	if item != nil {
		if DeleteItem(item) {
			return true
		}
	}
	return false
}
