package server

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"
)

// These are used to convert a REST request to a database command line.

type (
	commandBuilder func(method string, parameters []string) string

	requestInterpreter struct {
		validMethods []string
		builder      commandBuilder
	}
)

var (
	validRequests = map[string]requestInterpreter{
		"db": {
			validMethods: []string{http.MethodGet},
			builder:      dbBuilder,
		},
		"values": {
			validMethods: []string{http.MethodGet, http.MethodPut, http.MethodPatch, http.MethodDelete},
			builder:      valuesBuilder,
		},
		"sets": {
			validMethods: []string{http.MethodGet, http.MethodPut},
			builder:      setsBuilder,
		},
	}
)

func buildCommandLineFromREST(request *http.Request) (commandLine string) {
	// Join the command parts: url path + paramters + data

	var commandParts = strings.Split(request.URL.EscapedPath(), "/")[1:]

	for httpName, httpValue := range request.URL.Query() {
		commandParts = append(commandParts, []string{httpName, httpValue[0]}...)
	}

	var bodyBuffer = new(bytes.Buffer)

	if bodySize, err := bodyBuffer.ReadFrom(request.Body); (err == nil) && (bodySize > 0) {
		commandParts = append(commandParts, strings.Split(bodyBuffer.String(), " ")...)
	}

	// Get the root, check method and build the command.

	if len(commandParts) < 1 {
		return
	}

	var root = commandParts[0]
	var parameters = commandParts[1:]

	if interpreter, exists := validRequests[root]; exists {
		for index := range interpreter.validMethods {
			if request.Method == interpreter.validMethods[index] {
				commandLine = interpreter.builder(request.Method, parameters)
			}
		}
	}

	return
}

/*

DBSIZE
======
GET /db/size

*/

const (
	dbSizeParemeterName = "size"
)

func dbBuilder(method string, parameters []string) (command string) {
	if (len(parameters) == 1) && (strings.ToLower(parameters[0]) == dbSizeParemeterName) {
		command = "DBSIZE"
	}

	return
}

/*

SET​ key value
=============
PUT /values
key value

SET key value EX seconds
========================
PUT /values
key value seconds

GET​ key
=======
GET /values/key

DEL​ key [key ...]
=================
DELETE /values/key

INCR​ key
========
PATCH /values/key

*/

func valuesBuilder(method string, parameters []string) (command string) {
	switch method {
	case http.MethodGet:
		if len(parameters) == 1 {
			command = fmt.Sprintf("GET %s", parameters[0])
		}
	case http.MethodPut:
		switch len(parameters) {
		case 2:
			command = fmt.Sprintf("SET %s %s", parameters[0], parameters[1])
		case 3:
			command = fmt.Sprintf("SET %s %s EX %s", parameters[0], parameters[1], parameters[2])
		}
	case http.MethodPatch:
		if len(parameters) == 1 {
			command = fmt.Sprintf("INCR %s", parameters[0])
		}
	case http.MethodDelete:
		if len(parameters) == 1 {
			command = fmt.Sprintf("DEL %s", parameters[0])
		}
	}

	return
}

/*

ZADD​ key score member
=====================
PUT /sets
key score member

ZCARD​ key
=========
GET /sets/key/size

ZRANK​ key member
================
GET /sets/key/rank/member

ZRANGE key start stop
=====================
GET /sets/key
GET /sets/key?start=0&stop=2

*/

const (
	setSizeParameterName  = "size"
	setRankParameterName  = "rank"
	setStartParameterName = "start"
	setStopParameterName  = "stop"
)

func setsBuilder(method string, parameters []string) (command string) {
	switch method {
	case http.MethodGet:
		switch len(parameters) {
		case 1:
			command = fmt.Sprintf("ZRANGE %s 0 -1", parameters[0])
		case 2:
			if parameters[1] == setSizeParameterName {
				command = fmt.Sprintf("ZCARD %s", parameters[0])
			}
		case 3:
			switch parameters[1] {
			case setRankParameterName:
				command = fmt.Sprintf("ZRANK %s %s", parameters[0], parameters[2])
			case setStartParameterName:
				command = fmt.Sprintf("ZRANGE %s %s -1", parameters[0], parameters[2])
			case setStopParameterName:
				command = fmt.Sprintf("ZRANGE %s 0 %s", parameters[0], parameters[2])
			}
		case 5:
			if ((parameters[1] == setStartParameterName) && (parameters[3] == setStopParameterName)) ||
				((parameters[1] == setStopParameterName) && (parameters[3] == setStartParameterName)) {

				var startIndex = 2
				var stopIndex = 4

				if parameters[1] == setStopParameterName {
					startIndex = 4
					stopIndex = 2
				}

				command = fmt.Sprintf("ZRANGE %s %s %s", parameters[0], parameters[startIndex], parameters[stopIndex])
			}
		}
	case http.MethodPut:
		if (len(parameters) >= 3) && (len(parameters)%2 == 1) {
			command = "ZADD " + strings.Join(parameters, " ")
		}
	}

	return
}
