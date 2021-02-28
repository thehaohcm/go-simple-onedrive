package service

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/thehaohcm/go-simple-onedrive/models"
)

func (service *Service) GetChildItemsByPath(path string) ([]*models.ItemInfo, error) {
	RefreshTokenFunc(service)

	var url string
	if path == "" || path == "/" {
		//treat the root path as a special url
		url = strings.Replace(service.GetItemsPathEndPoint, ":{PATH}:", "", 1)
	} else {
		url = strings.Replace(service.GetItemsPathEndPoint, "{PATH}", path, 1)
	}
	getItemsRequest, _ := http.NewRequest("GET", url, nil)
	getItemsRequest.Header.Add("Content-Type", "application/json")
	getItemsRequest.Header.Add("Authorization", service.TokenType+" "+service.SavedToken.AccessToken)

	client := &http.Client{}
	resp, err := client.Do(getItemsRequest)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		var ItemsInfo *models.ItemsInfo
		json.Unmarshal(body, &ItemsInfo)

		return ItemsInfo.ItemInfos, nil
	}

	body, _ := ioutil.ReadAll(resp.Body)
	var errorResponse *models.ErrorResponse
	json.Unmarshal(body, &errorResponse)
	err = errors.New(errorResponse.Message)

	return nil, err
}

func (service *Service) GetItemByPath(path string) (*models.ItemInfo, error) {
	RefreshTokenFunc(service)

	url := strings.Replace(service.GetItemAPIEndPoint, "{PATH}", path, 1)
	getItemRequest, _ := http.NewRequest("GET", url, nil)
	getItemRequest.Header.Add("Content-Type", "application/json")
	getItemRequest.Header.Add("Authorization", service.TokenType+" "+service.SavedToken.AccessToken)

	client := &http.Client{}
	resp, err := client.Do(getItemRequest)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		var ItemInfo *models.ItemInfo
		json.Unmarshal(body, &ItemInfo)

		return ItemInfo, nil
	}

	body, _ := ioutil.ReadAll(resp.Body)
	var errorResponse *models.ErrorResponse
	json.Unmarshal(body, &errorResponse)
	err = errors.New(errorResponse.Message)
	return nil, err
}
