//
// Copyright (c) 2016 by Viacheslav Shynkarenko. All Rights Reserved.
//

package jo

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
)

// HTTPIntegrationTest is a collection of APIs to test framework via
// real HTTP requests over TCP connection.
type HTTPIntegrationTest struct {
	proto     string
	host      string
	socket    string
	transport *http.Transport
}

// NewHTTPIntegrationTest creates new instance of HTTP testing framework.
func NewHTTPIntegrationTest(host string) *HTTPIntegrationTest {
	httpTest := &HTTPIntegrationTest{}
	httpTest.host = host
	httpTest.proto = "http"
	return httpTest
}

// NewHTTPIntegrationTestTLS creates new instance of HTTP testing framework with TLS
// configuration, allowing to test over HTTPS.
func NewHTTPIntegrationTestTLS(host, crt, key string) *HTTPIntegrationTest {
	cert, err := tls.LoadX509KeyPair(crt, key)
	if err != nil {
		panic(err)
	}

	tlsConfig := &tls.Config{
		Certificates:       []tls.Certificate{cert},
		InsecureSkipVerify: true,
	}
	tlsConfig.BuildNameToCertificate()
	httpTest := NewHTTPIntegrationTest(host)
	httpTest.proto = "https"
	httpTest.transport = &http.Transport{
		TLSClientConfig:    tlsConfig,
		DisableCompression: true,
	}
	return httpTest
}

// NewHTTPIntegrationTestUnix creates new instance of HTTP testing framework to run
// requests via unix socket.
func NewHTTPIntegrationTestUnix(host, socket string) *HTTPIntegrationTest {
	httpTest := NewHTTPIntegrationTest(host)
	httpTest.proto = "http"
	httpTest.socket = socket
	httpTest.transport = &http.Transport{
		Dial: httpTest.dialUnixSocket,
	}
	return httpTest
}

// Get sends GET HTTP request to api endpoint and returns
// wrapped response for further testing.
func (ht *HTTPIntegrationTest) Get(url string) *Response {
	return ht.callAPI("GET", url, nil)
}

// Delete sends GET HTTP request to api endpoint and returns
// wrapped response for further testing.
func (ht *HTTPIntegrationTest) Delete(url string) *Response {
	return ht.callAPI("DELETE", url, nil)
}

// Post sends POST HTTP request to api endpoint and returns
// wrapped response for further testing.
func (ht *HTTPIntegrationTest) Post(url string, requestJSON interface{}) *Response {
	return ht.callAPI("POST", url, requestJSON)
}

// Put sends PUT HTTP request to api endpoint and returns
// wrapped response for further testing.
func (ht *HTTPIntegrationTest) Put(url string, requestJSON interface{}) *Response {
	return ht.callAPI("PUT", url, requestJSON)
}

// Patch sends PATCH HTTP request to api endpoint and returns
// wrapped response for further testing.
func (ht *HTTPIntegrationTest) Patch(url string, requestJSON interface{}) *Response {
	return ht.callAPI("PATCH", url, requestJSON)
}

func (ht *HTTPIntegrationTest) callAPI(
	method string, url string, requestJSON interface{}) *Response {
	fullURL := fmt.Sprintf("%s://%s%s", ht.proto, ht.host, url)
	request := createHTTPTestRequest(method, fullURL, requestJSON)
	client := &http.Client{}
	if ht.transport != nil {
		client.Transport = ht.transport
	}
	httpResponse, err := client.Do(request)
	defer httpResponse.Body.Close()
	if err != nil {
		panic(err)
	}
	response := ht.readResponse(httpResponse)
	return response
}

func (ht *HTTPIntegrationTest) readResponse(httpResponse *http.Response) *Response {
	body, err := ioutil.ReadAll(httpResponse.Body)
	if err != nil {
		panic(err)
	}
	response := &Response{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Fatalf("Couldn't parse response: %s", body)
	}
	response.HTTPCode = httpResponse.StatusCode
	return response
}

func (ht *HTTPIntegrationTest) dialUnixSocket(
	proto, addr string) (conn net.Conn, err error) {
	return net.Dial("unix", ht.socket)
}
