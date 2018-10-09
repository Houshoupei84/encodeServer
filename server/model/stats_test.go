package model

import "testing"

func TestGetStatsZero(t *testing.T) {
	s := NewStats()
	report := s.GetStats()
	if report.Total != 0 {
		t.Errorf("Expected report.requests to be %d, was %d", 0, report.Total)
	}
	if report.Average != 0 {
		t.Errorf("Expected report.averageTimeMicroseconds to be %d, was %d",
			0, report.Average)
	}
}

func TestGetStatsJson(t *testing.T) {
	s := NewStats()
	statsJson := s.GetStatsJson()
	expected := "{\"total\":0,\"average\":0}"
	if statsJson != expected {
		t.Errorf("Bad outputs. Expected '%s' got '%s'", expected, statsJson)
	}
}

func TestGetStatsTwoRequests(t *testing.T) {
	s := NewStats()
	s.AddRequest(100)
	s.AddRequest(200)
	report := s.GetStats()
	if report.Total != 2 {
		t.Errorf("Expected report.requests to be %d, was %d", 0, report.Total)
	}
	if report.Average != 150 {
		t.Errorf("Expected report.averageTimeMicroseconds to be %d, was %d",
			0, report.Average)
	}
}
