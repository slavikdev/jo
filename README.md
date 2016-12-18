# JO: Opinionated Go library to build RESTful JSON APIs.

Jo wraps around [gin](https://github.com/gin-gonic/gin) library and implements common patterns
useful in creating JSON APIs, such as strict response structure, authorization, logging, functional testing.
Basically I've extracted all the stuff I usually add building APIs. It wouldn't fit
everyone but that's exactly the point: to agree on common things once and use them as a framework.
If you need more flexibility–go get [gin](https://github.com/gin-gonic/gin).

[![Travis Build Status](https://travis-ci.org/slavikdev/jo.svg)](https://travis-ci.org/slavikdev/jo)
[![Appveyor Build Status](https://ci.appveyor.com/api/projects/status/h90m552en8cjxrv0?svg=true)](https://ci.appveyor.com/project/slavikdev/jo)
[![codecov](https://codecov.io/gh/slavikdev/jo/branch/master/graph/badge.svg)](https://codecov.io/gh/slavikdev/jo)
[![Go Report Card](https://goreportcard.com/badge/github.com/slavikdev/jo)](https://goreportcard.com/report/github.com/slavikdev/jo)
[![GoDoc](https://godoc.org/github.com/slavikdev/jo?status.svg)](https://godoc.org/github.com/slavikdev/jo)

## Example

```go
package main

import (
	"time"

	"github.com/slavikdev/jo"
)

func main() {
	// Create api.
	api := jo.NewAPI()

	// Map routes.
	api.Map("get", "/time", getTime)
	api.Map("get", "/secret", auth, getSecret)

	// Start api on port 9999.
	err := api.Run(":9999")
	if err != nil {
		panic(err)
	}
}

// Returns successful response with current time in data field.
func getTime(rc *jo.RequestContext) *jo.Response {
	currentTime := time.Now()
	return jo.Ok(currentTime)
}

// Returns successful response with a word "secret" in data field.
func getSecret(rc *jo.RequestContext) *jo.Response {
	return jo.Ok("secret")
}

const apiToken string = "0123456789"

// Checks if query string has apiToken argument with specific value.
// If it does--passes request to next handler. Otherwise returns 403 Forbidden error.
func auth(rc *jo.RequestContext) *jo.Response {
	if rc.Context.Query("apiToken") == apiToken {
		return jo.Next(nil)
	}
	return jo.Forbidden()
}
```

The example works as follows:
```
GET localhost:9999/time
```
```
HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8
Date: Sun, 18 Dec 2016 00:51:43 GMT
Content-Length: 64

{
    "data": "2016-12-18T02:51:43.07980668+02:00",
    "successful": true
}
```

```
GET localhost:9999/secret
```
```
HTTP/1.1 403 Forbidden
Content-Type: application/json; charset=utf-8
Date: Sun, 18 Dec 2016 00:51:10 GMT
Content-Length: 88

{
    "data": null,
    "error": {
        "code": 403,
        "message": "Forbidden",
        "data": null
    },
    "successful": false
}
```

```
GET localhost:9999/secret?apiToken=0123456789
```
```
HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8
Date: Sun, 18 Dec 2016 00:40:37 GMT
Content-Length: 36

{
    "data": "secret",
    "successful": true
}
```
