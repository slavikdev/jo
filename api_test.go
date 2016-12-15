//
// Copyright (c) 2016 by Viacheslav Shynkarenko. All Rights Reserved.
//

package jo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewAPI(t *testing.T) {
	api := NewAPI()
	assert.Nil(t, api.globalContext)
	assert.Empty(t, api.routes)
	assert.Nil(t, api.logger)
	assert.Nil(t, api.initRequestHandler)
	assert.Nil(t, api.endRequestHandler)
}

type dummyGlobalContext struct {
	Data string
}

func TestSetGlobalContext(t *testing.T) {
	api := NewAPI()
	assert.Nil(t, api.globalContext)

	globalContext := &dummyGlobalContext{Data: "test"}
	api.SetGlobalContext(globalContext)
	assert.NotNil(t, api.globalContext)
	assert.Equal(t, globalContext.Data, (api.globalContext.(*dummyGlobalContext)).Data)
}

func TestSetInitRequestHandler(t *testing.T) {
	api := NewAPI()
	assert.Nil(t, api.initRequestHandler)

	api.SetInitRequestHandler(emptyHandler)
	assert.NotNil(t, api.initRequestHandler)
	response := api.initRequestHandler(nil)
	assertOk(t, response)
}

func TestSetEndRequestHandler(t *testing.T) {
	api := NewAPI()
	assert.Nil(t, api.endRequestHandler)

	api.SetEndRequestHandler(emptyHandler)
	assert.NotNil(t, api.endRequestHandler)
	response := api.endRequestHandler(nil)
	assertOk(t, response)
}

func TestSetLogger(t *testing.T) {
	api := NewAPI()
	assert.Nil(t, api.logger)
	logger := &testLogger{}
	api.SetLogger(logger)
	assert.NotNil(t, api.logger)
}

func TestMap(t *testing.T) {
	api := NewAPI()
	assert.Empty(t, api.routes)

	api.Map("get", "/a", emptyHandler)
	api.Map("get,post,put", "/b", emptyHandler)
	api.Map("post, put, delete", "/c", emptyHandler)

	assert.NotEmpty(t, api.routes)
	assert.Equal(t, 7, len(api.routes), "There must be 7 routes for each HTTP method + path")
}

func assertOk(t *testing.T, response *Response) {
	assert.NotNil(t, response)
	assert.Equal(t, 200, response.HTTPCode)
	assert.True(t, response.Successful)
}

func emptyHandler(context *RequestContext) *Response {
	return Ok(nil)
}
