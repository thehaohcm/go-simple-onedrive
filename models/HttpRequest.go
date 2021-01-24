package models

import (
	"github.com/thehaohcm/go-simple-onedrive/enums"
)

type HttpRequest struct {
	HttpMethod enums.HttpRequestMethod
	Url        string
	Body       string
	Headers    []*HttpHeader
}

func InitHttpRequest(httpMethod enums.HttpRequestMethod, url string, body string, headers []*HttpHeader) *HttpRequest {
	var httpRequest HttpRequest
	httpRequest.HttpMethod = httpMethod
	httpRequest.Url = url
	httpRequest.Body = body
	httpRequest.Headers = headers
	return &httpRequest
}
