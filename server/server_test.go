package server

import (
	"net/http"
	"net/url"
	"strconv"
	"testing"
	"time"
)

const (
	host             = "localhost"
	port             = "8081"
	serverStartDelay = 5 * time.Millisecond
)

/*
These functional tests of the server are intended to validate that the Server shuts down
gracefully and cleanly, and that it handles poorly timed requests correctly.

Since an http.Server.ListenAndServe() call does not return until shutdown,
time.Sleep(serverStartDelay) calls are used to yield the execution context so that
the goroutine doing Server.Run() can execute.
This could result in test failures; in practice, the scheduler tends to do the right thing.
*/

func doGet(t *testing.T, id int) string {
	resp, err := http.Get("http://" + host + ":" + port + "/hash/" + strconv.Itoa(id))
	if err != nil {
		t.Errorf("Failed get request")
		return ""
	}
	body := make([]byte, 128)
	n, _ := resp.Body.Read(body)
	if n == 0 {
		t.Errorf("Failed read: %d bytes read", n)
	}
	resp.Body.Close()
	return string(body)
}

func doPost(t *testing.T) int {
	resp, err := http.PostForm("http://"+host+":"+port+"/hash",
		url.Values{"password": {"neato mosquito"}})
	if err != nil {
		t.Errorf("Post produced error %s", err)
		return -1
	}
	body := make([]byte, 1)
	n, err := resp.Body.Read(body)
	resp.Body.Close()
	if n != 1 {
		t.Errorf("Read incorrect body length %d", n)
	}
	result, err := strconv.Atoi(string(body[0:n]))
	if err != nil {
		t.Errorf("Response conversion produced error %s", err)
	}
	if result != 1 {
		t.Errorf("Post produced response %d", result)
	}
	return result
}

func doStats(t *testing.T) string {
	resp, err := http.Get("http://" + host + ":" + port + "/stats")
	if err != nil {
		t.Errorf("Stats produced error %s", err)
		return ""
	}
	body := make([]byte, 128)
	n, _ := resp.Body.Read(body)
	if n == 0 {
		t.Errorf("Failed read: %d bytes read", n)
	}
	resp.Body.Close()
	return string(body[0:n])
}

func doStatsUnavailable(t *testing.T) {
	_, err := http.Get("http://" + host + ":" + port + "/stats")
	if err == nil {
		t.Errorf("Expected error, didn't get one!")
	}
}

func doShutdown() {
	http.Get("http://" + host + ":" + port + "/shutdown")
}

func TestHandleShutdown(t *testing.T) {
	s := NewServer(port)
	go s.Run()
	time.Sleep(serverStartDelay)
	doShutdown()
	<-s.ShutdownComplete
}

func TestHandleDoubleShutdown(t *testing.T) {
	s := NewServer(port)
	go s.Run()
	time.Sleep(serverStartDelay)
	doShutdown()
	doShutdown()
	<-s.ShutdownComplete
}

func TestHandlePostRequest(t *testing.T) {
	s := NewServer(port)
	s.SetDelay(1 * time.Microsecond)
	go s.Run()
	time.Sleep(serverStartDelay)
	result := doPost(t)
	if result != 1 {
		t.Errorf("Bad result")
	}
	doShutdown()
	<-s.ShutdownComplete
}

func TestHandleDoubleShutdownWhileProcessing(t *testing.T) {
	s := NewServer(port)
	s.SetDelay(1 * time.Microsecond)
	go s.Run()
	time.Sleep(serverStartDelay)
	result := doPost(t)
	if result != 1 {
		t.Errorf("Bad result")
	}
	doShutdown()
	doShutdown()
	<-s.ShutdownComplete
}

func TestHandlePostGetShutdown(t *testing.T) {
	s := NewServer(port)
	s.SetDelay(1 * time.Microsecond)
	go s.Run()
	time.Sleep(serverStartDelay)
	postResult := doPost(t)
	getResult := doGet(t, postResult)
	if getResult == "" {
		t.Errorf("Empty response")
	}
	doShutdown()
	<-s.ShutdownComplete
}

func TestHandleStatsWithoutHashes(t *testing.T) {
	s := NewServer(port)
	go s.Run()
	time.Sleep(serverStartDelay)
	data := doStats(t)
	expected := "{\"total\":0,\"average\":0}"
	if data != expected {
		t.Errorf("Expected %s got %s", expected, data)
	}
	doShutdown()
	<-s.ShutdownComplete
}

func TestHandleStatsAfterShutdown(t *testing.T) {
	s := NewServer(port)
	go s.Run()
	time.Sleep(serverStartDelay)
	doShutdown()
	doStatsUnavailable(t)
	<-s.ShutdownComplete
}
