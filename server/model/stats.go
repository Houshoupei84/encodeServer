/*
Package model provides utility classes for storing web server data.
*/
package model

import (
	"encoding/json"
	"strings"
	"sync"
	"time"
)

/*
A Stats is a threadsafe tracker of total requests and average processing time.
*/
type Stats struct {
	requests  int
	totalTime time.Duration
	mutex     sync.Mutex
}

/*
A StatsReport is used for marshaling statistical output to JSON.
*/
type StatsReport struct {
	Total   int
	Average int
}

func NewStats() *Stats {
	return new(Stats)
}

func (s *Stats) AddRequest(t time.Duration) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.requests++
	s.totalTime += t
}

func (s *Stats) GetStats() StatsReport {
	averageTime := time.Duration(0)
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if s.requests != 0 {
		averageTime = s.totalTime / time.Duration(s.requests)
	}
	return StatsReport{s.requests, int(averageTime)}
}

func (s *Stats) GetStatsJson() string {
	report := s.GetStats()
	output, err := json.Marshal(report)
	if err != nil {
		return "ERROR"
	}
	result := strings.ToLower(string(output))
	return result
}
