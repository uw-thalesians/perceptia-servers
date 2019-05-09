package handler

import (
	"fmt"
	"net/http"
	"time"

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

// UsersDefaultHandler handles the default routes for the users collection.
//
// If the major version in the URL is not supported, request will return an error
func (cx *Context) UsersDefaultHandler(w http.ResponseWriter, r *http.Request) {
	reqVars := mux.Vars(r)
	if ver, ok := reqVars[ReqVarMajorVersion]; ok != false && ver != "v1" {
		cx.handleVersionNotSupported(w, r, "v1", ver)
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

	// Ensure email, if supplied meets requirements
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
	}

	// Ensure Username is not in use
	_, errGUN := cx.userStore.GetByUsername(newUser.Username)
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
	}

	userINS, errINS := cx.userStore.Insert(newUser)
	if errINS != nil {
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
	}

	sessState := NewSessionState(time.Now(), *userINS)
	// This adds the authorization header to the response as well
	_, errSID := session.BeginSession(cx.sessionSigningKey, cx.sessionStore, sessState, w)
	if errSID != nil {
		retErr := &Error{
			ClientError: false,
			ServerError: true,
			Message: fmt.Sprintf("user created with username: %s; error creating new session, "+
				"please log in with your username and password manually", userINS.Username),
			Code: 0,
		}
		cx.handleErrorJson(w, r, errSID, "error beginning new session",
			retErr,
			http.StatusInternalServerError)
		return
	}
	// Send response
	_, _ = cx.respondEncode(w, userINS, http.StatusCreated)
}
