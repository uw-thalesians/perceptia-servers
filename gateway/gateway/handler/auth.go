package handler

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	uuid "github.com/satori/go.uuid"

	"github.com/uw-thalesians/perceptia-servers/gateway/gateway/session"

	"github.com/uw-thalesians/perceptia-servers/gateway/gateway/user"

	"github.com/gorilla/mux"
)

// Expected json format to be provided by client
type newUserJson struct {
	Username    string `json:"username"`
	FullName    string `json:"fullName"`
	DisplayName string `json:"displayName"`
	Password    string `json:"password"`
	Email       string `json:"email,omitempty"`
}

type signInCredentialsJson struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// UsersDefaultHandler handles the default routes for the users collection.
//
// If the major version in the URL is not supported, request will return an error
func (cx *Context) UsersDefaultHandler(w http.ResponseWriter, r *http.Request) {
	reqVars := mux.Vars(r)
	if ver, ok := reqVars[ReqVarMajorVersion]; ok == true {
		if ver != "v1" {
			cx.handleVersionNotSupported(w, r, "v1", ver)
			return
		}
	} else {
		retErr := &Error{
			ClientError: false,
			ServerError: true,
			Message:     errUnexpected.Error(),
			Code:        0,
		}
		cx.handleErrorJson(w, r, nil, "did not find major version key in request vars", retErr, http.StatusInternalServerError)
		return
	}

	switch r.Method {
	case http.MethodPost:
		cx.usersHandlerV1Post(w, r)
		return
	default:
		cx.handleMethodNotAllowed(w, r)
		return
	}
}

// UsersDefaultHandler handles the authenticated routes for the users collection.
//
// If the major version in the URL is not supported, request will return an error
func (cx *Context) UsersSpecificHandler(w http.ResponseWriter, r *http.Request) {
	reqVars := mux.Vars(r)
	if ver, ok := reqVars[ReqVarMajorVersion]; ok == true {
		if ver != "v1" {
			cx.handleVersionNotSupported(w, r, "v1", ver)
			return
		}
	} else {
		retErr := &Error{
			ClientError: false,
			ServerError: true,
			Message:     errUnexpected.Error(),
			Code:        0,
		}
		cx.handleErrorJson(w, r, nil, "did not find major version key in request vars", retErr, http.StatusInternalServerError)
		return
	}

	//Get the authenticated user.
	userCx, ok := cx.getUserFromContext(w, r)
	if !ok {
		// Ends method execution if user was not found in the request context.
		return
	}

	switch r.Method {
	case http.MethodGet:
		cx.usersSpecificHandlerV1Get(w, r, userCx)
		return
	case http.MethodDelete:
		cx.usersSpecificHandlerV1Delete(w, r, userCx)

	default:
		cx.handleMethodNotAllowed(w, r)
		return
	}
}

// SessionsDefaultHandler handles the default routes for the sessions collection.
//
// If the major version in the URL is not supported, request will return an error
func (cx *Context) SessionsDefaultHandler(w http.ResponseWriter, r *http.Request) {
	reqVars := mux.Vars(r)
	if ver, ok := reqVars[ReqVarMajorVersion]; ok == true && ver != "v1" {
		cx.handleVersionNotSupported(w, r, "v1", ver)
		return
	}
	switch r.Method {
	case http.MethodPost:
		if strings.HasPrefix(r.Header.Get(HeaderContentType), ContentTypeJSON) {
			cx.sessionsHandlerV1PostAuthenticated(w, r)
		} else {
			cx.sessionsHandlerV1PostUnauthenticated(w, r)
		}

		return
	default:
		cx.handleMethodNotAllowed(w, r)
		return
	}
}

// SessionsSpecificHandler handles specific routes for the sessions collection.
//
// If the major version in the URL is not supported, request will return an error
func (cx *Context) SessionsSpecificHandler(w http.ResponseWriter, r *http.Request) {
	reqVars := mux.Vars(r)
	if ver, ok := reqVars[ReqVarMajorVersion]; ok != false && ver != "v1" {
		cx.handleVersionNotSupported(w, r, "v1", ver)
		return
	}

	switch r.Method {
	case http.MethodDelete:
		cx.sessionsSpecificHandlerV1Delete(w, r)

	default:
		cx.handleMethodNotAllowed(w, r)
		return
	}
}

