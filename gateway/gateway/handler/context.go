package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/uw-thalesians/perceptia-servers/gateway/gateway/utility"

	uuid "github.com/satori/go.uuid"

	kitlog "github.com/go-kit/kit/log"

	"github.com/uw-thalesians/perceptia-servers/gateway/gateway/session"
	"github.com/uw-thalesians/perceptia-servers/gateway/gateway/user"
)

// Context represents the shared resources amongst all http.Handler functions that receive this struct.
type Context struct {
	sessionSigningKey        string
	sessionStore             session.Store
	userStore                user.Store
	logger                   kitlog.Logger
	gatewayVersion           *utility.SemVer
	gatewayVersionsSupported map[int]*utility.SemVer
	environment              string
	apiInfo                  *ApiInfo
}

// NewContext creates a new Context, initialized using the provided handler context values.
// Returns a pointer to the created Context.
func NewContext(sessionStore session.Store, userStore user.Store,
	sessionSigningKey string, gatewayVersion *utility.SemVer,
	gatewayVersionsSupported map[int]*utility.SemVer, logger kitlog.Logger,
	apiInfo *ApiInfo) *Context {
	if sessionStore == nil || userStore == nil || len(sessionSigningKey) <= 0 {
		panic("all parameters must not be nil or empty")
	}
	environment, _ := utility.DefaultEnv("GATEWAY_ENVIRONMENT", "development")
	return &Context{sessionSigningKey: sessionSigningKey, sessionStore: sessionStore,
		userStore: userStore, logger: logger, gatewayVersion: gatewayVersion,
		gatewayVersionsSupported: gatewayVersionsSupported, environment: environment, apiInfo: apiInfo}
}

type Error struct {
	Reference   string `json:"reference"`
	ServerError bool   `json:"serverError"`
	ClientError bool   `json:"clientError"`
	Message     string `json:"message"`
	Context     string `json:"context"`
	Code        int    `json:"Code"`
}

type ApiInfo struct {
	Scheme string
	Host   string
	Port   string
}

// ensureJSONHeader is a helper method to handle checking for the application/json content-type header.
// Will return true if valid JSON header is present in the request.
func (cx *Context) ensureJSONHeader(w http.ResponseWriter, r *http.Request) bool {
	if !strings.HasPrefix(r.Header.Get(HeaderContentType), ContentTypeJSON) {
		retErr := &Error{
			ClientError: true,
			ServerError: false,
			Message:     errContentTypeNotJson.Error(),
			Context:     r.Method + " path:" + r.URL.Path,
			Code:        0,
		}
		cx.handleErrorJson(w, r, nil, "request Content-Type header was not application/json",
			retErr, http.StatusUnsupportedMediaType)
		return false
	}
	return true
}

/*// handleError will handle logging error and respond to client with correct message and status code.
// If len(clientErrorMessage) == 0 will only log error and will not send error to client.
// If you only need to log an error without sending error to client you should use logError instead.
func (cx *Context) handleError(w http.ResponseWriter, r *http.Request, errorToLog error, logContext,
	clientErrorMessage string,
	statusCode int) {
	logReference := cx.logError(errorToLog, logContext, clientErrorMessage, statusCode)

	// Only send error to client if clientErrorMessage provided.
	if len(clientErrorMessage) != 0 {
		completeClientMessage := fmt.Sprintf("error text: %s", clientErrorMessage)
		if len(logReference) != 0 {
			completeClientMessage = fmt.Sprintf("error reference: %s\n%s", logReference, completeClientMessage)
		}
		http.Error(w, completeClientMessage, statusCode)
	}
	return
}*/

// handleError will handle logging error and respond to client with correct message and status code.
// If len(clientErrorMessage) == 0 will only log error and will not send error to client.
// If you only need to log an error without sending error to client you should use logError instead.
func (cx *Context) handleErrorJson(w http.ResponseWriter, r *http.Request, errorToLog error, logContext string,
	clientErrorJson *Error,
	statusCode int) {
	logReference := cx.logError(errorToLog, logContext, clientErrorJson.Message, statusCode)
	clientErrorJson.Reference = logReference
	// Only send error to client if clientErrorMessage provided.
	if clientErrorJson != nil {
		_, _ = cx.respondEncode(w, clientErrorJson, statusCode)
	}
	return
}

