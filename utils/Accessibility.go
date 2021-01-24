package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"https://github.com/thehaohcm/go-simple-onedrive/enums"
	"https://github.com/thehaohcm/go-simple-onedrive/models"
	"https://github.com/thehaohcm/go-simple-onedrive/token"
	"strconv"
	"sync"
	"time"
)

var (
	uploadFolderPath       = "/Folder/File/"
	blockSize              = 0
	fileBytes              []byte
	fragSize               = 0
	fileSize               = 0
	uploadFinishedResponse models.UploadFinishedResponse
)

func ShareLinkFunc(uploadFinishedResponse *models.UploadFinishedResponse) {
	token.RefreshToken()
	if uploadFinishedResponse != nil && uploadFinishedResponse.Id != "" {
		//share the item's link
		sharedLinkAPIEndpoint := "https://graph.microsoft.com/v1.0/me/drive/items/" + uploadFinishedResponse.Id + "/createLink"
		payload := `{ "type": "view", "scope": "anonymous" }`
		var httpHeaders [](*models.HttpHeader)
		httpHeaders = append(httpHeaders, models.InitHttpHeader("Content-Type", "application/json"))
		httpRequest := models.InitHttpRequest(enums.POST, sharedLinkAPIEndpoint, payload, httpHeaders)

		var sharedLinkResponse models.SharedLinkResponse
		HandleHttpRequestForUploading(httpRequest, &sharedLinkResponse)

		if sharedLinkResponse.Link.WebUrl != "" {
			fmt.Println("The file " + uploadFinishedResponse.Name + " has been shared via URL: " + sharedLinkResponse.Link.WebUrl + " for every one")
		}
	}
}

func UploadFileWithSession(localFilePath string) {

	token.RefreshToken()
	fileName := filepath.Base(localFilePath)
	fragSize := 320 * 1024

	fi, err := os.Open(localFilePath)
	fileBytes, err := ioutil.ReadFile(localFilePath)
	if err != nil {
		fmt.Println("Error 2: " + err.Error())
		return
	}
	// defer fileBytes.Close()

	//get file size
	fileData, err := fi.Stat()
	if err != nil {
		fmt.Println("Error 2: " + err.Error())
		return
	}
	fileSize := int(fileData.Size())

	// read file into bytes
	blockSize := (fileSize + fragSize - 1) / fragSize

	sessionUrL := "https://graph.microsoft.com/v1.0/me/drive/root:" + uploadFolderPath + fileName + ":/createUploadSession"
	//Create an Upload Session
	// fmt.Println("Creating an Uploading Session, with token: " + token.SavedToken.AccessToken)
	var payload = []byte(`{"item":{"@microsoft.graph.conflictBehavior":"rename","name":"` + fileName + `"}}`)
	uploadSessionRequest, _ := http.NewRequest("POST", sessionUrL, bytes.NewBuffer(payload))
	uploadSessionRequest.Header.Add("Content-Type", "application/json")
	uploadSessionRequest.Header.Add("Authorization", "Bearer "+token.SavedToken.AccessToken)

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
				if time.Now().After(token.ExpiredTokenTime) || numberOfAttempt > 0 {
					fmt.Println("The Token is expired, refreshing...")
					token.RefreshToken()
				}

				isSuccess = true
				var byteBlock []byte
				endFlag := false
				if i == blockSize-1 {
					endFlag = true
					byteBlock = fileBytes[(i * fragSize):]
				} else {
					byteBlock = fileBytes[(i * fragSize):(i*fragSize + fragSize)]
				}
				sizeByteBlock := len(byteBlock)
				uploadBlockFileRequest, _ := http.NewRequest("PUT", uploadJSONResult.UploadUrl, bytes.NewBuffer(byteBlock))
				uploadBlockFileRequest.Header.Add("Content-Length", strconv.Itoa(sizeByteBlock))
				uploadBlockFileRequest.Header.Add("Authorization", "Bearer "+token.SavedToken.AccessToken)
				param := "bytes " + strconv.Itoa(i*fragSize) + "-" + strconv.Itoa(i*fragSize+fragSize-1) + "/" + strconv.Itoa(fileSize)
				if endFlag == true {
					param = "bytes " + strconv.Itoa(i*fragSize) + "-" + strconv.Itoa(fileSize-1) + "/" + strconv.Itoa(fileSize)
				}
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

func UploadFileWithSessionMultiThreading(localFilePath string) {
	fmt.Println("Refreshing Token before uploading file...")
	token.RefreshToken()
	fmt.Println("The Token has been valid...")
	fileName := filepath.Base(localFilePath)
	fragSize = 320 * 1024

	fi, err := os.Open(localFilePath)
	fileBytes, err = ioutil.ReadFile(localFilePath)
	if err != nil {
		fmt.Println("Error : " + err.Error())
		return
	}
	// defer fileBytes.Close()

	//get file size
	fileData, err := fi.Stat()
	if err != nil {
		fmt.Println("Error : " + err.Error())
		return
	}
	fileSize = int(fileData.Size())

	// read file into bytes
	blockSize = (fileSize + fragSize - 1) / fragSize

	sessionUrL := "https://graph.microsoft.com/v1.0/me/drive/root:" + uploadFolderPath + fileName + ":/createUploadSession"
	//Create an Upload Session
	// fmt.Println("Creating an Uploading Session, with token: " + token.SavedToken.AccessToken)
	var payload = []byte(`{"item":{"@microsoft.graph.conflictBehavior":"rename","name":"` + fileName + `"}}`)
	uploadSessionRequest, _ := http.NewRequest("POST", sessionUrL, bytes.NewBuffer(payload))
	uploadSessionRequest.Header.Add("Content-Type", "application/json")
	uploadSessionRequest.Header.Add("Authorization", "Bearer "+token.SavedToken.AccessToken)

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
		fmt.Println("Uploading the file " + fileName + " into OneDrive...")

		//init multi-thread
		var wg sync.WaitGroup
		wg.Add(blockSize - 1)

		//upload with multi-thread
		for i := 0; i < blockSize-1; i++ {
			go func(fragNumber int) {
				defer wg.Done()
				uploadFragment(fragNumber, uploadJSONResult)
				// uploadFragment(i, uploadJSONResult)
			}(i)
		}

		//upload the last fragment
		wg.Wait()
		time.Sleep(time.Duration(1) * time.Second)
		uploadFragment(blockSize-1, uploadJSONResult)

		ShareLinkFunc(&uploadFinishedResponse)
	}
}

