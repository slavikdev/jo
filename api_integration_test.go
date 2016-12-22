//
// Copyright (c) 2016 by Viacheslav Shynkarenko. All Rights Reserved.
//

package jo

import (
	"fmt"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const inTestHost = "localhost:1234"
const inTestHostTLS = "localhost:12345"
const inTestSocket = "/tmp/jo_test_socket.sock"
const inTestCRTFile = "test_files/certificate.crt"
const inTestKeyFile = "test_files/privateKey.key"

// Runs API on TCP port.
// NOTE here we merely test our wrapper around gin's Run. Gin has its own tests.
func TestRun(t *testing.T) {
	api := newInTestAPI()
	host := inTestHost
	go func() {
		assert.NoError(t, api.Run(host))
	}()
	waitServer()

	inTestDefaultRoute(t, host)
}

// Runs API on unix socket.
func TestRunUnix(t *testing.T) {
	if runtime.GOOS == "windows" {
		fmt.Println("Skipping this test because unix sockets don't work on Windows.")
		return
	}

	api := newInTestAPI()
	socket := inTestSocket
	go func() {
		assert.NoError(t, api.RunUnix(socket))
	}()
	waitServer()
}

// Runs API on bad unix socket.
func TestRunUnixBadSocket(t *testing.T) {
	if runtime.GOOS == "windows" {
		fmt.Println("Skipping this test because unix sockets don't work on Windows.")
		return
	}

	api := newInTestAPI()
	socket := "###" + inTestSocket
	assert.Error(t, api.RunUnix(socket))
}

// Runs API on TCP port via TLS.
func TestRunTLS(t *testing.T) {
	api := newInTestAPI()
	host := inTestHostTLS
	go func() {
		assert.NoError(t,
			api.RunTLS(host, inTestCRTFile, inTestKeyFile))
	}()
	waitServer()

	inTestDefaultRouteTLS(t, host)
}

func waitServer() {
	time.Sleep(10 * time.Millisecond)
}

func inTestDefaultRoute(t *testing.T, host string) {
	http := NewHTTPIntegrationTest(host)
	response := http.Get("/")
	AssertOk(t, response)
}

func inTestDefaultRouteTLS(t *testing.T, host string) {
	https := NewHTTPIntegrationTestTLS(host, inTestCRTFile, inTestKeyFile)
	response := https.Get("/")
	AssertOk(t, response)
}

func newInTestAPI() *API {
	api, handlers, _ := newAPITest()
	api.Map("get", "/", handlers.emptyHandler)
	return api
}