// logError will log an error provided to it, any context including message sent to client.
// Will return a string containing the log reference to be used by the caller to associate further logging or
// response to client with this logged error.
// statusCode should be the expected status code to be sent to the client with the clientErrorMessage.
func (cx *Context) logError(errorToLog error, logContext, clientErrorMessage string,
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
		retErr := &Error{
			ClientError: false,
			ServerError: true,
			Message:     errUnexpected.Error(),
			Context:     r.Method + " path:" + r.URL.Path,
			Code:        0,
		}
		cx.handleErrorJson(w, r, errGUC, "issue getting user from request context", retErr, http.StatusInternalServerError)
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
		retErr := &Error{
			ClientError: false,
			ServerError: true,
			Message:     errUnexpected.Error(),
			Context:     r.Method + " path:" + r.URL.Path,
			Code:        0,
		}
		cx.handleErrorJson(w, r, errGST, "issue getting session state from context", retErr, http.StatusInternalServerError)
		return nil, false
	}
	return sesSt, true
}

func (cx *Context) decodeJSON(w http.ResponseWriter, r *http.Request, obj interface{},
	desc string) bool {
	if err := json.NewDecoder(r.Body).Decode(obj); err != nil {
		retErr := &Error{
			ClientError: true,
			ServerError: false,
			Message:     errDecodingJson.Error(),
			Context:     r.Method + " path:" + r.URL.Path,
			Code:        0,
		}
		cx.handleErrorJson(w, r, err, fmt.Sprintf("error decoding %s from request body",
			desc), retErr, http.StatusBadRequest)
		return false
	}
	return true
}

func (cx *Context) handleMethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	retErr := &Error{
		ClientError: true,
		ServerError: false,
		Message:     errMethodNotAllowed.Error(),
		Context:     r.Method + " path:" + r.URL.Path,
		Code:        0,
	}
	cx.handleErrorJson(w, r, nil, fmt.Sprintf("the method (%s) is not allowed", r.Method), retErr, http.StatusMethodNotAllowed)
}

func (cx *Context) handleMajorVersionNotSupported(w http.ResponseWriter, r *http.Request, supported, requested string) {
	retErr := &Error{
		ClientError: true,
		ServerError: false,
		Message:     errMajorVersionNotSupported.Error(),
		Context:     r.Method + " path:" + r.URL.Path,
		Code:        0,
	}
	cx.handleErrorJson(w, r, nil, fmt.Sprintf("major version of API not supported; requested=%s supported=%s", requested, supported), retErr, http.StatusNotFound)
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
	return authSess.User, nil
}

func (cx *Context) getSessionUuidFromRequest(r *http.Request) (*uuid.UUID, error) {
	sessToken, errTK := session.GetSessionID(r, cx.sessionSigningKey)
	if errTK != nil {
		return nil, errTK
	}
	authSess := &SessionState{}
	errAuth := cx.sessionStore.Get(sessToken, authSess)
	if errAuth != nil {
		return nil, errAuth
	}
	return &authSess.SessionUuid, nil
}

func (cx *Context) getSessionStateFromRequest(r *http.Request) (*SessionState, error) {
	//validate the session token in the request,
	//fetch the session state from the session store,
	//and return the authenticated user
	//or an error if the user is not in a session
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

// respond allows sending a text or json object as the response, based on the item provided.
// For a text item, item should be of type string. All other types of item will be encoded as json.
// Respond will handle logging any errors that occur. respond will return an error if any errors occur,
// and a string that may contain the log reference.
func (cx *Context) respond(w http.ResponseWriter, item interface{}, statusCode int) (error, string) {
	switch item.(type) {
	case string:
		return cx.respondText(w, item.(string), statusCode)
	default:
		return cx.respondEncode(w, item, statusCode)
	}
}

// respondEncode will encode the provided object to the provided response stream.
// If an error occurs will log that error and return the error that occurred, and the logging reference string.
func (cx *Context) respondEncode(w http.ResponseWriter, objToEncode interface{}, statusCode int) (error,
	string) {
	w.Header().Set(HeaderContentType, ContentTypeJSON)
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(objToEncode)
	if err != nil {
		logReference := cx.logError(err, fmt.Sprintf("error encoding object of type: %T, object:%v", objToEncode,
			objToEncode), "", statusCode)
		return err, logReference
	}
	return nil, ""
}

func (cx *Context) respondText(w http.ResponseWriter, textToSend string, statusCode int) (error, string) {
	w.Header().Set(HeaderContentType, ContentTypeTextPlain)
	w.WriteHeader(statusCode)
	_, err := io.WriteString(w, textToSend)
	if err != nil {
		logReference := cx.logError(err, fmt.Sprintf("error writing textToSend to response stream"), "",
			statusCode)
		return err, logReference
	}
	return nil, ""
}
