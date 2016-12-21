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
	errorCode := 403
	errorData := ResponseError{Code: errorCode, Message: "Forbidden"}
	response := Fail(errorCode, nil, errorData)
	return response
}

// BadRequest creates 400 Bad Request HTTP response.
func BadRequest() *Response {
	errorCode := 400
	errorData := ResponseError{Code: errorCode, Message: "Bad Request"}
	response := Fail(errorCode, nil, errorData)
	return response
}

// Unauthorized creates 401 Unauthorized HTTP response.
func Unauthorized() *Response {
	errorCode := 401
	errorData := ResponseError{Code: errorCode, Message: "Unauthorized"}
	response := Fail(errorCode, nil, errorData)
	return response
}

// Error creates 500 internal error HTTP response.
func Error(err error) *Response {
	errorCode := 500
	errorData := ResponseError{Code: errorCode, Message: err.Error()}
	response := Fail(errorCode, nil, errorData)
	return response
}

// Next creates response for the next handler in chain.
func Next(data interface{}) *Response {
	response := &Response{}
	response.Data = data
	response.Successful = true
	response.EndRequest = false
	return response
}
