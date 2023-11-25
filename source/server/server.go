package server

import (
	"net/http"

	"arc/vm"
)

type (
	// Server defines a simple database server type.
	Server struct {
		address string
		runtime *vm.Runtime
	}
)

// Create creates a new database server.
func Create(address string, runtime *vm.Runtime) (server *Server) {
	return &Server{
		address: address,
		runtime: runtime,
	}
}

// Run starts the server and listens for connections and commands.
func (server *Server) Run() (err error) {
	return http.ListenAndServe(server.address, &httpServer{
		runtime: server.runtime,
	})
}
