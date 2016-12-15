//
// Copyright (c) 2016 by Viacheslav Shynkarenko. All Rights Reserved.
//

package jo

func emptyHandler(context *RequestContext) *Response {
	return Ok(nil)
}

func passGlobalContext(context *RequestContext) *Response {
	return Ok(context.GlobalContext)
}

func loggingHanler(context *RequestContext) *Response {
	context.Logger.Debug("Hello %s", "World")
	context.Logger.Info("Hello %s", "World")
	context.Logger.Warn("Hello %s", "World")
	context.Logger.Error("Hello %s", "World")
	return Ok(nil)
}

func authHandler(context *RequestContext) *Response {
	if context.Context.Query("token") == "secret" {
		return Next(nil)
	}

	err := ResponseError{Code: 403, Message: "Forbidden"}
	return Fail(err.Code, nil, err)
}
