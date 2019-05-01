package handler

import (
	"context"
	"errors"
	"net/http"

	"github.com/uw-thalesians/perceptia-servers/gateway/gateway/user"
)

//TODO: Decouple authenticating user/adding user to request and rejecting request.
// Create new middleware that specifically rejects a request if the user is not authenticated.
// That way authenticator can be wrapped around all structs,
// but the ensure authenticated can be wrapped around only the handlers that must only be accessed by authenticated
// user.

// MIDDLEWARE

// Key type for values added to an http.Request context.
type contextKey int

// authUserKey Key for the authenticated user stored in the http.Request context.
// TODO: tobe removed, replaced with session state key. Remove uses first!
const authUserKey contextKey = 0

const authSessionStateKey contextKey = 1

// authUserAuthenticatedKey Key only added to request context if the user has been authenticated.
// Should be used for resources that don't need the authenticated user,
// but should still only be accessed by authorized users.
const authUserAuthenticatedKey contextKey = 2

var ErrUserNotInContext = errors.New("authenticator: user not in context")
var ErrSessionNotInContext = errors.New("authenticator: SessionState not in context")

// EnsureAuth represents the current handler in the request/response cycle.
type EnsureAuth struct {
	handlerFunc http.HandlerFunc
	cx          *Context
}

// NewEnsureAuth constructs a new EnsureAuth struct with the provided handler.
func (cx *Context) NewEnsureAuth(handlerFunc http.HandlerFunc) *EnsureAuth {
	return &EnsureAuth{handlerFunc, cx}
}

// ServeHTTP handles confirming the user is authenticated,
// and passing the authenticated user's profile in a new http.Request object
func (ea *EnsureAuth) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !IsUserAuthenticated(r) {
		ea.cx.handleError(w, r, nil, "request to access authenticated resource, but user is not authenticated",
			"access not authorized, please sign-in", http.StatusUnauthorized)
		return
	}
	ea.handlerFunc.ServeHTTP(w, r)
}

// Authenticator represents the current handler in the request/response cycle.
type Authenticator struct {
	handler http.Handler
	cx      *Context
}

// NewAuthenticator constructs a new Authenticator struct with the provided handler and Context.
func (cx *Context) NewAuthenticator(handler http.Handler) http.Handler {
	return &Authenticator{handler, cx}
}

// ServeHTTP ,
// and passing the authenticated user's profile in a new http.Request object
func (au *Authenticator) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	sesSt, errGST := au.cx.getSessionStateFromRequest(r)
	if errGST != nil {
		au.handler.ServeHTTP(w, r)
		return
	}

	//create a new request context containing the authenticated user
	cxWithUser := context.WithValue(r.Context(), authSessionStateKey, sesSt)
	cxWithUserAuth := context.WithValue(cxWithUser, authUserAuthenticatedKey, true)
	//create a new request using that new context
	rWithUser := r.WithContext(cxWithUserAuth)

	//call the real handler, passing the new request
	au.handler.ServeHTTP(w, rWithUser)
}

// GetUserFromContext returns the user stored in the request context,
// or the error ErrUserNotInContext if the authenticated user
// was not in the request context.
func GetUserFromContext(r *http.Request) (*user.User, error) {
	//Get the authenticated user.
	sesSt, ok := r.Context().Value(authSessionStateKey).(*SessionState)
	if sesSt == nil || !ok || &sesSt.User == nil {
		return nil, ErrUserNotInContext
	}
	return &sesSt.User, nil
}

// GetSessionStateFromContext returns the SessionState stored in the request context,
// or the error ErrSessionNotInContext if the SessionState
// was not in the request context.
func GetSessionStateFromContext(r *http.Request) (*SessionState, error) {
	//Get the authenticated user.
	sesSt, ok := r.Context().Value(authSessionStateKey).(*SessionState)
	if sesSt == nil || !ok {
		return nil, ErrSessionNotInContext
	}
	return sesSt, nil
}

// IsUserAuthenticated will return true if the user is authenticated, and false if the user is not.
func IsUserAuthenticated(r *http.Request) bool {
	val := r.Context().Value(authUserAuthenticatedKey)
	if val == nil {
		return false
	}
	return true
}
