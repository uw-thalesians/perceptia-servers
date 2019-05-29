package handler

import (
	"net/http"
)

// RequestLogger represents the current handler in the request/response cycle.
type RequestLogger struct {
	handler http.Handler
	cx      *Context
}

// NewRequestLogger constructs a new RequestLogger struct with the provided handler and Context.
func (cx *Context) NewRequestLogger(handler http.Handler) http.Handler {
	return &RequestLogger{handler, cx}
}

// ServeHTTP logs all requests if in development environment
func (gv *RequestLogger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if gv.cx.environment == "development" {
		_ = gv.cx.logger.Log("requestAgent", r.UserAgent(), "requestMethod", r.Method, "requestPath", r.URL.Path)
	}
	//call the real handler
	gv.handler.ServeHTTP(w, r)
}
