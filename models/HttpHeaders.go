package models

type HttpHeaders []*HttpHeader

type HttpHeader struct {
	Key   string
	Value string
}

func InitHttpHeader(key string, value string) *HttpHeader {
	var httpHeader HttpHeader
	httpHeader.Key = key
	httpHeader.Value = value
	return &httpHeader
}
