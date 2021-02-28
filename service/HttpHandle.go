package service

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/thehaohcm/go-simple-onedrive/models"
)

func HandleHttpRequestForUploading(service *Service, httpRequest *models.HttpRequest, response interface{}) {
	bodyBytes := bytes.NewBuffer([]byte(httpRequest.Body))
	request, _ := http.NewRequest(string(httpRequest.HttpMethod), httpRequest.Url, bodyBytes)
	for _, header := range httpRequest.Headers {
		request.Header.Add(header.Key, header.Value)
	}
	request.Header.Add("Authorization", "Bearer "+service.SavedToken.AccessToken)
	httpClient := &http.Client{}
	resp, err := httpClient.Do(request)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	json.Unmarshal(body, &response)
}
