//
// Copyright (c) 2016 by Viacheslav Shynkarenko. All Rights Reserved.
//

package jo

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
)

// HTTPFunctionalTest is a collection of APIs to test framework via simulation of
// HTTP requests without any particular transport.
type HTTPFunctionalTest struct {
	api *API
}

// NewHTTPFunctionalTest creates new instance of HTTP testing framework.
func NewHTTPFunctionalTest(api *API) *HTTPFunctionalTest {
	httpTest := &HTTPFunctionalTest{}
	httpTest.api = api
	return httpTest
}

// Get sends GET HTTP request to api endpoint and returns
// wrapped response for further testing.
func (ht *HTTPFunctionalTest) Get(url string) *Response {
	return ht.callAPI("GET", url, nil)
}

// Delete sends GET HTTP request to api endpoint and returns
// wrapped response for further testing.
func (ht *HTTPFunctionalTest) Delete(url string) *Response {
	return ht.callAPI("DELETE", url, nil)
}

// Post sends POST HTTP request to api endpoint and returns
// wrapped response for further testing.
func (ht *HTTPFunctionalTest) Post(url string, requestJSON interface{}) *Response {
	return ht.callAPI("POST", url, requestJSON)
}

// Put sends PUT HTTP request to api endpoint and returns
// wrapped response for further testing.
func (ht *HTTPFunctionalTest) Put(url string, requestJSON interface{}) *Response {
	return ht.callAPI("PUT", url, requestJSON)
}

// Patch sends PATCH HTTP request to api endpoint and returns
// wrapped response for further testing.
func (ht *HTTPFunctionalTest) Patch(url string, requestJSON interface{}) *Response {
	return ht.callAPI("PATCH", url, requestJSON)
}

func (ht *HTTPFunctionalTest) callAPI(
	method string, url string, requestJSON interface{}) *Response {
	request := ht.createRequest(method, url, requestJSON)
	response := ht.getResponse(request)
	return response
}

func (ht *HTTPFunctionalTest) createRequest(
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

func (ht *HTTPFunctionalTest) getResponse(request *http.Request) *Response {
	recorder := httptest.NewRecorder()
	engine := ht.api.buildEngine()
	engine.ServeHTTP(recorder, request)
	response := ht.readResponse(recorder)
	return response
}

func jsonToBytes(requestJSON interface{}) []byte {
	jsonBytes, err := json.Marshal(requestJSON)
	if err != nil {
		jsonBytes = []byte{}
	}
	return jsonBytes
}

func (ht *HTTPFunctionalTest) readResponse(recorder *httptest.ResponseRecorder) *Response {
	response := &Response{}
	responseStr := recorder.Body.String()
	err := json.Unmarshal([]byte(responseStr), &response)
	if nil != err {
		log.Fatalf("Couldn't parse response: %s", responseStr)
	}
	response.HTTPCode = recorder.Code
	return response
}
