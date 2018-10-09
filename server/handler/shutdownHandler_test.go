package handler

import (
	"net/http"
	"sync"
	"testing"
	"time"
)

type MockShutdowner struct {
	Called bool
}

func (m *MockShutdowner) Shutdown() {
	m.Called = true
}

func TestHandleShutdown(t *testing.T) {
	s := make([]Shutdowner, 1)
	mockShutdowner := new(MockShutdowner)
	s[0] = mockShutdowner
	wg := new(sync.WaitGroup)
	killFunc := func() {}
	h := NewShutdownHandler(s, wg, killFunc)
	h.shutdown()
	if !mockShutdowner.Called {
		t.Errorf("Shutdown function was not called")
	}
}

func TestHandleIdleRequest(t *testing.T) {
	s := make([]Shutdowner, 0)
	wg := new(sync.WaitGroup)
	killFuncCalled := false
	killFunc := func() {
		killFuncCalled = true
	}
	h := NewShutdownHandler(s, wg, killFunc)
	writer := new(MockResponseWriter)
	req := new(http.Request)
	h.HandleRequest(writer, req)
	time.Sleep(10 * time.Millisecond)
	if !killFuncCalled {
		t.Errorf("Kill function was not called")
	}
}

func TestHandleBusyRequest(t *testing.T) {
	s := make([]Shutdowner, 0)
	wg := new(sync.WaitGroup)
	wg.Add(1)
	killFunc := func() {}
	h := NewShutdownHandler(s, wg, killFunc)
	writer := new(MockResponseWriter)
	req := new(http.Request)
	go func() {
		time.Sleep(10 * time.Millisecond)
		wg.Done()
	}()
	h.HandleRequest(writer, req)
}
