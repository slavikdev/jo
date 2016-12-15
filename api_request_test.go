//
// Copyright (c) 2016 by Viacheslav Shynkarenko. All Rights Reserved.
//

package jo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Creates API with dummy Ok handler and sends request to it checking response.
func TestSimpleRequest(t *testing.T) {
	api := NewAPI()
	api.Map("get", "/simple-request", emptyHandler)

	ht := NewHTTPTest(api)
	response := ht.Get("/simple-request")

	AssertOk(t, response)
}

// Creates API with predefined global context. Adds route handler which sends global context
// as response. Checks response to match initially set global context.
func TestGlobalContext(t *testing.T) {
	api := NewAPI()
	globalContext := make(map[string]string)
	globalContext["bebe"] = "meme"
	api.SetGlobalContext(globalContext)
	api.Map("get", "/global-context", passGlobalContext)

	ht := NewHTTPTest(api)
	response := ht.Get("/global-context")

	AssertOk(t, response)
	responseData := response.Data.(map[string]interface{})
	assert.NotNil(t, responseData)
	assert.Equal(t, globalContext["bebe"], responseData["bebe"].(string))
}

// Creates API with dummy logger and adds route handler that uses logger.
func TestLoggingRequest(t *testing.T) {
	api := NewAPI()
	logger := testLogger{}
	api.SetLogger(&logger)
	api.Map("get", "/log", loggingHanler)

	ht := NewHTTPTest(api)
	response := ht.Get("/log")

	AssertOk(t, response)
}

// Creates API with a route handled by 2 handlers. First one authorizes request, second
// can be only executed for authorized requests.
func TestAuthentication(t *testing.T) {
	api := NewAPI()
	api.Map("get", "/secured", authHandler, emptyHandler)

	ht := NewHTTPTest(api)

	// First we send request without token to get error.
	response := ht.Get("/secured")
	AssertFail(t, response)
	assert.Equal(t, 403, response.HTTPCode)
	assert.Equal(t, 403, response.Error.Code)
	assert.Equal(t, "Forbidden", response.Error.Message)

	// Now let's pass token to get through.
	response = ht.Get("/secured?token=secret")
	AssertOk(t, response)
}
