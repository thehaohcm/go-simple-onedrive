package upload

import (
	"bytes"
	"net/http"
	"strings"

	"github.com/thehaohcm/go-simple-onedrive/config"
	"github.com/thehaohcm/go-simple-onedrive/models"
)

func MoveItem(itemInfo *models.ItemInfo, parentFolderInfo *models.ItemInfo, newItemName string) bool {
	config.RefreshTokenFunc()

	url := strings.Replace(config.MoveItemAPIEndPoint, "{ITEM_ID}", itemInfo.ID, 1)
	reqBody := strings.Replace(config.MoveItemBodyJSON, "{NEW_PARENT_FOLDER_ID}", parentFolderInfo.ID, 1)
	reqBody = strings.Replace(reqBody, "{NEW_ITEM_NAME}", newItemName, 1)
	var payload = []byte(reqBody)
	moveItemRequest, _ := http.NewRequest("PATCH", url, bytes.NewBuffer(payload))
	moveItemRequest.Header.Add("Content-Type", "application/json")
	moveItemRequest.Header.Add("Authorization", config.TokenType+" "+config.SavedToken.AccessToken)

	client := &http.Client{}
	resp, err := client.Do(moveItemRequest)
	defer resp.Body.Close()
	if err != nil {
		panic(err)
	}

	if resp.StatusCode == 200 {
		return true
	}

	return false
}

func MoveItemByPathWithNewName(itemPath string, parentFolderPath string, newFileName string) bool {
	item, err := GetItemByPath(itemPath)
	if err != nil {
		panic(err)
	}

	parentFolderItem, err := GetItemByPath(parentFolderPath)
	if err != nil {
		panic(err)
	}

	if len(newFileName) == 0 {
		newFileName = item.Name
	}

	return MoveItem(item, parentFolderItem, newFileName)
}

func MoveItemByPath(itemPath string, parentFolderPath string) bool {
	return MoveItemByPathWithNewName(itemPath, parentFolderPath, "")
}
