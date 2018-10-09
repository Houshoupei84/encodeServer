package handler

import (
	"testing"
	"time"

	"github.com/ifIMust/encodeServer/server/model"
)

func TestAddOneRequest(t *testing.T) {
	stats := model.NewStats()
	h := NewStatsHandler(stats)
	duration := time.Duration(50)
	h.addRequest(duration)
	requests, meanTime := h.getStats()
	if requests != 1 {
		t.Errorf("Expected %d requests, had %d", 1, requests)
	}
	if meanTime != 50 {
		t.Errorf("Expected %d duration, had %d", duration, meanTime)
	}
}

func TestAddTwoRequest(t *testing.T) {
	stats := model.NewStats()
	h := NewStatsHandler(stats)
	h.addRequest(time.Duration(50))
	h.addRequest(time.Duration(100))
	requests, meanTime := h.getStats()
	if requests != 2 {
		t.Errorf("Expected %d requests, had %d", 2, requests)
	}
	expectedTime := 75
	if meanTime != expectedTime {
		t.Errorf("Expected %d duration, had %d", expectedTime, meanTime)
	}
}

func TestAddManyRequests(t *testing.T) {
	stats := model.NewStats()
	h := NewStatsHandler(stats)
	madeRequests := 1001
	for i := 0; i < madeRequests; i++ {
		h.addRequest(time.Duration(i))
	}
	requests, meanTime := h.getStats()

	if requests != madeRequests {
		t.Errorf("Expected %d requests, had %d", madeRequests, requests)
	}
	expectedTime := 500
	if meanTime != expectedTime {
		t.Errorf("Expected %d duration, had %d", expectedTime, meanTime)
	}
}
