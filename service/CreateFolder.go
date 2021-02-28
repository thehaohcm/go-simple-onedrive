package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/thehaohcm/go-simple-onedrive/models"
)

func (service *Service) CreateFolder(newFolderName string, parentFolder *models.ItemInfo) (*models.ItemInfo, error) {
	RefreshTokenFunc(service)

	url := strings.Replace(service.CreateFolderAPIEndPoint, "{PARENT_FOLDER_ID}", parentFolder.ID, 1)
	var payload = []byte(strings.Replace(service.CreateFolderBodyJSON, "{FOLDER_NAME}", newFolderName, 1))
	createFolderRequest, _ := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	createFolderRequest.Header.Add("Content-Type", "application/json")
	createFolderRequest.Header.Add("Authorization", service.TokenType+" "+service.SavedToken.AccessToken)

	client := &http.Client{}
	resp, err := client.Do(createFolderRequest)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == 201 {
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

func (service *Service) CreateFolderForParentName(newFolderName string, parentFolderName string) {
	parentFolderItem, err := service.GetItemByPath("/" + parentFolderName)
	if err != nil {
		panic(err)
	}

	if parentFolderItem != nil {
		newFolder, err := service.CreateFolder(newFolderName, parentFolderItem)
		if err != nil {
			panic(err)
		}
		fmt.Println("new Folder " + newFolder.Name + " is just created insider folder " + parentFolderItem.Name)
	}
}
