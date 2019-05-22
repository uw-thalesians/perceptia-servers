package handler

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/uw-thalesians/perceptia-servers/gateway/gateway/session"

	"github.com/uw-thalesians/perceptia-servers/gateway/gateway/user"
)

// MIDDLEWARE

// Key type for values added to an http.Request context.
type contextKey int

// authSessionActiveKey is the key used to indicate if the session is active
const authSessionActiveKey contextKey = 3030

// authSessionStateKey is the key used to retrieve the session state once added to request
const authSessionStateKey contextKey = 4040

// authUserAuthenticatedKey Key only added to request context if the user has been authenticated.
// Should be used for resources that don't need the authenticated user,
// but should still only be accessed by authorized users.
const authUserAuthenticatedKey contextKey = 5050

var ErrUserNotInContext = errors.New("authenticator: user not in context")
var ErrSessionNotInContext = errors.New("authenticator: SessionState not in context")

// EnsureAuth represents the current handler in the request/response cycle.
type EnsureAuth struct {
	handler http.Handler
	cx      *Context
}

// NewEnsureAuth constructs a new EnsureAuth struct with the provided handler.
func (cx *Context) NewEnsureAuth(handler http.Handler) http.Handler {
	return &EnsureAuth{handler, cx}
}

// ServeHTTP handles confirming the user is authenticated,
// and passing the authenticated user's profile in a new http.Request object
func (ea *EnsureAuth) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !IsUserAuthenticated(r) {
		retErr := &Error{
			ClientError: true,
			ServerError: false,
			Message:     errUnauthorized.Error(),
			Context:     fmt.Sprintf("method=%s path=%s", r.Method, r.URL.Path),
			Code:        0,
		}
		ea.cx.handleErrorJson(w, r, nil, "request to access authenticated resource, but user is not authenticated",
			retErr, http.StatusUnauthorized)
		return
	}
	ea.handler.ServeHTTP(w, r)
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
		if errGST != session.ErrNoSessionId {
			au.cx.logError(errGST, "issue getting session from request", "",
				http.StatusInternalServerError)
		}
		cxWithSessionActive := context.WithValue(r.Context(), authSessionActiveKey, false)
		cxWithUserAuthFalse := context.WithValue(cxWithSessionActive, authUserAuthenticatedKey, false)

		rWithUserAuthFalse := r.WithContext(cxWithUserAuthFalse)
		au.handler.ServeHTTP(w, rWithUserAuthFalse)
		return
	}

	//create a new request context containing the authenticated user
	cxWithSessionActive := context.WithValue(r.Context(), authSessionActiveKey, true)
	cxWithSessionState := context.WithValue(cxWithSessionActive, authSessionStateKey, sesSt)

	cxWithKey := context.WithValue(cxWithSessionState, authUserAuthenticatedKey, sesSt.Authenticated)

	//create a new request using that new context
	rWithSession := r.WithContext(cxWithKey)

	//call the real handler, passing the new request
	au.handler.ServeHTTP(w, rWithSession)
}

// GetUserFromContext returns the user stored in the request context,
// or the error ErrUserNotInContext if the authenticated user
// was not in the request context.
func GetUserFromContext(r *http.Request) (*user.User, error) {
	//Get the authenticated user.
	if authenticated, ok := r.Context().Value(authUserAuthenticatedKey).(bool); ok {
		if !authenticated {
			return nil, ErrUserNotInContext
		}
	} else {
		return nil, ErrUserNotInContext
	}
	sesSt, ok := r.Context().Value(authSessionStateKey).(*SessionState)
	if sesSt == nil || !ok || sesSt.User == nil {
		return nil, ErrUserNotInContext
	}
	return sesSt.User, nil
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
	switch val.(type) {
	case bool:
		return val.(bool)
	default:
		return false
	}
}

// IsSession will return true if the user is in a session, and false if the user is not.
func IsSession(r *http.Request) bool {
	val := r.Context().Value(authSessionActiveKey)
	switch val.(type) {
	case bool:
		return val.(bool)
	default:
		return false
	}
}
