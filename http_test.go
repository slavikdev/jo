package jo

import (
	"bytes"
	"log"
	"net/http"
)

// HTTPTest is a definition of HTTP testing API.
type HTTPTest interface {
	Get(url string) *Response
	Delete(url string) *Response
	Post(url string, requestJSON interface{}) *Response
	Put(url string, requestJSON interface{}) *Response
	Patch(url string, requestJSON interface{}) *Response
}

func createHTTPTestRequest(
	method string, url string, requestJSON interface{}) *http.Request {
	jsonBytes := jsonToBytes(requestJSON)
	byteReader := bytes.NewReader(jsonBytes)
	request, err := http.NewRequest(method, url, byteReader)
	if nil != err {
		log.Fatalf("Couldn't create request: %s %s", method, url)
	}
	request.Header.Add("Content-Type", "application/json")
	return request
}
