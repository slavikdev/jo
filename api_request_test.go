//
// Copyright (c) 2016 by Viacheslav Shynkarenko. All Rights Reserved.
//

package jo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Creates things we need in every test.
func newAPITest() (*API, *testHandlers, *HTTPTest) {
	api := NewAPI()
	handlers := newTestHandlers()
	ht := NewHTTPTest(api)
	return api, handlers, ht
}

// Creates API with dummy Ok handler and sends request to it checking response.
func TestSimpleRequest(t *testing.T) {
	api, handlers, http := newAPITest()
	api.Map("get", "/simple-request", handlers.emptyHandler)

	response := http.Get("/simple-request")

	AssertOk(t, response)
}

// Creates API with predefined global context. Adds route handler which sends global context
// as response. Checks response to match initially set global context.
func TestGlobalContext(t *testing.T) {
	api, handlers, http := newAPITest()
	globalContext := make(map[string]string)
	globalContext["bebe"] = "meme"
	api.SetGlobalContext(globalContext)
	api.Map("get", "/global-context", handlers.passGlobalContext)

	response := http.Get("/global-context")

	AssertOk(t, response)
	responseData := response.Data.(map[string]interface{})
	assert.NotNil(t, responseData)
	assert.Equal(t, globalContext["bebe"], responseData["bebe"].(string))
}

// Creates API with dummy logger and adds route handler that uses logger.
func TestLoggingRequest(t *testing.T) {
	api, handlers, http := newAPITest()
	logger := testLogger{}
	api.SetLogger(&logger)
	api.Map("get", "/log", handlers.loggingHanler)

	response := http.Get("/log")

	AssertOk(t, response)
}

// Creates API with a route handled by 2 handlers. First one authorizes request, second
// can be only executed for authorized requests.
func TestAuthentication(t *testing.T) {
	api, handlers, http := newAPITest()
	api.Map("get", "/secured", handlers.authHandler, handlers.emptyHandler)

	// First we send request without token to get error.
	response := http.Get("/secured")
	AssertForbidden(t, response)

	// Now let's pass token to get through.
	response = http.Get("/secured?token=secret")
	AssertOk(t, response)
}

// Creates API with init request handler that validates request body.
func TestRequestValidation(t *testing.T) {
	api, handlers, http := newAPITest()
	api.SetInitRequestHandler(handlers.validateRequestHandler)
	endpoints := []string{"/path1", "/path2"}
	api.Map("post", endpoints[0], handlers.emptyHandler)
	api.Map("post", endpoints[1], handlers.emptyMessageHandler)

	for _, endpoint := range endpoints {
		// Passing no data should return bad request error.
		response := http.Post(endpoint, nil)
		AssertBadRequest(t, response)

		// Passing token but no session_id would still return bad request error.
		response = http.Post(endpoint+"?token=123", nil)
		AssertBadRequest(t, response)

		// Passing session_id but no token would still return bad request error.
		requestJSON := make(map[string]string)
		requestJSON["session_id"] = "qwerty"
		response = http.Post(endpoint, requestJSON)
		AssertBadRequest(t, response)

		// When token and session_id are present we should get through.
		response = http.Post(endpoint+"?token=123", requestJSON)
		AssertOk(t, response)
	}
}
