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

func (service *Service) CopyItem(itemInfo *models.ItemInfo, newItemName string) (bool, error) {
	RefreshTokenFunc(service)

	url := strings.Replace(service.CopyItemAPIEndPoint, "{ITEM_ID}", itemInfo.ID, 1)
	var payload = []byte(strings.Replace(service.CopyItemBodyJSON, "{NEW_FILE_NAME}", newItemName, 1))
	copyItemRequest, _ := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	copyItemRequest.Header.Add("Content-Type", "application/json")
	copyItemRequest.Header.Add("Authorization", service.TokenType+" "+service.SavedToken.AccessToken)

	client := &http.Client{}
	resp, err := client.Do(copyItemRequest)
	defer resp.Body.Close()
	if err != nil {
		panic(err)
	}

	//TODO: handle error message
	// body, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode == 202 {
		return true, nil
	}

	body, _ := ioutil.ReadAll(resp.Body)
	var errorResponse *models.ErrorResponse
	json.Unmarshal(body, &errorResponse)
	err = errors.New(errorResponse.Message)

	return false, err
}

func (service *Service) CopyItemByPath(itemPath string, newFileName string) (bool, error) {
	item, err := service.GetItemByPath(itemPath)
	if err != nil {
		panic(err)
	}

	if len(newFileName) == 0 {
		newFileName = item.Name
	}

	return service.CopyItem(item, newFileName)
}
