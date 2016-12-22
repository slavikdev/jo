package jo

// HTTPTest is a definition of HTTP testing API.
type HTTPTest interface {
	Get(url string) *Response
	Delete(url string) *Response
	Post(url string, requestJSON interface{}) *Response
	Put(url string, requestJSON interface{}) *Response
	Patch(url string, requestJSON interface{}) *Response
}
