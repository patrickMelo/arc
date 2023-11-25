package server

import (
	"crypto/md5"
	"encoding/hex"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"arc/vm"
)

type (
	httpServer struct {
		http.Handler
		runtime *vm.Runtime
	}
)

func (server *httpServer) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	var idString = request.RequestURI + strconv.FormatInt(time.Now().Unix(), 10)
	var requestHash = md5.Sum([]byte(idString))
	var requestID = hex.EncodeToString(requestHash[:])
	var commandLine = ""

	log.Printf("REQ(%s): %s", requestID, request.RequestURI)

	if (request.Method == http.MethodGet) && (request.URL.EscapedPath() == "/") {
		commandLine = request.URL.Query().Get("cmd")
	} else {
		commandLine = buildCommandLineFromREST(request)
	}

	if commandLine == "" {
		log.Printf("RESP(%s): 400", requestID)
		response.WriteHeader(400)
		return
	}

	log.Printf("REQ(%s): %s", requestID, commandLine)

	if result := server.runtime.Execute(commandLine); result != nil {
		var resultString = strings.Join(result, " ")
		log.Printf("RESP(%s): %s", requestID, resultString)
		response.Write([]byte(resultString))
	} else {
		log.Printf("RESP(%s): 500", requestID)
		response.WriteHeader(500)
	}
}
