//
// Copyright (c) 2016 by Viacheslav Shynkarenko. All Rights Reserved.
//

package jo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// AssertOk checks expected properties of Ok response.
func AssertOk(t *testing.T, response *Response) {
	assert.NotNil(t, response)
	assert.Equal(t, 200, response.HTTPCode)
	assert.True(t, response.Successful)
}

// AssertFail checks expected properties of Fail response.
func AssertFail(t *testing.T, response *Response) {
	assert.NotNil(t, response)
	assert.NotEqual(t, 200, response.HTTPCode)
	assert.False(t, response.Successful)
}

// AssertForbidden checks expected properties of Forbidden response.
func AssertForbidden(t *testing.T, response *Response, messages ...string) {
	AssertHTTPError(403, "Forbidden", t, response, messages...)
}

// AssertBadRequest checks expected properties of BadRequest response.
func AssertBadRequest(t *testing.T, response *Response, messages ...string) {
	AssertHTTPError(400, "Bad Request", t, response, messages...)
}

// AssertUnauthorized checks expected properties of Unauthorized response.
func AssertUnauthorized(t *testing.T, response *Response, messages ...string) {
	AssertHTTPError(401, "Unauthorized", t, response, messages...)
}

// AssertResponseError checks expected properties of Error response.
func AssertResponseError(t *testing.T, response *Response, messages ...string) {
	AssertHTTPError(500, "Internal Error", t, response, messages...)
}

// AssertHTTPError checks expected properties of HTTP error response.
func AssertHTTPError(
	code int,
	defMessage string,
	t *testing.T,
	response *Response,
	messages ...string) {
	message := getHTTPErrorMessage(defMessage, messages...)
	AssertFail(t, response)
	assert.Equal(t, code, response.HTTPCode)
	assert.Equal(t, code, response.Error.Code)
	assert.Equal(t, message, response.Error.Message)
}

func getHTTPErrorMessage(def string, messages ...string) string {
	message := def
	if len(messages) > 0 {
		message = messages[0]
	}
	return message
}
