//
// Copyright (c) 2016 by Viacheslav Shynkarenko. All Rights Reserved.
//

package jo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Creates API instance and make sure all fields are empty right after creation.
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

// Creates API instance, makes sure global context is nil after creation,
// then adds dummy global contexts and checks whether it was properly stored.
func TestSetGlobalContext(t *testing.T) {
	api := NewAPI()
	assert.Nil(t, api.globalContext)

	globalContext := &dummyGlobalContext{Data: "test"}
	api.SetGlobalContext(globalContext)
	assert.NotNil(t, api.globalContext)
	assert.Equal(t, globalContext.Data, (api.globalContext.(*dummyGlobalContext)).Data)
}

// Creates API instance, makes sure init request handler is nil,
// then sets the handler to a function that returns empty Ok response,
// then calls the function and validates response.
func TestSetInitRequestHandler(t *testing.T) {
	api := NewAPI()
	assert.Nil(t, api.initRequestHandler)

	handlers := newTestHandlers()
	api.SetInitRequestHandler(handlers.emptyHandler)
	assert.NotNil(t, api.initRequestHandler)
	response := api.initRequestHandler(nil)
	AssertOk(t, response)
}

// Creates API instance, makes sure end request handler is nil,
// then sets the handler to a function that returns empty Ok response,
// then calls the function and validates response.
func TestSetEndRequestHandler(t *testing.T) {
	api := NewAPI()
	assert.Nil(t, api.endRequestHandler)

	handlers := newTestHandlers()
	api.SetEndRequestHandler(handlers.emptyHandler)
	assert.NotNil(t, api.endRequestHandler)
	response := api.endRequestHandler(nil)
	AssertOk(t, response)
}

// Creates API instance, makes sure logger is nil,
// then sets logger and checks whether it was properly stored.
func TestSetLogger(t *testing.T) {
	api := NewAPI()
	assert.Nil(t, api.logger)
	logger := &testLogger{}
	api.SetLogger(logger)
	assert.NotNil(t, api.logger)
}

// Create API instance, makes sure routes are empty.
// Adds 3 mappings with several HTTP methods which eventually maps handlers to 7 endpoints.
// Checks whether correct amount of routes was added.
func TestMap(t *testing.T) {
	api := NewAPI()
	assert.Empty(t, api.routes)

	handlers := newTestHandlers()
	api.Map("get", "/a", handlers.emptyHandler)
	api.Map("get,post,put", "/b", handlers.emptyHandler)
	api.Map("post, put, delete", "/c", handlers.emptyHandler)

	assert.NotEmpty(t, api.routes)
	assert.Equal(t, 7, len(api.routes), "There must be 7 routes for each HTTP method + path")
}

// Creates API, adds simple route mapping and builds gin engine, creating gin-specific
// handler for the route mapped.
// Makes sure the engine isn't nil.
func TestBuildEngine(t *testing.T) {
	api := NewAPI()
	handlers := newTestHandlers()
	api.Map("get", "/", handlers.emptyHandler)
	engine := api.buildEngine()
	assert.NotNil(t, engine)
}
