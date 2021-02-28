package service

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/thehaohcm/go-simple-onedrive/models"
)

func (service *Service) DeleteItem(itemInfo *models.ItemInfo) (bool, error) {
	RefreshTokenFunc(service)

	url := strings.Replace(service.DeleteItemAPIEndPoint, "{ITEM_ID}", itemInfo.ID, 1)
	deleteItemRequest, _ := http.NewRequest("DELETE", url, nil)
	deleteItemRequest.Header.Add("Content-Type", "application/json")
	deleteItemRequest.Header.Add("Authorization", service.TokenType+" "+service.SavedToken.AccessToken)

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
	err = errors.New(errorResponse.Message)

	return false, err
}

func (service *Service) DeleteItemByPath(itemPath string) (bool, error) {
	item, err := service.GetItemByPath(itemPath)
	if err != nil {
		return false, err
	}

	return service.DeleteItem(item)
}
