package handler

import (
	"net/http"
	"sync"
)

/*
A ShutdownHandler responds to any request by cleanly bringing down the web seb server.
*/
type ShutdownHandler struct {
	duties            []Shutdowner
	shutdownWaitGroup *sync.WaitGroup
	killFunc          func()
}

/*
NewShutdownHandler initializes and returns a new ShutdownHandler.
When it handles a shutdown request, it will call Shutdown() on each entry in handlers,
call Wait() on shutdownWaitGroup, and invoke killFunc as a new goroutine.
*/
func NewShutdownHandler(handlers []Shutdowner, shutdownWaitGroup *sync.WaitGroup, killFunc func()) *ShutdownHandler {
	s := new(ShutdownHandler)
	s.duties = handlers
	s.shutdownWaitGroup = shutdownWaitGroup
	s.killFunc = killFunc
	return s
}

func (s *ShutdownHandler) shutdown() {
	for _, s := range s.duties {
		s.Shutdown()
	}
	s.shutdownWaitGroup.Wait()
}

/*
HandleRequest is an http request handler intended for use with http.ServeMux.
*/
func (s *ShutdownHandler) HandleRequest(w http.ResponseWriter, request *http.Request) {
	w.WriteHeader(http.StatusOK)
	s.shutdown()
	// Detach this goroutine so that the handler can complete, allowing the server to shutdown.
	// In the normal usage, this function is bound to http.Server.Shutdown, The completion of the
	// function is signaled by a send to the Server.ShutdownComplete channel.
	go s.killFunc()
}
