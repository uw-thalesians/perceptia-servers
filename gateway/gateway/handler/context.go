package handler

import (
	"log"

	"github.com/uw-thalesians/perceptia-servers/gateway/gateway/session"
	"github.com/uw-thalesians/perceptia-servers/gateway/gateway/user"
)

// HandlerContext represents the shared resources amongst all http.Handler functions that receive this struct.
type HandlerContext struct {
	signingKey   string
	sessionStore session.Store
	userStore    user.Store
	logger       log.Logger
}
