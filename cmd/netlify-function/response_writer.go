package main

import "net/http"

// Placeholder struct implementing http.ResponseWriter
// to transport between http and lambda Request/Response
type respWriter struct {
	header     http.Header
	b          []byte
	statusCode int
}

func (w *respWriter) Header() http.Header {
	return w.header
}

func (w *respWriter) Write(b []byte) (int, error) {
	w.b = b
	return len(b), nil
}

func (w *respWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
}
