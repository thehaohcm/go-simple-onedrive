package upload

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/thehaohcm/go-simple-onedrive/config"
	"github.com/thehaohcm/go-simple-onedrive/models"
	"github.com/thehaohcm/go-simple-onedrive/utils"
)

func GetItemsByPath(path string) ([]*models.Value, error) {
	var url string
	if path == "" || path == "/" {
		//treat the root path as a special url
		url = strings.Replace(config.GetItemsPathEndPoint, ":{PATH}:", "", 1)
	} else {
		url = strings.Replace(config.GetItemsPathEndPoint, "{PATH}", path, 1)
	}
	// fmt.Println("url: " + url)
	getItemsRequest, _ := http.NewRequest("GET", url, nil)
	getItemsRequest.Header.Add("Content-Type", "application/json")
	getItemsRequest.Header.Add("Authorization", config.TokenType+" "+config.SavedToken.AccessToken)

	config.RefreshTokenFunc()

	client := &http.Client{}
	resp, err := client.Do(getItemsRequest)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("body: " + string(body))
	var ItemsIfoJSONResult models.ItemsInfoResponse
	json.Unmarshal(body, &ItemsIfoJSONResult)

	if ItemsIfoJSONResult.Value != nil {
		fmt.Println("The Path \"" + path + "\" has these files/folders: ")
		for _, item := range ItemsIfoJSONResult.Value {
			itemTxt := "- "
			if item.Folder != nil {
				itemTxt += item.Name + "[Folder] "
			} else {
				itemTxt += item.Name + " (" + utils.GetReadableFileCapacity(item.Size) + ")"
			}
			fmt.Println(itemTxt)
		}
	}

	return ItemsIfoJSONResult.Value, nil
}