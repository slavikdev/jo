//
// Copyright (c) 2016 by Viacheslav Shynkarenko. All Rights Reserved.
//

package jo

import "log"

type testLogger struct {
}

func (l *testLogger) Debug(format string, v ...interface{}) {
	l.log("debug", format, v)
}

func (l *testLogger) Info(format string, v ...interface{}) {
	l.log("info", format, v)
}

func (l *testLogger) Warn(format string, v ...interface{}) {
	l.log("warn", format, v)
}

func (l *testLogger) Error(format string, v ...interface{}) {
	l.log("error", format, v)
}

func (l *testLogger) log(category string, format string, v ...interface{}) {
	log.Printf("["+category+"] "+format, v)
}
