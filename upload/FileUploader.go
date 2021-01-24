package upload

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/thehaohcm/go-simple-onedrive/config"
	"github.com/thehaohcm/go-simple-onedrive/enums"
	"github.com/thehaohcm/go-simple-onedrive/models"
	"github.com/thehaohcm/go-simple-onedrive/token"
	"github.com/thehaohcm/go-simple-onedrive/utils"
)

var (
	blockSize              = 0
	fileBytes              []byte
	fileSize               = 0
	uploadFinishedResponse models.UploadFinishedResponse
)

func ShareLinkFunc(uploadFinishedResponse *models.UploadFinishedResponse) {
	token.RefreshToken()
	if uploadFinishedResponse != nil && uploadFinishedResponse.Id != "" {
		//share the item's link
		sharedLinkAPIEndpoint := "https://graph.microsoft.com/v1.0/me/drive/items/" + uploadFinishedResponse.Id + "/createLink"
		var httpHeaders [](*models.HttpHeader)
		httpHeaders = append(httpHeaders, models.InitHttpHeader("Content-Type", "application/json"))
		httpRequest := models.InitHttpRequest(enums.POST, sharedLinkAPIEndpoint, config.ShareBodyJSON, httpHeaders)

		var sharedLinkResponse models.SharedLinkResponse
		utils.HandleHttpRequestForUploading(httpRequest, &sharedLinkResponse)

		if sharedLinkResponse.Link.WebUrl != "" {
			fmt.Println("The file " + uploadFinishedResponse.Name + " has been shared via URL: " + sharedLinkResponse.Link.WebUrl + " for every one")
		}
	}
}

func UploadFile(localFilePath string) {

	token.RefreshToken()
	fileName := filepath.Base(localFilePath)

	fi, err := os.Open(localFilePath)
	fileBytes, err = ioutil.ReadFile(localFilePath)
	if err != nil {
		fmt.Println("Error 2: " + err.Error())
		return
	}

	//get file size
	fileData, err := fi.Stat()
	if err != nil {
		fmt.Println("Error 2: " + err.Error())
		return
	}
	fileSize = int(fileData.Size())

	// read file into bytes
	blockSize := (fileSize + config.FragSize - 1) / config.FragSize

	sessionUrL := "https://graph.microsoft.com/v1.0/me/drive/root:" + config.UploadFolderPath + fileName + ":/createUploadSession"
	//Create an Upload Session
	// fmt.Println("Creating an Uploading Session, with token: " + token.SavedToken.AccessToken)
	var payload = []byte(`{"item":{"@microsoft.graph.conflictBehavior":"rename","name":"` + fileName + `"}}`)
	uploadSessionRequest, _ := http.NewRequest("POST", sessionUrL, bytes.NewBuffer(payload))
	uploadSessionRequest.Header.Add("Content-Type", "application/json")
	uploadSessionRequest.Header.Add("Authorization", config.TokenType+" "+token.SavedToken.AccessToken)

	client := &http.Client{}
	resp, err := client.Do(uploadSessionRequest)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	var uploadJSONResult models.UploadSessionResponse
	json.Unmarshal(body, &uploadJSONResult)
	if uploadJSONResult.UploadUrl != "" {
		var uploadFinishedResponse models.UploadFinishedResponse
		fmt.Println("Uploading the file " + fileName + " into OneDrive...")
		for i := 0; i < blockSize; i++ {
			isSuccess := false
			numberOfAttempt := -1
			finishedPercent := float64(float64(i+1)/float64(blockSize)) * 100
			finishedPercentText := fmt.Sprintf("%.2f", finishedPercent)
			for !isSuccess {
				fmt.Println("Uploading fragment number: " + strconv.Itoa(i+1) + "/" + strconv.Itoa(blockSize) + "....(" + finishedPercentText + "%)")
				numberOfAttempt++
				if numberOfAttempt > 0 {
					fmt.Println("Uploading was failed, trying to upload again (Number of attempting: " + strconv.Itoa(numberOfAttempt) + "....")
				}
				//Compute the passTime to check if the token has to be refreshed or not
				// fmt.Println("subTime: " + strconv.Itoa(int(time.Now().Sub(token.RefreshTokenStartTime))))
				if time.Now().After(config.ExpiredTime) || numberOfAttempt > 0 {
					fmt.Println("The Token is expired, refreshing...")
					token.RefreshToken()
				}

				isSuccess = true
				var byteBlock []byte
				byteBlock = fileBytes[(i * config.FragSize):]
				param := "bytes " + strconv.Itoa(i*config.FragSize) + "-" + strconv.Itoa(fileSize-1) + "/" + strconv.Itoa(fileSize)
				if i < blockSize-1 {
					byteBlock = fileBytes[(i * config.FragSize):(i*config.FragSize + config.FragSize)]
					param = "bytes " + strconv.Itoa(i*config.FragSize) + "-" + strconv.Itoa(i*config.FragSize+config.FragSize-1) + "/" + strconv.Itoa(fileSize)
				}
				sizeByteBlock := len(byteBlock)
				uploadBlockFileRequest, _ := http.NewRequest("PUT", uploadJSONResult.UploadUrl, bytes.NewBuffer(byteBlock))
				uploadBlockFileRequest.Header.Add("Content-Length", strconv.Itoa(sizeByteBlock))
				uploadBlockFileRequest.Header.Add("Authorization", config.TokenType+" "+token.SavedToken.AccessToken)
				uploadBlockFileRequest.Header.Add("Content-Range", param)

				client := &http.Client{}
				resp, err := client.Do(uploadBlockFileRequest)
				if err != nil {
					isSuccess = false
					fmt.Println("The fragment File number " + strconv.Itoa(i+1) + " has error while uploading. Trying again...")
					continue
				}
				defer resp.Body.Close()
				body, _ := ioutil.ReadAll(resp.Body)
				if i < blockSize-1 {
					var uploadBlockResponse models.UploadBlockResponse
					json.Unmarshal(body, &uploadBlockResponse)
				} else {
					json.Unmarshal(body, &uploadFinishedResponse)
					fmt.Println("uploading is finished, file: " + uploadFinishedResponse.Name)
					fmt.Println("Download link: " + uploadFinishedResponse.DownloadUrl)
				}
			}
		}

		ShareLinkFunc(&uploadFinishedResponse)
	}
}