// usersHandlerPost is a helper method for UsersDefaultHandler to handle Post requests to the users collection.
func (cx *Context) usersHandlerV1Post(w http.ResponseWriter, r *http.Request) {
	// Test if json header is present, if not, return
	if !cx.ensureJSONHeader(w, r) {
		return
	}

	newUser := &user.NewUser{}
	newUserFromClient := &newUserJson{}

	if !cx.decodeJSON(w, r, newUserFromClient, "NewUserJson") {
		// return if unable to decode provided json object
		return
	}
	newUser.DisplayName = newUserFromClient.DisplayName
	newUser.Username = newUserFromClient.Username
	newUser.FullName = newUserFromClient.FullName

	// Clean up user supplied names
	newUser.PrepNewUser()

	// Ensure password meets requirements
	if err := user.ValidatePassword(newUserFromClient.Password); err != nil {
		retErr := &Error{
			ClientError: true,
			ServerError: false,
			Message:     fmt.Sprintf("the provided password is not a valid password: %s", err.Error()),
			Context:     "POST path:" + r.URL.Path,
			Code:        0,
		}
		cx.handleErrorJson(w, r, err, "error: the provided password is not a valid password",
			retErr, http.StatusBadRequest)
		return
	}

	passHash, errCEH := user.CreateEncodedHash(newUserFromClient.Password)
	if errCEH != nil {
		retErr := &Error{
			ClientError: false,
			ServerError: true,
			Message:     errUnexpected.Error(),
			Context:     "POST path:" + r.URL.Path,
			Code:        0,
		}
		cx.handleErrorJson(w, r, errCEH, "error: unable to create hash of provided password",
			retErr, http.StatusInternalServerError)
		return
	}
	newUser.EncodedHash = passHash

	// Ensure user supplied values meet requirements
	if err := newUser.ValidateNewUser(); err != nil {
		retErr := &Error{
			ClientError: true,
			ServerError: false,
			Message:     fmt.Sprintf("the provided new user is not a valid user: %s", err.Error()),
			Context:     "POST path:" + r.URL.Path,
			Code:        0,
		}
		cx.handleErrorJson(w, r, err, "error: the provided NewUser is not a valid user",
			retErr, http.StatusBadRequest)
		return
	}

	/*	// Ensure email, if supplied meets requirements
		userEmail := ""
		if len(newUserFromClient.Email) != 0 {
			var errCE error
			userEmail, errCE = user.CleanEmail(newUserFromClient.Email)
			if errCE != nil {
				retErr := &Error{
					ClientError: true,
					ServerError: false,
					Message:     fmt.Sprintf("error: the provided email is not a valid email: %s", errCE.Error()),
					Context:     "POST path:" + r.URL.Path,
					Code:        0,
				}
				cx.handleErrorJson(w, r, errCE, "error: the provided email is not a valid email",
					retErr, http.StatusBadRequest)
				return
			}
		}*/

	// Ensure Username is not in use
	_, errGUN := cx.userStore.ReadUserUuid(newUser.Username)
	if errGUN == nil {
		retErr := &Error{
			ClientError: true,
			ServerError: false,
			Message:     errAccountUserNameUnavailable.Error(),
			Context:     "POST path:" + r.URL.Path,
			Code:        0,
		}
		cx.handleErrorJson(w, r, nil, fmt.Sprintf("user with that username already exists: %s", newUser.Username),
			retErr,
			http.StatusConflict)
		return
	} else if errGUN != user.ErrUserNotFound {
		retErr := &Error{
			ClientError: true,
			ServerError: false,
			Message:     errUnexpected.Error(),
			Context:     "POST path:" + r.URL.Path,
			Code:        0,
		}
		cx.handleErrorJson(w, r, nil, "error occurred trying to search for user",
			retErr,
			http.StatusInternalServerError)
		return
	}

	userINS, errINS := cx.userStore.CreateUser(newUser)
	if errINS != nil {
		if errINS == user.ErrUserAlreadyExists {
			retErr := &Error{
				ClientError: false,
				ServerError: true,
				Message:     errUnexpected.Error(),
				Context:     "POST path:" + r.URL.Path,
				Code:        0,
			}
			cx.handleErrorJson(w, r, errINS, "uuid matched existing user uuid", retErr,
				http.StatusInternalServerError)
			return
		} else if errINS == user.ErrUsernameUnavailable {
			retErr := &Error{
				ClientError: true,
				ServerError: false,
				Message:     errAccountUserNameUnavailable.Error(),
				Context:     "POST path:" + r.URL.Path,
				Code:        0,
			}
			cx.handleErrorJson(w, r, nil, fmt.Sprintf("user with that username already exists: %s", newUser.Username),
				retErr,
				http.StatusConflict)
			return
		}
		retErr := &Error{
			ClientError: false,
			ServerError: true,
			Message:     errUnexpected.Error(),
			Context:     "POST path:" + r.URL.Path,
			Code:        0,
		}
		cx.handleErrorJson(w, r, errINS, "error adding user to database", retErr,
			http.StatusInternalServerError)
		return
	}

	/*
	   TODO: Add email to account
	   if len(userEmail) > 0 {
	   		errAE := cx.userStore.InsertEmail(userINS.Uuid, userEmail)
	   		if errAE != nil {
	   			retErr := &Error{
	   				ClientError: false,
	   				ServerError: true,
	   				Message: fmt.Sprintf("user created with username: %s; error adding users email, "+
	   					"please log in with your username and password manually", userINS.Username),
	   				Context: "POST path:" + r.URL.Path,
	   				Code:    0,
	   			}
	   			cx.handleErrorJson(w, r, errINS, "error adding users email to database",
	   				retErr,
	   				http.StatusInternalServerError)
	   			return
	   		}
	   	}*/

	sesId, sesUuid, errSID := session.CreateSession(cx.sessionSigningKey)
	if errSID != nil {
		retErr := &Error{
			ClientError: false,
			ServerError: true,
			Message: fmt.Sprintf("user created with username: %s; error creating new session, "+
				"please log in with your username and password manually", userINS.Username),
			Context: "POST path:" + r.URL.Path,
			Code:    0,
		}
		cx.handleErrorJson(w, r, errSID, "error beginning new session",
			retErr,
			http.StatusInternalServerError)
		return
	}
	sessState := NewSessionState(time.Now(), userINS, sesUuid, sesId, true)
	// This adds the authorization header to the response as well
	errBS := session.BeginSession(sesId, cx.sessionStore, sessState, w)
	if errBS != nil {
		retErr := &Error{
			ClientError: false,
			ServerError: true,
			Message: fmt.Sprintf("user created with username: %s; error creating new session, "+
				"please log in with your username and password manually", userINS.Username),
			Context: "POST path:" + r.URL.Path,
			Code:    0,
		}
		cx.handleErrorJson(w, r, errBS, "error beginning new session",
			retErr,
			http.StatusInternalServerError)
		return
	}
	// Send response
	_, _ = cx.respondEncode(w, userINS, http.StatusCreated)
}

