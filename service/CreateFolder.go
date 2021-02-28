package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/thehaohcm/go-simple-onedrive/config"
	"github.com/thehaohcm/go-simple-onedrive/models"
)

func CreateFolder(newFolderName string, parentFolder *models.ItemInfo) (*models.ItemInfo, error) {
	config.RefreshTokenFunc()

	url := strings.Replace(config.CreateFolderAPIEndPoint, "{PARENT_FOLDER_ID}", parentFolder.ID, 1)
	var payload = []byte(strings.Replace(config.CreateFolderBodyJSON, "{FOLDER_NAME}", newFolderName, 1))
	createFolderRequest, _ := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	createFolderRequest.Header.Add("Content-Type", "application/json")
	createFolderRequest.Header.Add("Authorization", config.TokenType+" "+config.SavedToken.AccessToken)

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
	err := errors.New(errorResponse.Message)

	return nil, err
}

func CreateFolderForParentName(newFolderName string, parentFolderName string) {
	parentFolderItem, err := GetItemByPath("/" + parentFolderName)
	if err != nil {
		panic(err)
	}

	if parentFolderItem != nil {
		newFolder, err := CreateFolder(newFolderName, parentFolderItem)
		if err != nil {
			panic(err)
		}
		fmt.Println("new Folder " + newFolder.Name + " is just created insider folder " + parentFolderItem.Name)
	}
}
