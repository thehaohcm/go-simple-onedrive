package service

import (
	"bytes"
	"net/http"
	"strings"

	"github.com/thehaohcm/go-simple-onedrive/config"
	"github.com/thehaohcm/go-simple-onedrive/models"
)

func CopyItem(itemInfo *models.ItemInfo, newItemName string) (bool, error) {
	config.RefreshTokenFunc()

	url := strings.Replace(config.CopyItemAPIEndPoint, "{ITEM_ID}", itemInfo.ID, 1)
	var payload = []byte(strings.Replace(config.CopyItemBodyJSON, "{NEW_FILE_NAME}", newItemName, 1))
	copyItemRequest, _ := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	copyItemRequest.Header.Add("Content-Type", "application/json")
	copyItemRequest.Header.Add("Authorization", config.TokenType+" "+config.SavedToken.AccessToken)

	client := &http.Client{}
	resp, err := client.Do(copyItemRequest)
	defer resp.Body.Close()
	if err != nil {
		panic(err)
	}

	//TODO: handle error message
	// body, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode == 202 {
		return (true, nil)
	}

	body, _ := ioutil.ReadAll(resp.Body)
	var errorResponse *models.ErrorResponse
	json.Unmarshal(body, &errorResponse)
	err:=errors.New(errorResponse.Message)

	return false,err
}

func CopyItemByPath(itemPath string, newFileName string) (bool,error) {
	item, err := GetItemByPath(itemPath)
	if err != nil {
		panic(err)
	}

	if len(newFileName) == 0 {
		newFileName = item.Name
	}

	return CopyItem(item, newFileName)
}
