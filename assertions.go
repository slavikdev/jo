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
