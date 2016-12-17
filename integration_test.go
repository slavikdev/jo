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

// Runs API on TCP port.
// NOTE here we merely test our wrapper around gin's Run. Gin has its own tests.
func TestRun(t *testing.T) {
	api := newAPITestIntegration()
	go func() {
		assert.NoError(t, api.Run(":9988"))
	}()
	waitServer()
}

// Runs API on unix socket.
func TestRunUnix(t *testing.T) {
	if runtime.GOOS == "windows" {
		fmt.Println("Skipping this test because unix sockets don't work on Windows.")
		return
	}

	api := newAPITestIntegration()
	socket := "/tmp/jo_test_socket.sock"
	go func() {
		assert.NoError(t, api.RunUnix(socket))
	}()
	waitServer()

	// unixTest := NewHTTPUnixTest(socket)
	// response := unixTest.Get("/")
	// AssertOk(t, response)
}

// Runs API on TCP port via TLS.
func TestRunTLS(t *testing.T) {
	api := newAPITestIntegration()
	go func() {
		assert.NoError(t,
			api.RunTLS(":12345", "test_files/certificate.crt", "test_files/privateKey.key"))
	}()
	waitServer()
}

func waitServer() {
	time.Sleep(10 * time.Millisecond)
}

func newAPITestIntegration() *API {
	api, handlers, _ := newAPITest()
	api.Map("get", "/", handlers.emptyHandler)
	return api
}
