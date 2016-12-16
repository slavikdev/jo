//
// Copyright (c) 2016 by Viacheslav Shynkarenko. All Rights Reserved.
//

package jo

import (
	"strings"

	"gopkg.in/gin-gonic/gin.v1"
)

// API provides functionality to map handler functions on specific HTTP routes.
type API struct {
	globalContext      interface{}
	routes             []*Route
	logger             ILogger
	initRequestHandler RouteHandler
	endRequestHandler  RouteHandler
}

// NewAPI creates new instance of API structure.
func NewAPI() *API {
	api := &API{}
	return api
}

// SetGlobalContext stores an object of any kind as global context.
// This context is then passed to every request.
// For example app configuration can be passed via global context to each handler.
func (api *API) SetGlobalContext(context interface{}) {
	api.globalContext = context
}

// SetInitRequestHandler sets a handler function which
// will be called before user defined handlers on every request specified in routes.
// This handler won't be called on routes without custom handlers.
func (api *API) SetInitRequestHandler(handler RouteHandler) {
	api.initRequestHandler = handler
}

// SetEndRequestHandler sets a handler function which
// will be called after user defined handlers on every request specified in routes.
// This handler won't be called on routes without custom handlers.
func (api *API) SetEndRequestHandler(handler RouteHandler) {
	api.endRequestHandler = handler
}

// SetLogger sets user defined logger to be made available for every request handler.
func (api *API) SetLogger(logger ILogger) {
	api.logger = logger
}

// Map assigns a chain of handlers to a specific URL path
// available via specific HTTP methods.
// List of HTTP methods should be specified as a string e.g. "get,post,put" or just "get".
// Path format is compatible with gin https://github.com/gin-gonic/gin#parameters-in-path
// Handlers must be specified in consequent order. They are called in that same order
// and response might be returned on every handler
// if it returns response with EndRequest flag.
func (api *API) Map(httpMethods string, path string, handlers ...RouteHandler) {
	httpMethodsSplit := strings.Split(httpMethods, ",")
	for _, httpMethod := range httpMethodsSplit {
		route := &Route{}
		route.Method = strings.ToLower(strings.TrimSpace(httpMethod))
		route.Path = path
		route.Handlers = handlers
		api.routes = append(api.routes, route)
	}
}

// Run starts API on specified TCP address.
func (api *API) Run(addr string) error {
	engine := api.buildEngine()
	return engine.Run(addr)
}

// RunTLS starts API on specified TCP address, serving requests via TLS.
// Certificate and key file paths must be specified.
func (api *API) RunTLS(addr string, certFile string, keyFile string) error {
	engine := api.buildEngine()
	return engine.RunTLS(addr, certFile, keyFile)
}

// RunUnix starts API on unix socket.
func (api *API) RunUnix(file string) error {
	engine := api.buildEngine()
	return engine.RunUnix(file)
}

// buildEngine creates instance of a gin engine and adds routes to it.
func (api *API) buildEngine() *gin.Engine {
	engine := gin.Default()
	api.buildRoutes(engine)
	return engine
}

// buildRoutes goes through a list of defined routes and builds them into gin engine.
func (api *API) buildRoutes(engine *gin.Engine) {
	for _, route := range api.routes {
		api.mapRoute(route, engine)
	}
}

// mapRoute creates gin-specific handler wrapped around specified handlers and maps in
// on route.
func (api *API) mapRoute(route *Route, engine *gin.Engine) {
	handlerWrapper := api.createHandlerWrapper(route.Handlers)
	mapRouteHandler(route, handlerWrapper, engine)
}

// createHandlerWrapper creates gin-specific handler wrapper.
func (api *API) createHandlerWrapper(handlers []RouteHandler) gin.HandlerFunc {
	return func(innerContext *gin.Context) {
		response, context := api.initRequest(innerContext)
		if response.EndRequest {
			api.endRequest(innerContext, context, response)
			return
		}

		for _, handler := range handlers {
			response = handler(context)
			if response.EndRequest {
				api.endRequest(innerContext, context, response)
				return
			}
			context.PrevHandlerResponse = response
		}
	}
}

// initRequest is called at the beginning of every handled request.
// initRequestHandler is called if specified.
func (api *API) initRequest(innerContext *gin.Context) (*Response, *RequestContext) {
	response := Next(nil)
	context := api.createRequestContext(innerContext, response)
	if api.initRequestHandler != nil {
		response = api.initRequestHandler(context)
		context.PrevHandlerResponse = response
	}
	return response, context
}

// endRequest is called at the end of every handled request.
// endRequestHandler is called if specified.
func (api *API) endRequest(
	innerContext *gin.Context,
	context *RequestContext,
	response *Response) {
	if api.endRequestHandler != nil {
		context.PrevHandlerResponse = response
		response = api.endRequestHandler(context)
	}
	// Every request eventually returns JSON no matter of its status.
	innerContext.JSON(response.HTTPCode, response)
	innerContext.Abort()
}

// createRequestContext creates context passed to request handlers.
func (api *API) createRequestContext(
	innerContext *gin.Context, prevResponse *Response) *RequestContext {
	context := &RequestContext{}
	context.Context = innerContext
	context.GlobalContext = api.globalContext
	context.PrevHandlerResponse = prevResponse
	context.Logger = api.logger
	return context
}

// mapRouteHandler maps gin handler on specific path and HTTP method.
func mapRouteHandler(route *Route, handler gin.HandlerFunc, engine *gin.Engine) {
	switch route.Method {
	case "get":
		engine.GET(route.Path, handler)
	case "post":
		engine.POST(route.Path, handler)
	case "put":
		engine.PUT(route.Path, handler)
	case "patch":
		engine.PATCH(route.Path, handler)
	case "delete":
		engine.DELETE(route.Path, handler)
	}
}
