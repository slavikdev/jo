//
// Copyright (c) 2016 by Viacheslav Shynkarenko. All Rights Reserved.
//

package jo

import "gopkg.in/gin-gonic/gin.v1"

// Request contains request specific data.
type Request struct {
	// Context is a gin-specific request context.
	// TODO it's better to hide it and expose via public functions or properties.
	Context *gin.Context

	// GlobalContext is something user desided to pass to every request.
	// Typecally it's some sort of global configuration.
	GlobalContext interface{}

	// PrevHandlerResponse is a response of a previous handler.
	// For the first handler it's nil.
	PrevHandlerResponse *Response

	// Logger is a user defined logging interface.
	Logger ILogger
}

// GetQuery returns request query string value by specified argument name.
func (request *Request) GetQuery(name string) string {
	return request.Context.Query(name)
}

// GetJSON deserializes JSON request into specified object.
func (request *Request) GetJSON(output interface{}) {
	request.Context.BindJSON(output)
}
