package handler

import (
	"bytes"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/ifIMust/encodeServer/server/model"
)

func TestGetNonexistantHash(t *testing.T) {
	stats := model.NewStats()
	wg := new(sync.WaitGroup)
	h := NewHashHandler(stats, wg)
	id := 77
	retrieved := h.getHash(id)
	expected := ""
	if retrieved != expected {
		t.Errorf("Expected %s, retrieved %s", expected, retrieved)
	}
}

func TestAddGetHash(t *testing.T) {
	stats := model.NewStats()
	wg := new(sync.WaitGroup)
	h := NewHashHandler(stats, wg)
	id := 77
	pwd := "3.14159265358979"
	setDelay := 1 * time.Millisecond
	h.SetDelay(setDelay)
	h.waitGroup.Add(1)
	go h.delayedHash(id, pwd)
	time.Sleep(setDelay * 2)
	retrieved := h.getHash(id)
	if len(retrieved) == 0 {
		t.Errorf("Expected nonzero hash string")
	}
}

func TestGetHashTooQuickly(t *testing.T) {
	stats := model.NewStats()
	wg := new(sync.WaitGroup)
	h := NewHashHandler(stats, wg)
	id := 77
	pwd := "777"
	setDelay := 10 * time.Millisecond
	h.SetDelay(setDelay)
	h.waitGroup.Add(1)
	go h.delayedHash(id, pwd)
	retrieved := h.getHash(id)
	if retrieved != "" {
		t.Errorf("Hash was processed before delay expired.")
	}
}

func TestHandlePost(t *testing.T) {
	stats := model.NewStats()
	wg := new(sync.WaitGroup)
	h := NewHashHandler(stats, wg)
	writer := new(MockResponseWriter)
	body := bytes.NewBufferString("malformed post")
	req, err := http.NewRequest("POST", "http://12.34.56.78:4321/hash", body)
	if err != nil {
		t.Errorf("Failed to construct POST request")
	}
	h.HandleRequest(writer, req)
	if len(writer.LastData) == 0 {
		t.Errorf("Handler didn't write anything")
	}
}
