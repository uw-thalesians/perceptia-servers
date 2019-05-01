package handler

import "net/http"

// Cors struct contains all connection scoped.
type Cors struct {
	handler http.Handler
}

// NewCors initializes and returns a new Cors struct.
func NewCors(handler http.Handler) http.Handler {
	return &Cors{handler}
}

// ServeHTTP handles adding the required CORS headers to every response passed to it.
// Will only pass response and request onto its handler when the request does not have the OPTIONS method.
func (crs *Cors) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(HeaderACAllowOrigin, ACAllowOriginAll)
	w.Header().Set(HeaderACAllowMethods, ACAllowMethods)
	w.Header().Set(HeaderACAllowHeaders, ACAllowHeaders)
	w.Header().Set(HeaderACExposeHeaders, ACExposeHeaders)
	w.Header().Set(HeaderACMaxAge, ACMaxAge)

	if r.Method != http.MethodOptions {
		crs.handler.ServeHTTP(w, r)
	}
}
