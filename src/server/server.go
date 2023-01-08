package server

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

type Server struct {
	Port, Timeout int16
}

func (server *Server) createServer() *http.Server {
	return &http.Server{
		Addr:              fmt.Sprintf(":%d", server.Port),
		ReadHeaderTimeout: time.Duration(server.Timeout) * time.Second,
	}
}

func (server *Server) StartServer() {
	httpServer := server.createServer()
	log.Printf("Server running at port %d\n", server.Port)
	log.Fatal(httpServer.ListenAndServe())
}
