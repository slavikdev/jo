//
// Copyright (c) 2016 by Viacheslav Shynkarenko. All Rights Reserved.
//

package jo

// RouteHandler is a definition of a function which handles request on specific route.
type RouteHandler func(context *RequestContext) *Response

// Route is a definition of API endpoint with specific path,
// HTTP method and chain of handlers.
type Route struct {
	Method   string
	Path     string
	Handlers []RouteHandler
}