func uploadFragment(fragNumber int, uploadJSONResult models.UploadSessionResponse) {
	isSuccess := false
	numberOfAttempt := -1
	finishedPercent := float64(float64(fragNumber+1)/float64(blockSize)) * 100
	finishedPercentText := fmt.Sprintf("%.2f", finishedPercent)
	for !isSuccess {
		fmt.Println("Uploading fragment number: " + strconv.Itoa(fragNumber+1) + "/" + strconv.Itoa(blockSize) + "....(" + finishedPercentText + "%)")
		numberOfAttempt++
		if numberOfAttempt > 0 {
			fmt.Println("Uploading was failed, trying to upload again (Number of attempting: " + strconv.Itoa(numberOfAttempt) + "....")
		}
		//Compute the passTime to check if the token has to be refreshed or not
		// fmt.Println("subTime: " + strconv.Itoa(int(time.Now().Sub(token.RefreshTokenStartTime))))
		if time.Now().After(token.ExpiredTokenTime) || numberOfAttempt > 0 {
			// if time.Now().After(token.ExpiredTokenTime) {
			fmt.Println("The Token is expired, refreshing...")
			token.RefreshToken()
		}

		isSuccess = true
		var byteBlock []byte
		endFlag := false
		byteBlock = fileBytes[(fragNumber * fragSize):(fragNumber*fragSize + fragSize)]
		if fragNumber == blockSize-1 {
			endFlag = true
			byteBlock = fileBytes[(fragNumber * fragSize):]
		}
		sizeByteBlock := len(byteBlock)
		uploadBlockFileRequest, _ := http.NewRequest("PUT", uploadJSONResult.UploadUrl, bytes.NewBuffer(byteBlock))
		uploadBlockFileRequest.Header.Add("Content-Length", strconv.Itoa(sizeByteBlock))
		uploadBlockFileRequest.Header.Add("Authorization", "Bearer "+token.SavedToken.AccessToken)
		param := "bytes " + strconv.Itoa(fragNumber*fragSize) + "-" + strconv.Itoa(fragNumber*fragSize+fragSize-1) + "/" + strconv.Itoa(fileSize)
		if endFlag == true {
			param = "bytes " + strconv.Itoa(fragNumber*fragSize) + "-" + strconv.Itoa(fileSize-1) + "/" + strconv.Itoa(fileSize)
		}
		uploadBlockFileRequest.Header.Add("Content-Range", param)

		client := &http.Client{}
		resp, err := client.Do(uploadBlockFileRequest)
		if err != nil {
			isSuccess = false
			fmt.Println("The fragment File number " + strconv.Itoa(fragNumber+1) + " has error while uploading. Trying again...")
			continue
		}
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		if fragNumber < blockSize-1 {
			var uploadBlockResponse models.UploadBlockResponse
			json.Unmarshal(body, &uploadBlockResponse)
		} else {
			json.Unmarshal(body, &uploadFinishedResponse)
			fmt.Println("uploading is finished, file: " + uploadFinishedResponse.Name)
			fmt.Println("Download link: " + uploadFinishedResponse.DownloadUrl)
		}
	}
}
