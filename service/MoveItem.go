package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/thehaohcm/go-simple-onedrive/models"
)

func (service *Service) MoveItem(itemInfo *models.ItemInfo, parentFolderInfo *models.ItemInfo, newItemName string) (bool, error) {
	RefreshTokenFunc(service)

	url := strings.Replace(service.MoveItemAPIEndPoint, "{ITEM_ID}", itemInfo.ID, 1)
	reqBody := strings.Replace(service.MoveItemBodyJSON, "{NEW_PARENT_FOLDER_ID}", parentFolderInfo.ID, 1)
	reqBody = strings.Replace(reqBody, "{NEW_ITEM_NAME}", newItemName, 1)
	var payload = []byte(reqBody)
	moveItemRequest, _ := http.NewRequest("PATCH", url, bytes.NewBuffer(payload))
	moveItemRequest.Header.Add("Content-Type", "application/json")
	moveItemRequest.Header.Add("Authorization", service.TokenType+" "+service.SavedToken.AccessToken)

	client := &http.Client{}
	resp, err := client.Do(moveItemRequest)
	defer resp.Body.Close()
	if err != nil {
		panic(err)
	}

	if resp.StatusCode == 200 {
		return true, nil
	}

	body, _ := ioutil.ReadAll(resp.Body)
	var errorResponse *models.ErrorResponse
	json.Unmarshal(body, &errorResponse)
	err = errors.New(errorResponse.Message)

	return false, err
}

func (service *Service) MoveItemByPathWithNewName(itemPath string, parentFolderPath string, newFileName string) (bool, error) {
	item, err := service.GetItemByPath(itemPath)
	if err != nil {
		panic(err)
	}

	parentFolderItem, err := service.GetItemByPath(parentFolderPath)
	if err != nil {
		panic(err)
	}

	if len(newFileName) == 0 {
		newFileName = item.Name
	}

	return service.MoveItem(item, parentFolderItem, newFileName)
}

func (service *Service) MoveItemByPath(itemPath string, parentFolderPath string) (bool, error) {
	return service.MoveItemByPathWithNewName(itemPath, parentFolderPath, "")
}