// usersSpecificHandlerV1Get is a helper method for UsersSpecificHandler to handle Get requests to the users collection.
func (cx *Context) usersSpecificHandlerV1Get(w http.ResponseWriter, r *http.Request, userCx *user.User) {
	reqVars := mux.Vars(r)
	reqUserUuidString, ok := reqVars[ReqVarUserUuid]
	if !ok {
		retErr := &Error{
			ClientError: false,
			ServerError: true,
			Message:     fmt.Sprintf("uuid not extracted from request"),
			Context:     "GET path:" + r.URL.Path,
			Code:        0,
		}
		cx.handleErrorJson(w, r, nil, "uuid expected in path, but not found in mux vars", retErr, http.StatusInternalServerError)
		return
	}

	reqUserUuid, errUFS := uuid.FromString(reqUserUuidString)
	if errUFS != nil {
		retErr := &Error{
			ClientError: false,
			ServerError: true,
			Message:     fmt.Sprintf("unable to get valid uuid from path"),
			Context:     "GET path:" + r.URL.Path,
			Code:        0,
		}
		cx.handleErrorJson(w, r, errUFS, "issue converting string to valid uuid", retErr, http.StatusInternalServerError)
		return
	}
	if userCx != nil && !uuid.Equal(userCx.Uuid, reqUserUuid) {
		retErr := &Error{
			ClientError: true,
			ServerError: false,
			Message:     errActionNotAuthorized.Error(),
			Context:     "GET path:" + r.URL.Path,
			Code:        0,
		}
		cx.handleErrorJson(w, r, errUFS, "user attempted to get info about another user", retErr, http.StatusUnauthorized)
		return
	}

	userProfile, errGID := cx.userStore.ReadUserInfo(reqUserUuid)
	if errGID != nil {
		if errGID == user.ErrUserNotFound {
			retErr := &Error{
				ClientError: true,
				ServerError: false,
				Message:     errUserNotFound.Error(),
				Context:     "GET path:" + r.URL.Path,
				Code:        0,
			}
			cx.handleErrorJson(w, r, errGID, fmt.Sprintf("requested user not found in database: uuid=%s", reqUserUuid.String()), retErr, http.StatusNotFound)
		} else {
			retErr := &Error{
				ClientError: false,
				ServerError: true,
				Message:     errUnexpected.Error(),
				Context:     "GET path:" + r.URL.Path,
				Code:        0,
			}
			cx.handleErrorJson(w, r, errGID, fmt.Sprintf("issue retrieving user from database: uuid=%s", reqUserUuid.String()), retErr, http.StatusInternalServerError)
		}
		return
	}
	// Send response
	_, _ = cx.respondEncode(w, userProfile, http.StatusOK)
}

