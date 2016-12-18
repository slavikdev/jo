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
