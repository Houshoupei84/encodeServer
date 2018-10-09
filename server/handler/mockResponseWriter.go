package handler

import "net/http"

/*
A MockResponseWriter is used for unit testing http handler functions.
*/
type MockResponseWriter struct {
	LastData []byte
}

func (m *MockResponseWriter) Header() http.Header {
	return make(http.Header)
}

func (m *MockResponseWriter) Write(data []byte) (int, error) {
	m.LastData = data
	return len(data), nil
}

func (m *MockResponseWriter) WriteHeader(int) {
}