/*// usersSpecificHandlerV1Patch is a helper method for SpecificUserHandler to handle Patch requests to the users collection.
func (cx *Context) usersSpecificHandlerV1Patch(w http.ResponseWriter, r *http.Request, userCx *user.User) {
	reqVars := mux.Vars(r)
	reqUserUuidString, ok := reqVars[ReqVarUserUuid]
	if !ok {
		retErr := &Error{
			ClientError: false,
			ServerError: true,
			Message:     fmt.Sprintf("uuid not extracted from request"),
			Code:        0,
		}
		cx.handleErrorJson(w, r, nil, "uuid expected in path, but not found in mux vars", retErr, http.StatusInternalServerError)
		return
	}

	reqUserUuid, errUFS := uuid.FromString(reqUserUuidString)
	if errUFS != nil {
		retErr := &Error{
			ClientError: false,
			ServerError: true,
			Message:     fmt.Sprintf("unable to get valid uuid from path"),
			Code:        0,
		}
		cx.handleErrorJson(w, r, errUFS, "issue converting string to valid uuid", retErr, http.StatusInternalServerError)
		return
	}

	if !uuid.Equal(reqUserUuid, userCx.Uuid) {
		retErr := &Error{
			ClientError: true,
			ServerError: false,
			Message:     errActionNotAuthorized.Error(),
			Code:        0,
		}
		cx.handleErrorJson(w, r, nil,
			fmt.Sprintf("logged in user tried to update a different users profile: user=%s userToUpdate=%s",
				userCx.Uuid.String(), reqUserUuid.String()), retErr, http.StatusForbidden)
		return
	}

	// Test if json header is present, if not, return
	if !cx.ensureJSONHeader(w, r) {
		return
	}

	updatesClient := &userUpdatesJson{}
	if !cx.decodeJSON(w, r, updatesClient, "userUpdatesJson") {
		// return if unable to decode updates
		return
	}
	updates := &user.Updates{}
	updates.FullName = updatesClient.FullName
	updates.DisplayName = updatesClient.DisplayName
	updates.PrepUpdates()
	if err := updates.ValidateUpdates(); err != nil {
		retErr := &Error{
			ClientError: true,
			ServerError: false,
			Message:     err.Error(),
			Code:        0,
		}
		cx.handleErrorJson(w, r, err,
			"provided updates invalid", retErr, http.StatusBadRequest)
	}
	var updatedUser *user.User
	if len(updates.FullName) > 0 {
		var err error
		updatedUser, err = cx.userStore.UpdateFullName(userCx.Uuid, updates.FullName)
		if err != nil {
			retErr := &Error{
				ClientError: false,
				ServerError: true,
				Message:     err.Error(),
				Code:        0,
			}
			cx.handleErrorJson(w, r, err,
				"error occurred when updating fullName", retErr, http.StatusInternalServerError)
			return
		}
	}
	if len(updates.DisplayName) > 0 {
		var err error
		updatedUser, err = cx.userStore.UpdateDisplayName(userCx.Uuid, updates.DisplayName)
		if err != nil {
			retErr := &Error{
				ClientError: false,
				ServerError: true,
				Message:     err.Error(),
				Code:        0,
			}
			cx.handleErrorJson(w, r, err,
				"error occurred when updating displayName", retErr, http.StatusInternalServerError)
			return
		}
	}

	sesSt, errGSR := cx.getSessionStateFromRequest(r)
	if errGSR != nil && sesSt != nil {
		_ = cx.sessionStore.Save(sesSt.SessionID, SessionState{
			User:      *updatedUser,
			SessionID: sesSt.SessionID,
			StartTime: sesSt.StartTime,
		})
	}
	// Send response to client.
	_, _ = cx.respondEncode(w, updatedUser, http.StatusOK)
}*/

