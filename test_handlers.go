//
// Copyright (c) 2016 by Viacheslav Shynkarenko. All Rights Reserved.
//

package jo

import (
	"strings"
	"time"
)

// testHandlers is a collection of route handlers used in tests.
type testHandlers struct {
}

func newTestHandlers() *testHandlers {
	handlers := &testHandlers{}
	return handlers
}

// Returns Ok response. Used when handler logic is irrelevant.
func (handlers *testHandlers) emptyHandler(request *Request) *Response {
	return Ok(true)
}

// Returns Ok response with a message. Used when handler logic is irrelevant.
func (handlers *testHandlers) emptyMessageHandler(request *Request) *Response {
	return Ok("hello")
}

// Takes global context from request context and returns it as Data in successful response.
func (handlers *testHandlers) passGlobalContext(request *Request) *Response {
	return Ok(request.GlobalContext)
}

// Expects external logger to be set on API level.
// Calls every logging function and returns ok.
func (handlers *testHandlers) loggingHanler(request *Request) *Response {
	request.Logger.Debug("Hello %s", "World")
	request.Logger.Info("Hello %s", "World")
	request.Logger.Warn("Hello %s", "World")
	request.Logger.Error("Hello %s", "World")
	return Ok(true)
}

// Checks whether request has special query string argument "token" with value secret.
// If it does--request goes to the next handler. Otherwise 403 Forbidden is returned.
func (handlers *testHandlers) authHandler(request *Request) *Response {
	if request.GetQuery("token") == "secret" {
		return Next()
	}
	return Unauthorized()
}

// Validates request to be of specific structure. Used as init request handler.
func (handlers *testHandlers) validateRequestHandler(request *Request) *Response {
	// Every request must have token
	response := handlers.validateToken(request)
	if response.EndRequest {
		return response
	}
	return handlers.validateSessionID(request)
}

// Validates presence of a token.
func (handlers *testHandlers) validateToken(request *Request) *Response {
	token := request.GetQuery("token")
	if len(token) == 0 {
		return ForbiddenMessage("Token required")
	}

	// Token must begin with letter S.
	if strings.Index(token, "S") == 0 {
		return Next()
	}

	return BadRequestMessage("Invalid token")
}

// Validates presence of session_id
func (handlers *testHandlers) validateSessionID(request *Request) *Response {
	json := make(map[string]interface{})
	request.GetJSON(&json)
	if json == nil {
		return BadRequestMessage("No body")
	}

	if json["session_id"] == nil || len(json["session_id"].(string)) == 0 {
		return ForbiddenMessage("Session required")
	}

	sessionID := json["session_id"].(string)
	// Session ID must begin with ID.
	if strings.Index(sessionID, "ID") == 0 {
		return Next(sessionID)
	}

	return BadRequestMessage("Invalid session")
}

// Patches previous response data.
func (handlers *testHandlers) patchResponse(request *Request) *Response {
	response := request.PrevHandlerResponse
	wrapper := make(map[string]interface{})
	wrapper["date"] = time.Now().String()
	wrapper["response"] = response.Data
	response.Data = wrapper
	return response
}
