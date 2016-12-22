//
// Copyright (c) 2016 by Viacheslav Shynkarenko. All Rights Reserved.
//

package jo

import (
	"strings"
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
	AssertUnauthorized(t, response)

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
		testSecuredEndpoint(t, http, endpoint)
	}
}

// Creates API with a route protected by two validation handlers.
// This test ensures chained handlers work properly.
func TestChainedRequestValidation(t *testing.T) {
	api, handlers, http := newAPITest()
	endpoint := "/protected/area"
	api.Map(
		"post",
		endpoint,
		handlers.validateToken,
		handlers.validateSessionID,
		handlers.emptyMessageHandler)
	testSecuredEndpoint(t, http, endpoint)
}

func testSecuredEndpoint(t *testing.T, http *HTTPTest, endpoint string) {
	// Passing no data should return Forbidden error.
	response := http.Post(endpoint, nil)
	AssertForbidden(t, response, "Token required")

	// Passing invalid token results in bad request.
	response = http.Post(endpoint+"?token=123", nil)
	AssertBadRequest(t, response, "Invalid token")

	// Passing valid token results in bad request because there's no JSON body.
	response = http.Post(endpoint+"?token=S123", nil)
	AssertBadRequest(t, response, "No body")

	// Passing valid token and some JSON body results in forbidden.
	response = http.Post(endpoint+"?token=S123", make(map[string]string))
	AssertForbidden(t, response, "Session required")

	// Passing invalid session id returns bad request.
	requestJSON := make(map[string]string)
	requestJSON["session_id"] = "qwerty"
	response = http.Post(endpoint+"?token=S123", requestJSON)
	AssertBadRequest(t, response, "Invalid session")

	// Both token and session id are valid, should be OK.
	requestJSON["session_id"] = "ID12345"
	response = http.Post(endpoint+"?token=S123", requestJSON)
	AssertOk(t, response)
}

// Creates API with end request handler that patches response data.
func TestResponsePatching(t *testing.T) {
	api, handlers, http := newAPITest()
	api.SetEndRequestHandler(handlers.patchResponse)
	endpoints := []string{"/path1", "/path2"}
	api.Map("get", endpoints[0], handlers.emptyHandler)
	api.Map("get", endpoints[1], handlers.emptyMessageHandler)

	for _, endpoint := range endpoints {
		// Passing no data should return bad request error.
		response := http.Get(endpoint)
		AssertOk(t, response)

		responseData := response.Data.(map[string]interface{})
		assert.NotNil(t, responseData)
		assert.NotNil(t, responseData["date"])
		assert.NotNil(t, responseData["response"])
		assert.NotEmpty(t, responseData["date"].(string))
	}
}

// Creates API with mappings for every supported HTTP method.
func TestEachHTTPMethod(t *testing.T) {
	api, handlers, http := newAPITest()
	api.SetEndRequestHandler(handlers.patchResponse)
	endpoint := "/a/b/c/d/e"
	httpMethods := "get,post,put,patch,delete"
	api.Map(httpMethods, endpoint, handlers.emptyHandler)

	httpMethodsSplit := strings.Split(httpMethods, ",")
	requestBody := make(map[string]string)
	requestBody["data"] = "dnepr"
	for _, httpMethod := range httpMethodsSplit {
		var response *Response
		switch httpMethod {
		case "get":
			response = http.Get(endpoint)
		case "post":
			response = http.Post(endpoint, requestBody)
		case "put":
			response = http.Put(endpoint, requestBody)
		case "patch":
			response = http.Patch(endpoint, requestBody)
		case "delete":
			response = http.Delete(endpoint)
		}

		AssertOk(t, response)
	}
}