// usersSpecificHandlerV1Delete is a helper method for SpecificUserHandler to handle Delete requests to the users collection.
func (cx *Context) usersSpecificHandlerV1Delete(w http.ResponseWriter, r *http.Request, userCx *user.User) {
	reqVars := mux.Vars(r)
	reqUserUuidString, ok := reqVars[ReqVarUserUuid]
	if !ok {
		retErr := &Error{
			ClientError: false,
			ServerError: true,
			Message:     fmt.Sprintf("uuid not extracted from request"),
			Context:     "DELETE path:" + r.URL.Path,
			Code:        0,
		}
		cx.handleErrorJson(w, r, nil, "uuid expected in path, but not found in mux vars", retErr, http.StatusInternalServerError)
		return
	}

	reqUserUuid, errUFS := uuid.FromString(reqUserUuidString)
	if errUFS != nil {
		retErr := &Error{
			ClientError: false,
			ServerError: true,
			Message:     fmt.Sprintf("unable to get valid uuid from path"),
			Context:     "DELETE path:" + r.URL.Path,
			Code:        0,
		}
		cx.handleErrorJson(w, r, errUFS, "issue converting string to valid uuid", retErr, http.StatusInternalServerError)
		return
	}

	if userCx != nil && !uuid.Equal(reqUserUuid, userCx.Uuid) {
		retErr := &Error{
			ClientError: true,
			ServerError: false,
			Message:     errActionNotAuthorized.Error(),
			Context:     "DELETE path:" + r.URL.Path,
			Code:        0,
		}
		cx.handleErrorJson(w, r, nil,
			fmt.Sprintf("logged in user tried to update a different users profile: user=%s userToUpdate=%s",
				userCx.Uuid.String(), reqUserUuid.String()), retErr, http.StatusForbidden)
		return
	}

	errDU := cx.userStore.DeleteUser(reqUserUuid)
	if errDU != nil {
		retErr := &Error{
			ClientError: false,
			ServerError: true,
			Message:     errUnexpected.Error(),
			Context:     "DELETE path:" + r.URL.Path,
			Code:        0,
		}
		cx.handleErrorJson(w, r, errDU,
			"error occurred while attempting to delete user account", retErr, http.StatusInternalServerError)
		return
	}
	sesSt, errGSR := cx.getSessionStateFromRequest(r)
	if errGSR == nil && sesSt != nil {
		_ = cx.sessionStore.Delete(sesSt.SessionID)
	}
	// Send response to client.
	_, _ = cx.respond(w, "account deleted successfully", http.StatusOK)
}

