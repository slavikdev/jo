//
// Copyright (c) 2016 by Viacheslav Shynkarenko. All Rights Reserved.
//

package jo

// ILogger is a definition of an object capable of logging.
// When handling a request we quite often find ourselves in need
// to log some error or event. Implementation of such logging varies.
// So in this library we allow users to provide their own loggers,
// compatible with this interface.
type ILogger interface {
	Debug(format string, v ...interface{})
	Info(format string, v ...interface{})
	Warn(format string, v ...interface{})
	Error(format string, v ...interface{})
}
