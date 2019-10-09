package http

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"pkg.glorieux.io/mantra"
)

// RoutesFunc is a function defining HTTP routes
type RoutesFunc func(*mux.Router)

// Server is an HTTP server
type Server interface {
	mantra.Service

	URL() string
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
		listener: listener(),
	}
}

func listener() net.Listener {
	port := 4242

	for {
		ln, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
		if err == nil {
			return ln
		}
		port++
	}
}

func (s *httpServer) Receive(mux mantra.ServeMux) {}

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
