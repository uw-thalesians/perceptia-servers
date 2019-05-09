package handler

import (
	"fmt"
	"net/http"
	"time"

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

	// Ensure user supplied values meet requirements
	if err := newUser.ValidateNewUser(); err != nil {
		cx.handleError(w, r, err, "error: the provided NewUser is not a valid user",
			fmt.Sprintf("error: the provided new user is not a valid user: %s", err.Error()), http.StatusBadRequest)
		return
	}

	// Ensure password meets requirements
	if err := user.ValidatePassword(newUserFromClient.Password); err != nil {
		cx.handleError(w, r, err, "error: the provided password is not a valid password",
			fmt.Sprintf("error: the provided password is not a valid password: %s", err.Error()), http.StatusBadRequest)
		return
	}
	// Ensure email, if supplied meets requirements
	userEmail := ""
	if len(newUserFromClient.Email) != 0 {
		userEmail, errCE := user.CleanEmail(newUserFromClient.Email)
		if errCE != nil {
			cx.handleError(w, r, errCE, "error: the provided email is not a valid email",
				fmt.Sprintf("error: the provided email is not a valid email: %s", errCE.Error()), http.StatusBadRequest)
			return
		}
	}

	// Ensure UserName is not in use
	_, errGUN := cx.userStore.GetByUsername(newUser.Username)
	if errGUN == nil {
		cx.handleError(w, r, nil, fmt.Sprintf("user with that username already exists: %s", newUser.Username),
			errAccountUserNameUnavailable.Error(),
			http.StatusConflict)
		return
	}

	passHash := user.CreateEncodedHash(newUserFromClient.Password)
	// Ensure password meets requirements
	if err := user.ValidatePassword(newUserFromClient.Password); err != nil {
		cx.handleError(w, r, err, "error: the provided password is not a valid password",
			fmt.Sprintf("error: the provided password is not a valid password: %s", err.Error()), http.StatusBadRequest)
		return
	}

	userINS, errINS := cx.userStore.Insert()
	if errINS != nil {
		cx.handleError(w, r, errINS, "error adding user to database", errUnexpected,
			http.StatusInternalServerError)
		return
	}
	// Load new user into the trie
	cx.trie.LoadTrieFromUser(*userINS)

	sessState := NewSessionState(time.Now(), *userINS)
	// This adds the authorization header to the response as well
	_, errSID := sessions.BeginSession(cx.signingKey, cx.sessionStore, sessState, w)
	if errSID != nil {
		cx.handleError(w, r, errSID, "error beginning new session",
			fmt.Sprintf("user created with id: %s; error creating new session, "+
				"please log in with your email and password manually", string(userINS.ID)),
			http.StatusInternalServerError)
		return
	}
	// Send response
	cx.respondEncode(w, userINS, http.StatusCreated)
}
