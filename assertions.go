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
func AssertForbidden(t *testing.T, response *Response) {
	AssertFail(t, response)
	assert.Equal(t, 403, response.HTTPCode)
	assert.Equal(t, 403, response.Error.Code)
	assert.Equal(t, "Forbidden", response.Error.Message)
}

// AssertBadRequest checks expected properties of BadRequest response.
func AssertBadRequest(t *testing.T, response *Response) {
	AssertFail(t, response)
	assert.Equal(t, 400, response.HTTPCode)
	assert.Equal(t, 400, response.Error.Code)
	assert.Equal(t, "Bad Request", response.Error.Message)
}
