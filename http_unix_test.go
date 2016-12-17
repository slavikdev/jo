//
// Copyright (c) 2016 by Viacheslav Shynkarenko. All Rights Reserved.
//

package jo

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"runtime"
	"strings"
)

// HTTPUnixTest is a collection of functions to test API via HTTP requests to Unix socket.
type HTTPUnixTest struct {
	socket string
}

// NewHTTPTest creates new instance of HTTP testing framework.
func NewHTTPUnixTest(socket string) *HTTPUnixTest {
	if runtime.GOOS == "windows" {
		panic("HTTPUnixTest doesn't work on Windows.")
	}

	httpTest := &HTTPUnixTest{}
	httpTest.socket = socket
	return httpTest
}

// Get sends GET HTTP request to api endpoint and returns
// wrapped response for further testing.
func (ht *HTTPUnixTest) Get(url string) *Response {
	return ht.callAPI("GET", url, nil)
}

func (ht *HTTPUnixTest) callAPI(method string, url string, requestJSON interface{}) *Response {
	conn, _ := net.Dial("unix", ht.socket)
	request := ht.createRequest(method, url, requestJSON)
	fmt.Fprintf(conn, request)
	scanner := bufio.NewScanner(conn)
	var responseRaw string
	for scanner.Scan() {
		responseRaw += scanner.Text()
	}
	return ht.readRawResponse(responseRaw)
}

func (ht *HTTPUnixTest) createRequest(method string, url string, requestJSON interface{}) string {
	httpRequest := fmt.Sprintf(
		"%s %s HTTP/1.1\r\nHost: localhost\r\nContent-Type: application/json\r\n", method, url)
	if method != "get" && method != "delete" {
		jsonBytes := jsonToBytes(requestJSON)
		requestBody := bytesToString(jsonBytes)
		httpRequest = fmt.Sprintf("%s\r\n%s", httpRequest, requestBody)
	}
	return httpRequest
}

func (ht *HTTPUnixTest) readRawResponse(responseRaw string) *Response {
	fmt.Println(responseRaw)
	stringReader := strings.NewReader(responseRaw)
	bufReader := bufio.NewReader(stringReader)
	dummyRequest := http.Request{}
	response, err := http.ReadResponse(bufReader, &dummyRequest)
	if err != nil {
		log.Fatalf("Couldn't read response: %s", err.Error())
	}

	responseBuffer := bytes.NewBuffer(make([]byte, 0, response.ContentLength))
	responseBuffer.ReadFrom(response.Body)
	responseBytes := responseBuffer.Bytes()
	responseStr := bytesToString(responseBytes)
	apiResponse := &Response{}
	err = json.Unmarshal([]byte(responseStr), apiResponse)
	if nil != err {
		log.Fatalf("Couldn't parse response: %s", responseStr)
	}
	apiResponse.HTTPCode = response.StatusCode
	return apiResponse
}

func bytesToString(bytesArray []byte) string {
	buffer := bytes.NewBuffer(bytesArray)
	return buffer.String()
}
