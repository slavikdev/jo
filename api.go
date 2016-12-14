package jo

import (
	"strings"

	"gopkg.in/gin-gonic/gin.v1"
)

type Api struct {
	globalContext      interface{}
	routes             []*Route
	logger             *ILogger
	initRequestHandler RouteHandler
	endRequestHandler  RouteHandler
}

func NewApi() *Api {
	api := &Api{}
	return api
}

func (api *Api) SetGlobalContext(context interface{}) {
	api.globalContext = context
}

func (api *Api) SetInitRequestHandler(handler RouteHandler) {
	api.initRequestHandler = handler
}

func (api *Api) SetEndRequestHandler(handler RouteHandler) {
	api.endRequestHandler = handler
}

func (api *Api) SetLogger(logger *ILogger) {
	api.logger = logger
}

func (api *Api) Map(httpMethods string, path string, handlers ...RouteHandler) {
	httpMethodsSplit := strings.Split(httpMethods, ",")
	for _, httpMethod := range httpMethodsSplit {
		route := &Route{}
		route.Method = strings.ToLower(strings.TrimSpace(httpMethod))
		route.Path = path
		route.Handlers = handlers
		api.routes = append(api.routes, route)
	}
}

func (api *Api) Run(addr string) error {
	engine := api.buildEngine()
	return engine.Run(addr)
}

func (api *Api) RunTLS(addr string, certFile string, keyFile string) error {
	engine := api.buildEngine()
	return engine.RunTLS(addr, certFile, keyFile)
}

func (api *Api) RunUnix(file string) error {
	engine := api.buildEngine()
	return engine.RunUnix(file)
}

func (api *Api) buildEngine() *gin.Engine {
	engine := gin.Default()
	api.buildRoutes(engine)
	return engine
}

func (api *Api) buildRoutes(engine *gin.Engine) {
	for _, route := range api.routes {
		api.mapRoute(route, engine)
	}
}

func (api *Api) mapRoute(route *Route, engine *gin.Engine) {
	handlerWrapper := api.createHandlerWrapper(route.Handlers)
	mapRouteHandler(route, handlerWrapper, engine)
}

func (api *Api) createHandlerWrapper(handlers []RouteHandler) gin.HandlerFunc {
	return func(innerContext *gin.Context) {
		prevResponse, context := api.initRequest(innerContext)
		for _, handler := range handlers {
			response := handler(context)
			if response.EndRequest {
				api.endRequest(innerContext, context, response)
				return
			}

			prevResponse = response
			context.PrevHandlerResponse = prevResponse
		}
	}
}

func (api *Api) initRequest(innerContext *gin.Context) (*Response, *RequestContext) {
	response := Next(nil)
	context := api.createRequestContext(innerContext, response)
	if api.initRequestHandler != nil {
		response = api.initRequestHandler(context)
		context.PrevHandlerResponse = response
	}
	return response, context
}

func (api *Api) endRequest(
	innerContext *gin.Context,
	context *RequestContext,
	response *Response) {
	if api.endRequestHandler != nil {
		context.PrevHandlerResponse = response
		response = api.endRequestHandler(context)
	}
	innerContext.JSON(response.HttpCode, response)
	innerContext.Abort()
}

func (api *Api) createRequestContext(
	innerContext *gin.Context, prevResponse *Response) *RequestContext {
	context := &RequestContext{}
	context.Context = innerContext
	context.GlobalContext = api.globalContext
	context.PrevHandlerResponse = prevResponse
	context.Logger = api.logger
	return context
}

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
