package jo

// Response describes common service response format.
// Every HTTP response is wrapped in this structure.
type Response struct {
	Successful bool          `json:"successful"`
	Error      ResponseError `json:"error"`
	Data       interface{}   `json:"data"`

	HttpCode   int
	EndRequest bool
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
	response.HttpCode = 200
	response.EndRequest = true
	return response
}

// Fail creates unsuccessful response with error information.
func Fail(code int, data interface{}, errorData ResponseError) *Response {
	response := &Response{}
	response.Successful = false
	response.Error = errorData
	response.Data = data
	response.HttpCode = code
	response.EndRequest = true
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
