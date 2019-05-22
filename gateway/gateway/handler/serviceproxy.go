package handler

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"os"
)

// NewServiceProxy is an http proxy that forwards requests on to the appropriate microservice,
// as noted by the hostname provided.
func (cx *Context) NewServiceProxy(hostname, port string) *httputil.ReverseProxy {
	if len(hostname) == 0 {
		_ = cx.logger.Log("error", "there must be at least one microservice address provided", "result", "exit")
		os.Exit(1)
	}

	hostnameAndPort := fmt.Sprintf("%s:%s", hostname, port)

	return &httputil.ReverseProxy{
		Director: func(r *http.Request) {
			r.URL.Scheme = "http"
			r.URL.Host = hostnameAndPort
			// Remove existing User Uuid header
			r.Header.Del(HeaderPerceptiaUserUuid)
			r.Header.Del(HeaderPerceptiaSessionUuid)

			if user, errGAU := cx.getUserFromRequest(r); errGAU == nil && user != nil {
				r.Header.Set(HeaderPerceptiaUserUuid, user.Uuid.String())
			}
			if sesUuid, errGAU := cx.getSessionUuidFromRequest(r); errGAU == nil && sesUuid != nil {
				r.Header.Set(HeaderPerceptiaSessionUuid, sesUuid.String())
			}
		},
	}
}
