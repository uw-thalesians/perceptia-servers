package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	uuid "github.com/satori/go.uuid"

	kitlog "github.com/go-kit/kit/log"

	"github.com/uw-thalesians/perceptia-servers/gateway/gateway/session"
	"github.com/uw-thalesians/perceptia-servers/gateway/gateway/user"
)

// Context represents the shared resources amongst all http.Handler functions that receive this struct.
type Context struct {
	sessionSigningKey string
	sessionStore      session.Store
	userStore         user.Store
	logger            kitlog.Logger
}

// NewContext creates a new Context, initialized using the provided handler context values.
// Returns a pointer to the created Context.
func NewContext(sessionStore session.Store, userStore user.Store,
	sessionSigningKey string, logger kitlog.Logger) *Context {
	if sessionStore == nil || userStore == nil || len(sessionSigningKey) <= 0 {
		panic("all parameters must not be nil or empty")
	}
	return &Context{sessionSigningKey, sessionStore, userStore, logger}
}

// ensureJSONHeader is a helper method to handle checking for the application/json content-type header.
// Will return true if valid JSON header is present in the request.
func (cx *Context) ensureJSONHeader(w http.ResponseWriter, r *http.Request) bool {
	if !strings.HasPrefix(r.Header.Get(HeaderContentType), ContentTypeJSON) {
		cx.handleError(w, r, nil, "",
			fmt.Sprintf("error: %s (%s) not supported, request body must have %s of (%s)",
				HeaderContentType, r.Header.Get(HeaderContentType), HeaderContentType, ContentTypeJSON),
			http.StatusUnsupportedMediaType)
		return false
	}
	return true
}

// handleError will handle logging error and respond to client with correct message and status code.
// If len(clientErrorMessage) == 0 will only log error and will not send error to client.
// If you only need to log an error without sending error to client you should use logError instead.
func (cx *Context) handleError(w http.ResponseWriter, r *http.Request, errorToLog error, logContext,
	clientErrorMessage string,
	statusCode int) {
	logReference := cx.logError(r, errorToLog, logContext, clientErrorMessage, statusCode)

	// Only send error to client if clientErrorMessage provided.
	if len(clientErrorMessage) != 0 {
		completeClientMessage := fmt.Sprintf("error text: %s", clientErrorMessage)
		if len(logReference) != 0 {
			completeClientMessage = fmt.Sprintf("error reference: %s\n%s", logReference, completeClientMessage)
		}
		http.Error(w, completeClientMessage, statusCode)
	}
	return
}

// logError will log an error provided to it, any context including message sent to client.
// Will return a string containing the log reference to be used by the caller to associate further logging or
// response to client with this logged error.
// statusCode should be the expected status code to be sent to the client with the clientErrorMessage.
func (cx *Context) logError(r *http.Request, errorToLog error, logContext, clientErrorMessage string,
	statusCode int) string {
	logReference := uuid.NewV4().String()
	_ = cx.logger.Log("logReference", logReference,
		"context", logContext, "error", errorToLog,
		"messageToClient", clientErrorMessage,
		"httpStatusCode", statusCode)

	return logReference
}

// getUserFromContext will extract the user that was added to the request by the authenticator middleware.
// If user was not found, will respond to caller with an error and the function will return false. If false,
// calling function should return.
func (cx *Context) getUserFromContext(w http.ResponseWriter, r *http.Request) (*user.User, bool) {
	//Get the authenticated user.
	userCx, errGUC := GetUserFromContext(r)
	if errGUC != nil {
		cx.handleError(w, r, errGUC, "issue getting user from request context",
			errUnexpected.Error(), http.StatusInternalServerError)
		return nil, false
	}
	return userCx, true
}

// getSessionStateFromContext will extract the user that was added to the request by the authenticator middleware.
// If user was not found, will respond to caller with an error and the function will return false. If false,
// calling function should return.
func (cx *Context) getSessionStateFromContext(w http.ResponseWriter, r *http.Request) (*SessionState, bool) {
	//Get the authenticated user.
	sesSt, errGST := GetSessionStateFromContext(r)
	if errGST != nil {
		cx.handleError(w, r, errGST, "issue getting SessionState from request context",
			errUnexpected.Error(), http.StatusInternalServerError)
		return nil, false
	}
	return sesSt, true
}

func (cx *Context) decodeJSON(w http.ResponseWriter, r *http.Request, obj interface{},
	desc string) bool {
	if err := json.NewDecoder(r.Body).Decode(obj); err != nil {
		cx.handleError(w, r, err, fmt.Sprintf("error decoding %s from request body",
			desc), fmt.Sprintf("error extracting %s from request body",
			desc), http.StatusBadRequest)
		return false
	}
	return true
}

func (cx *Context) handleMethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	cx.handleError(w, r, nil, fmt.Sprintf("the method (%s) is not allowed", r.Method),
		errMethodNotAllowed.Error(), http.StatusMethodNotAllowed)
}

func (cx *Context) getUserFromRequest(r *http.Request) (*user.User, error) {
	//validate the session token in the request,
	//fetch the session state from the session store,
	//and return the authenticated user
	//or an error if the user is not authenticated
	sessToken, errTK := session.GetSessionID(r, cx.sessionSigningKey)
	if errTK != nil {
		return nil, errTK
	}
	authSess := &SessionState{}
	errAuth := cx.sessionStore.Get(sessToken, authSess)
	if errAuth != nil {
		return nil, errAuth
	}
	return &authSess.User, nil
}

func (cx *Context) getSessionStateFromRequest(r *http.Request) (*SessionState, error) {
	//validate the session token in the request,
	//fetch the session state from the session store,
	//and return the authenticated user
	//or an error if the user is not authenticated
	sessToken, errTK := session.GetSessionID(r, cx.sessionSigningKey)
	if errTK != nil {
		return nil, errTK
	}
	authSess := &SessionState{}
	errAuth := cx.sessionStore.Get(sessToken, authSess)
	if errAuth != nil {
		return nil, errAuth
	}
	authSess.SessionID = sessToken
	return authSess, nil
}
