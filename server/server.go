/*
Package server provides a Server type; a hashing service that can be used to store and retrieve hashed values.
*/
package server

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/ifIMust/encodeServer/server/handler"
	"github.com/ifIMust/encodeServer/server/model"
)

/*
A Server implements an http server that stores and retrieves password hashes, and provides
basic usage statistics.
*/
type Server struct {
	// To ensure a clean server shutdown, receive a value from ShutdownComplete before program exit.
	ShutdownComplete chan int
	server           *http.Server
	hashHandler      *handler.HashHandler
}

func NewServer(port string) *Server {
	s := new(Server)
	s.ShutdownComplete = make(chan int)
	shutdownWaitGroup := new(sync.WaitGroup)
	mux := http.NewServeMux()
	stats := model.NewStats()
	s.hashHandler = handler.NewHashHandler(stats, shutdownWaitGroup)
	statsHandler := handler.NewStatsHandler(stats)
	handlers := []handler.Shutdowner{s.hashHandler, statsHandler}
	killFunc := func() {
		s.shutdown()
	}
	shutdownHandler := handler.NewShutdownHandler(handlers, shutdownWaitGroup, killFunc)

	mux.HandleFunc("/hash", getWrappedHandler(s.hashHandler.HandleRequest, shutdownWaitGroup))
	mux.HandleFunc("/hash/", getWrappedHandler(s.hashHandler.HandleRequest, shutdownWaitGroup))
	mux.HandleFunc("/stats", getWrappedHandler(statsHandler.HandleRequest, shutdownWaitGroup))
	mux.HandleFunc("/shutdown", shutdownHandler.HandleRequest)

	s.server = &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
	return s
}

/*
Run starts the web server, and blocks until the server has been signaled to shut down.
To guarantee a completely clean shutdown, receive a value from Server.ShutdownComplete
after this function returns.
Example:
<-server.Shutdowncomplete
*/
func (s *Server) Run() error {
	return s.server.ListenAndServe()
}

/*
SetDelay modifies the hashing delay used by the server. The default delay is 5s.
*/
func (s *Server) SetDelay(delay time.Duration) {
	s.hashHandler.SetDelay(delay)
}

func (s *Server) shutdown() error {
	err := s.server.Shutdown(context.Background())
	s.ShutdownComplete <- 1
	return err
}

func getWrappedHandler(handler func(http.ResponseWriter, *http.Request), wg *sync.WaitGroup) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, request *http.Request) {
		wg.Add(1)
		handler(w, request)
		wg.Done()
	}
}
