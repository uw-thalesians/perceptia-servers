// Package session provides tools for managing user sessions.
package session

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

const HeaderAuthorization = "Authorization"
const ParamAuthorization = "access_token"
const SchemeBearer = "Bearer "

// ErrNoSessionID is used when no session ID was found in the Authorization header.
var ErrNoSessionID = errors.New("session: no session ID found in " + HeaderAuthorization + " header")

// ErrInvalidScheme is used when the authorization scheme is not supported.
var ErrInvalidScheme = errors.New("session: authorization scheme not supported")

// BeginSession creates a new SessionID, saves the `sessionState` to the store, adds an
// Authorization header to the response with the SessionID, and returns the new SessionID.
func BeginSession(signingKey string, store Store, sessionState interface{}, w http.ResponseWriter) (SessionID, error) {
	sesID, err := NewSessionID(signingKey)
	if err != nil {
		return InvalidSessionID, fmt.Errorf("BeginSession: error generating session id: %s", err.Error())
	}
	errSS := store.Save(sesID, sessionState)
	if errSS != nil {
		return InvalidSessionID, fmt.Errorf("BeginSession: error saving session: %s", errSS.Error())
	}
	w.Header().Add(HeaderAuthorization, SchemeBearer+string(sesID))
	return sesID, nil
}

// GetSessionID extracts and validates the SessionID from the request headers.
func GetSessionID(r *http.Request, signingKey string) (SessionID, error) {
	authString := r.Header.Get(HeaderAuthorization)
	if len(authString) == 0 {
		authString = r.FormValue(ParamAuthorization)
		if len(authString) == 0 {
			return InvalidSessionID, ErrNoSessionID
		}
	}
	if strings.HasPrefix(authString, SchemeBearer) {
		authString = strings.TrimSpace(strings.Replace(authString, SchemeBearer, "", 1))
	} else {
		return InvalidSessionID, ErrInvalidScheme
	}
	sesID, errVID := ValidateID(authString, signingKey)
	if errVID != nil {
		return InvalidSessionID, errVID
	}
	return sesID, nil
}

// GetState extracts the SessionID from the request, gets the associated state from the provided store into
// the `sessionState` parameter, and returns the SessionID.
func GetState(r *http.Request, signingKey string, store Store, sessionState interface{}) (SessionID, error) {
	sesID, errGSID := GetSessionID(r, signingKey)
	if errGSID != nil {
		return InvalidSessionID, errGSID
	}
	errSG := store.Get(sesID, sessionState)
	if errSG != nil {
		return sesID, errSG
	}
	return sesID, nil
}

// EndSession extracts the SessionID from the request, and deletes the associated data in the provided store,
// returning the extracted SessionID.
func EndSession(r *http.Request, signingKey string, store Store) (SessionID, error) {
	sesID, errGSID := GetSessionID(r, signingKey)
	if errGSID != nil {
		return InvalidSessionID, errGSID
	}
	errSD := store.Delete(sesID)
	if errSD != nil {
		return sesID, errSD
	}
	return sesID, nil
}
