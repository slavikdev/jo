package jo

import "gopkg.in/gin-gonic/gin.v1"

// RequestContext contains request specific data and implements useful methods
// to use on current request.
type RequestContext struct {
	Context             *gin.Context
	GlobalContext       interface{}
	PrevHandlerResponse *Response
	Logger              *ILogger
}
