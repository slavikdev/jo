//
// Copyright (c) 2016 by Viacheslav Shynkarenko. All Rights Reserved.
//

package jo

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOk(t *testing.T) {
	api, _, http := newAPITest()
	path := "/ok"
	api.Map("get", path, func(r *Request) *Response {
		return Ok(true)
	})
	response := http.Get(path)
	AssertOk(t, response)
}

func TestFail(t *testing.T) {
	api, _, http := newAPITest()
	httpCode := 501
	data := "something"
	errCode := 1
	errMessage := "error"
	path := "/fail"
	api.Map("get", path, func(r *Request) *Response {
		return Fail(httpCode, data, ResponseError{Code: errCode, Message: errMessage})
	})
	response := http.Get(path)
	AssertFail(t, response)
	assert.Equal(t, httpCode, response.HTTPCode)
	assert.Equal(t, data, response.Data.(string))
	assert.Equal(t, errCode, response.Error.Code)
	assert.Equal(t, errMessage, response.Error.Message)
}

func TestForbidden(t *testing.T) {
	api, _, http := newAPITest()
	path := "/forbidden"
	api.Map("get", path, func(r *Request) *Response {
		return Forbidden()
	})
	response := http.Get(path)
	AssertForbidden(t, response)
}

func TestBadRequest(t *testing.T) {
	api, _, http := newAPITest()
	path := "/bad-request"
	api.Map("get", path, func(r *Request) *Response {
		return BadRequest()
	})
	response := http.Get(path)
	AssertBadRequest(t, response)
}

func TestUnauthorized(t *testing.T) {
	api, _, http := newAPITest()
	path := "/unauthorized"
	api.Map("get", path, func(r *Request) *Response {
		return Unauthorized()
	})
	response := http.Get(path)
	AssertUnauthorized(t, response)
}

func TestError(t *testing.T) {
	api, _, http := newAPITest()
	path := "/error"
	errMessage := "error occured"
	api.Map("get", path, func(r *Request) *Response {
		err := errors.New(errMessage)
		return Error(err)
	})
	response := http.Get(path)
	AssertResponseError(t, response, errMessage)
}
