# JO: Opinionated Go library to build RESTful JSON APIs.

Jo wraps around [Gin](https://github.com/gin-gonic/gin) library and implements common patterns
useful in creating JSON APIs, such as strict response structure, authorization, logging.
Basically I've extracted here all the stuff I usually add building APIs. It wouldn't fit
everyone but that's exactly the point: to agree on common things and let developers
focus on specific implementation. If you need more flexibility–go get [gin](https://github.com/gin-gonic/gin).

[![Travis Build Status](https://travis-ci.org/slavikdev/jo.svg)](https://travis-ci.org/slavikdev/jo)
[![Appveyor Build Status](https://ci.appveyor.com/api/projects/status/h90m552en8cjxrv0?svg=true)](https://ci.appveyor.com/project/slavikdev/jo)
[![codecov](https://codecov.io/gh/slavikdev/jo/branch/master/graph/badge.svg)](https://codecov.io/gh/slavikdev/jo)
[![Go Report Card](https://goreportcard.com/badge/github.com/slavikdev/jo)](https://goreportcard.com/report/github.com/slavikdev/jo)
[![GoDoc](https://godoc.org/github.com/slavikdev/jo?status.svg)](https://godoc.org/github.com/slavikdev/jo)
