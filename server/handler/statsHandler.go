package handler

import (
	"io"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/ifIMust/encodeServer/server/model"
)

/*
A StatsHandler responds to GET requests to the '/stats' endpoint.
The response contains the sum and average processing time of POST requests to the '/hash/' endpoint.
*/
type StatsHandler struct {
	stats *model.Stats
	run   atomic.Value
}

/*
NewStatsHandler initializes and returns a new StatsHandler.
Parameter stats will be used by the new object for reading and writing statistics.
*/
func NewStatsHandler(stats *model.Stats) *StatsHandler {
	s := new(StatsHandler)
	s.stats = stats
	s.run.Store(true)
	return s
}

/*
HandleRequest is an http request handler intended for use with http.ServeMux.
*/
func (s *StatsHandler) HandleRequest(w http.ResponseWriter, request *http.Request) {
	if s.run.Load().(bool) {
		if request.Method == "GET" {
			output := s.stats.GetStatsJson()
			io.WriteString(w, string(output))
		} else {
			http.NotFound(w, request)
		}
	} else {
		w.WriteHeader(http.StatusServiceUnavailable)
	}
}

/*
Shutdown disables further handling of requests by this handler.
*/
func (s *StatsHandler) Shutdown() {
	s.run.Store(false)
}

func (s *StatsHandler) addRequest(t time.Duration) {
	s.stats.AddRequest(t)
}

func (s *StatsHandler) getStats() (int, int) {
	report := s.stats.GetStats()
	return report.Total, report.Average
}
