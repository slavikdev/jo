//
// Copyright (c) 2016 by Viacheslav Shynkarenko. All Rights Reserved.
//

package jo

import "testing"

func TestSimpleRequest(t *testing.T) {
	api := NewAPI()
	api.Map("get", "/", emptyHandler)
	ht := NewHTTPTest(api)
	response := ht.Get("/")
	AssertOk(t, response)
}
