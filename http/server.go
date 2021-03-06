package http

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"pkg.glorieux.io/mantra"
	internalNet "pkg.glorieux.io/mantra/internal/net"
)

// RoutesFunc is a function defining HTTP routes
type RoutesFunc func(*mux.Router)

// Server is an HTTP server
type Server interface {
	mantra.Service

	URL() string
	Serve()
}

type httpServer struct {
	*http.Server

	listener net.Listener
}

// NewServer returns a supervised HTTP server
func NewServer(routes RoutesFunc) Server {
	router := mux.NewRouter()
	routes(router)
	return &httpServer{
		Server: &http.Server{
			Handler:        router,
			ReadTimeout:    5 * time.Second,
			WriteTimeout:   10 * time.Second,
			IdleTimeout:    120 * time.Second,
			MaxHeaderBytes: 1 << 20,
			// TODO Add SSL support
		},
		listener: internalNet.ListenerFromPort(8080),
	}
}

func (s *httpServer) Serve() {
	go func() {
		err := s.Server.Serve(s.listener)
		if err != http.ErrServerClosed {
			panic(err)
		}
	}()
}

func (s *httpServer) Stop() error {
	return s.Shutdown(context.Background())
}

func (s *httpServer) URL() string {
	return fmt.Sprintf("http://%s", s.listener.Addr().String())
}

func (*httpServer) String() string {
	return "HTTPServer"
}
