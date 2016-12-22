//
// Copyright (c) 2016 by Viacheslav Shynkarenko. All Rights Reserved.
//

package jo

// Response describes common service response format.
// Every HTTP response is wrapped in this structure.
// NOTE this structure isn't directly serialized to JSON anywhere except tests.
// Instead we take required fields because error field, for example,
// isn't always needed in response and it wounldn't be nice to bloat responses
// with redundant data.
type Response struct {
	Successful bool          `json:"successful"`
	Error      ResponseError `json:"error"`
	Data       interface{}   `json:"data"`

	HTTPCode int `json:"-"`

	// EndRequest specifies whether this response should be considered as final.
	EndRequest bool `json:"-"`
}

// ResponseError describes error information returned in response.
type ResponseError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// Ok creates successful response.
func Ok(data interface{}) *Response {
	response := &Response{}
	response.Successful = true
	response.Data = data
	response.HTTPCode = 200
	response.EndRequest = true
	return response
}

// Fail creates unsuccessful response with error information.
func Fail(code int, data interface{}, errorData ResponseError) *Response {
	response := &Response{}
	response.Successful = false
	response.Error = errorData
	response.Data = data
	response.HTTPCode = code
	response.EndRequest = true
	return response
}

// Forbidden creates 403 Forbidden HTTP response.
func Forbidden() *Response {
	return ForbiddenMessage("Forbidden")
}

// ForbiddenMessage creates 403 Forbidden HTTP response with specified message.
func ForbiddenMessage(message string) *Response {
	return createHTTPErrorResponse(403, message)
}

// BadRequest creates 400 Bad Request HTTP response.
func BadRequest() *Response {
	return BadRequestMessage("Bad Request")
}

// BadRequestMessage creates 400 Bad Request HTTP response with specified message.
func BadRequestMessage(message string) *Response {
	return createHTTPErrorResponse(400, message)
}

// Unauthorized creates 401 Unauthorized HTTP response.
func Unauthorized() *Response {
	return UnauthorizedMessage("Unauthorized")
}

// UnauthorizedMessage creates 401 Unauthorized HTTP response with specified message.
func UnauthorizedMessage(message string) *Response {
	return createHTTPErrorResponse(401, message)
}

// Error creates 500 internal error HTTP response.
func Error(err error) *Response {
	return ErrorMessage(err.Error())
}

// ErrorMessage creates 500 internal error HTTP response with specified message.
func ErrorMessage(message string) *Response {
	return createHTTPErrorResponse(500, message)
}

// Next creates response which indicates that next handler in chain should be called.
func Next(data ...interface{}) *Response {
	response := &Response{}
	response.Data = data
	response.Successful = true
	response.EndRequest = false
	return response
}

func createHTTPErrorResponse(code int, message string) *Response {
	errorCode := code
	errorData := ResponseError{Code: errorCode, Message: message}
	response := Fail(errorCode, nil, errorData)
	return response
}