// sessionsHandlerV1Post is a helper method for SessionsHandler to handle Post requests to the sessions collection.
func (cx *Context) sessionsHandlerV1PostAuthenticated(w http.ResponseWriter, r *http.Request) {

	signInCredentials := &signInCredentialsJson{}
	if !cx.decodeJSON(w, r, signInCredentials, "signInCredentials") {
		// return if unable to decode credentials
		return
	}
	credentials := &user.SignInCredentials{}
	credentials.Username = signInCredentials.Username
	credentials.Password = signInCredentials.Password
	errVC := credentials.ValidateSignInCredentials()
	if errVC != nil {
		retErr := &Error{
			ClientError: true,
			ServerError: false,
			Message:     errInvalidCredentials.Error(),
			Context:     "POST path:" + r.URL.Path,
			Code:        0,
		}
		cx.handleErrorJson(w, r, errVC,
			"provided credentials are not valid", retErr, http.StatusBadRequest)
		return
	}
	validUserHash, errGEH := cx.userStore.ReadUserEncodedHash(credentials.Username)
	// If user does not exist, will attempt to compare provided credentials
	// against a fictitious valid user.
	if errGEH != nil {
		retErr := &Error{
			ClientError: true,
			ServerError: false,
			Message:     errInvalidCredentials.Error(),
			Context:     "POST path:" + r.URL.Path,
			Code:        0,
		}
		cx.handleErrorJson(w, r, errGEH,
			"provided credentials are not valid", retErr, http.StatusUnauthorized)
		return
	}
	valid, errAuth := user.Authenticate(credentials.Password, validUserHash)
	if errAuth != nil {
		retErr := &Error{
			ClientError: true,
			ServerError: false,
			Message:     errUnexpected.Error(),
			Context:     "POST path:" + r.URL.Path,
			Code:        0,
		}
		cx.handleErrorJson(w, r, errAuth,
			"error occurred when trying to validate credentials", retErr, http.StatusInternalServerError)
		return
	}
	if !valid {
		retErr := &Error{
			ClientError: true,
			ServerError: false,
			Message:     errInvalidCredentials.Error(),
			Context:     "POST path:" + r.URL.Path,
			Code:        0,
		}
		cx.handleErrorJson(w, r, errAuth,
			"provided credentials are not the same as were used to create the hash for this user",
			retErr, http.StatusUnauthorized)
		return
	}
	userUuid, errRU := cx.userStore.ReadUserUuid(credentials.Username)
	if errRU != nil {
		retErr := &Error{
			ClientError: true,
			ServerError: false,
			Message:     errUserNotFound.Error(),
			Code:        0,
		}
		cx.handleErrorJson(w, r, errRU,
			"user was not found in database but should be in database", retErr, http.StatusBadRequest)
	}
	userPro, errGUUN := cx.userStore.ReadUserInfo(*userUuid)
	if errGUUN != nil {
		retErr := &Error{
			ClientError: true,
			ServerError: false,
			Message:     errUserNotFound.Error(),
			Code:        0,
		}
		cx.handleErrorJson(w, r, errGUUN,
			"user was not found in database but should be in database", retErr, http.StatusBadRequest)
		return
	}
	// Begin new session
	var sesId session.SessionID
	var sesUuid uuid.UUID
	var errSID error
	if IsSession(r) {
		sesSt, errGSC := GetSessionStateFromContext(r)
		if errGSC != nil || sesSt == nil {
			retErr := &Error{
				ClientError: false,
				ServerError: true,
				Message:     "error creating new session, please try again",
				Context:     "POST path:" + r.URL.Path,
				Code:        0,
			}
			cx.handleErrorJson(w, r, errGSC, "error getting session state from context",
				retErr,
				http.StatusInternalServerError)
			return
		}
		if !IsUserAuthenticated(r) {
			sesId = sesSt.SessionID
			sesUuid = sesSt.SessionUuid
		} else if sesSt.User != nil && uuid.Equal(sesSt.User.Uuid, *userUuid) {
			sesId = sesSt.SessionID
			sesUuid = sesSt.SessionUuid
		} else {
			sesId, sesUuid, errSID = session.CreateSession(cx.sessionSigningKey)
		}

	} else {
		sesId, sesUuid, errSID = session.CreateSession(cx.sessionSigningKey)
	}

	if errSID != nil {
		retErr := &Error{
			ClientError: false,
			ServerError: true,
			Message:     "error creating new session, please try again",
			Context:     "POST path:" + r.URL.Path,
			Code:        0,
		}
		cx.handleErrorJson(w, r, errSID, "error beginning new session",
			retErr,
			http.StatusInternalServerError)
		return
	}
	sessState := NewSessionState(time.Now(), userPro, sesUuid, sesId, true)
	// This adds the authorization header to the response as well
	errBS := session.BeginSession(sesId, cx.sessionStore, sessState, w)
	if errBS != nil {
		retErr := &Error{
			ClientError: false,
			ServerError: true,
			Message:     "error creating new session, please try again",
			Context:     "POST path:" + r.URL.Path,
			Code:        0,
		}
		cx.handleErrorJson(w, r, errBS, "error beginning new session",
			retErr,
			http.StatusInternalServerError)
		return
	}
	// Send response
	_, _ = cx.respondEncode(w, userPro, http.StatusCreated)
}

