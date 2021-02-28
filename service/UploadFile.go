package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/thehaohcm/go-simple-onedrive/config"
	"github.com/thehaohcm/go-simple-onedrive/enums"
	"github.com/thehaohcm/go-simple-onedrive/models"
	"github.com/thehaohcm/go-simple-onedrive/utils"
)

var (
	blockSize              = 0
	fileBytes              []byte
	fileSize               int64 = 0
	uploadFinishedResponse models.UploadFinishedResponse
)

func ShareLinkFunc(uploadFinishedResponse *models.UploadFinishedResponse) *models.SharedLinkResponse {
	config.RefreshTokenFunc()
	if uploadFinishedResponse != nil && uploadFinishedResponse.Id != "" {
		//share the item's link
		sharedLinkAPIEndpoint := strings.Replace(config.ShareAPIEndPoint, "{UPLOADED_FILE_ID}", uploadFinishedResponse.Id, 1)
		var httpHeaders [](*models.HttpHeader)
		httpHeaders = append(httpHeaders, models.InitHttpHeader("Content-Type", "application/json"))
		httpRequest := models.InitHttpRequest(enums.POST, sharedLinkAPIEndpoint, config.ShareBodyJSON, httpHeaders)

		var sharedLinkResponse models.SharedLinkResponse
		utils.HandleHttpRequestForUploading(httpRequest, &sharedLinkResponse)
		if sharedLinkResponse.Link.WebUrl != "" {
			fmt.Println("The file " + uploadFinishedResponse.Name + " (size: " + utils.GetReadableFileCapacity(fileSize) + ") has been shared via URL: " + sharedLinkResponse.Link.WebUrl + " for every one")
			return &sharedLinkResponse
		}
	}
	return nil
}

func UploadFile(localFilePath string) (*models.UploadFinishedResponse, error) {
	config.RefreshTokenFunc()
	fileName := filepath.Base(localFilePath)

	fi, err := os.Open(localFilePath)
	fileBytes, err = ioutil.ReadFile(localFilePath)
	if err != nil {
		return nil, err
	}

	//get file size
	fileData, err := fi.Stat()
	if err != nil {
		return nil, err
	}
	fileSize = fileData.Size()

	// read file into bytes
	blockSize = 1
	if int(fileSize) > config.FragSize {
		blockSize = (int(fileSize) + config.FragSize - 1) / config.FragSize
	}

	sessionURL := strings.Replace(config.UploadAPIEndPoint, "{UPLOAD_FOLDER_PATH}", config.UploadFolderPath, 1)
	sessionURL = strings.Replace(sessionURL, "{FILE_NAME}", fileName, 1)
	//Create an Upload Session
	// fmt.Println("Creating an Uploading Session, with token: " + config.SavedToken.AccessToken)
	var payload = []byte(strings.Replace(config.UploadBodyJSON, "{FILE_NAME}", fileName, 1))
	uploadSessionRequest, _ := http.NewRequest("POST", sessionURL, bytes.NewBuffer(payload))
	uploadSessionRequest.Header.Add("Content-Type", "application/json")
	uploadSessionRequest.Header.Add("Authorization", config.TokenType+" "+config.SavedToken.AccessToken)

	client := &http.Client{}
	resp, err := client.Do(uploadSessionRequest)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	var uploadJSONResult models.UploadSessionResponse
	json.Unmarshal(body, &uploadJSONResult)
	if uploadJSONResult.UploadUrl != "" {
		var uploadFinishedResponse models.UploadFinishedResponse
		fmt.Println("Creating a session for Uploading the file " + fileName + " into OneDrive...")
		for i := 0; i < blockSize; i++ {
			isSuccess := false
			numberOfAttempt := -1
			finishedPercent := float64(float64(i+1)/float64(blockSize)) * 100
			finishedPercentText := fmt.Sprintf("%.2f", finishedPercent)
			for !isSuccess {
				numberOfAttempt++
				if numberOfAttempt > 0 {
					fmt.Println("Uploading was failed, trying to upload again (Number of attempting: " + strconv.Itoa(numberOfAttempt) + "....")
				}

				//refresh token if it is over time
				if time.Now().After(config.ExpiredTime) || numberOfAttempt > 0 {
					fmt.Println("The Token is expired, refreshing...")
					config.RefreshTokenFunc()
				}

				isSuccess = true
				var byteBlock []byte
				byteBlock = fileBytes[(i * config.FragSize):]
				param := "bytes " + strconv.Itoa(i*config.FragSize) + "-" + strconv.FormatInt(fileSize-1, 10) + "/" + strconv.FormatInt(fileSize, 10)
				if i < blockSize-1 {
					byteBlock = fileBytes[(i * config.FragSize):(i*config.FragSize + config.FragSize)]
					param = "bytes " + strconv.Itoa(i*config.FragSize) + "-" + strconv.Itoa(i*config.FragSize+config.FragSize-1) + "/" + strconv.FormatInt(fileSize, 10)
				}
				sizeByteBlock := len(byteBlock)
				uploadBlockFileRequest, _ := http.NewRequest("PUT", uploadJSONResult.UploadUrl, bytes.NewBuffer(byteBlock))
				uploadBlockFileRequest.Header.Add("Content-Length", strconv.Itoa(sizeByteBlock))
				uploadBlockFileRequest.Header.Add("Authorization", config.TokenType+" "+config.SavedToken.AccessToken)
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
				fmt.Println("Uploaded fragment number: " + strconv.Itoa(i+1) + "/" + strconv.Itoa(blockSize) + "....(" + finishedPercentText + "%)")
				if i < blockSize-1 {
					var uploadBlockResponse models.UploadBlockResponse
					json.Unmarshal(body, &uploadBlockResponse)
				} else {
					json.Unmarshal(body, &uploadFinishedResponse)
					fmt.Println("uploading is finished, file: " + uploadFinishedResponse.Name)
					fmt.Println("Download link: " + uploadFinishedResponse.DownloadUrl)

					//create a link for everyone can access
					// ShareLinkFunc(uploadFinishedResponse)
					return &uploadFinishedResponse, nil
				}
			}
		}
	}

	return nil, nil
}