// sessionsHandlerV1PostUnauthenticated is a helper method for SessionsHandler to handle Post requests to the sessions
// collection from a non authenticating user.
func (cx *Context) sessionsHandlerV1PostUnauthenticated(w http.ResponseWriter, r *http.Request) {

	// Begin new session
	var sesId session.SessionID
	var sesUuid uuid.UUID
	var errSID error

	sesId, sesUuid, errSID = session.CreateSession(cx.sessionSigningKey)

	if errSID != nil {
		retErr := &Error{
			ClientError: false,
			ServerError: true,
			Message:     "error creating new session, please try again",
			Context:     "POST path:" + r.URL.Path,
			Code:        0,
		}
		cx.handleErrorJson(w, r, errSID, "error beginning new session",
			retErr,
			http.StatusInternalServerError)
		return
	}
	sessState := NewSessionState(time.Now(), nil, sesUuid, sesId, false)
	// This adds the authorization header to the response as well
	errBS := session.BeginSession(sesId, cx.sessionStore, sessState, w)
	if errBS != nil {
		retErr := &Error{
			ClientError: false,
			ServerError: true,
			Message:     "error creating new session, please try again",
			Context:     "POST path:" + r.URL.Path,
			Code:        0,
		}
		cx.handleErrorJson(w, r, errBS, "error beginning new session",
			retErr,
			http.StatusInternalServerError)
		return
	}
	// Send response
	w.WriteHeader(http.StatusCreated)
	return
}

// sessionsSpecificHandlerV1Delete is a helper method for SpecificSessionHandler to handle Delete requests to the
// sessions collection.
func (cx *Context) sessionsSpecificHandlerV1Delete(w http.ResponseWriter, r *http.Request) {
	reqVars := mux.Vars(r)
	sesVar, ok := reqVars[ReqVarSession]
	if !ok {
		retErr := &Error{
			ClientError: false,
			ServerError: true,
			Message:     fmt.Sprintf("session identifier not extracted from path"),
			Code:        0,
		}
		cx.handleErrorJson(w, r, nil, "session identifier expected in path, but not found in mux vars", retErr, http.StatusInternalServerError)
		return
	}

	if sesVar != SpecificSessionHandlerDeleteUserAlias {
		retErr := &Error{
			ClientError: true,
			ServerError: false,
			Message:     fmt.Sprintf("session identifier provided not valid"),
			Code:        0,
		}
		cx.handleErrorJson(w, r, nil, "session identifier provided but not a valid identifier", retErr, http.StatusBadRequest)
		return
	}

	sessID, errGSID := session.GetSessionID(r, cx.sessionSigningKey)
	if errGSID != nil {
		retErr := &Error{
			ClientError: false,
			ServerError: true,
			Message:     errUnexpected.Error(),
			Code:        0,
		}
		cx.handleErrorJson(w, r, errGSID, "unable to extract sessionID from request", retErr, http.StatusInternalServerError)
		return
	}
	if ok, err := cx.sessionStore.Exists(sessID); err != nil || ok == false {
		retErr := &Error{
			ClientError: true,
			ServerError: false,
			Message:     errSessionNotFound.Error(),
			Code:        0,
		}
		cx.handleErrorJson(w, r, errGSID, "session key not in session database", retErr, http.StatusBadRequest)
		return
	}
	errDSID := cx.sessionStore.Delete(sessID)
	if errDSID != nil {
		retErr := &Error{
			ClientError: false,
			ServerError: true,
			Message:     errUnexpected.Error(),
			Code:        0,
		}
		cx.handleErrorJson(w, r, errDSID, "unable to delete user session", retErr, http.StatusInternalServerError)
		return
	}

	// Send response
	_, _ = cx.respondText(w, "session successfully ended", http.StatusOK)
}
