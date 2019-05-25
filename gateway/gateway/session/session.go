// Package session provides tools for managing user sessions.
package session

import (
	"errors"
	"net/http"
	"strings"
	"time"

	uuid "github.com/satori/go.uuid"
)

const HeaderAuthorization = "Authorization"
const ParamAuthorization = "access_token"
const AuthHeaderSchemeBearerPrefix = "Bearer "

// ErrNoSessionId is used when no session ID was found in the Authorization header.
var ErrNoSessionId = errors.New("session: no session ID found in header " +
	HeaderAuthorization + " or query params " + ParamAuthorization)

// ErrInvalidScheme is used when the authorization scheme is not supported.
var ErrInvalidScheme = errors.New("session: authorization scheme not supported")

var ErrInvalidSessionId = errors.New("invalid session id")

var ErrUnexpected = errors.New("session: unexpected error occurred")

type SessionInfo struct {
	Uuid      uuid.UUID
	SessionId SessionID
	Created   time.Time
}

type Sessions map[SessionID]*SessionInfo

// CreateSession creates a new SessionID and session uuid.
func CreateSession(signingKey string) (SessionID, uuid.UUID, error) {
	sesID, err := NewSessionID(signingKey)
	sesUuid := uuid.NewV4()
	if err != nil {
		return InvalidSessionID, sesUuid, ErrUnexpected
	}
	return sesID, sesUuid, nil
}

// BeginSession saves the `sessionState` to the store, adds an
// Authorization header to the response with the SessionID, and returns the new SessionID.
func BeginSession(sessionId SessionID, sessionUuid uuid.UUID, store Store, sessionState interface{}, w http.ResponseWriter) error {

	errSS := store.Save(sessionId, sessionUuid, sessionState)
	if errSS != nil {
		return ErrUnexpected
	}
	w.Header().Add(HeaderAuthorization, AuthHeaderSchemeBearerPrefix+string(sessionId))
	return nil
}

// GetSessionID extracts and validates the SessionID from the request headers.
func GetSessionID(r *http.Request, signingKey string) (SessionID, error) {
	authString := r.Header.Get(HeaderAuthorization)
	if len(authString) == 0 {
		authString = r.FormValue(ParamAuthorization)
		if len(authString) == 0 {
			return InvalidSessionID, ErrNoSessionId
		}
	} else {
		if strings.HasPrefix(authString, AuthHeaderSchemeBearerPrefix) {
			authString = strings.TrimSpace(strings.Replace(authString, AuthHeaderSchemeBearerPrefix, "", 1))
		} else {
			return InvalidSessionID, ErrInvalidScheme
		}
	}
	sesID, errVID := ValidateID(authString, signingKey)
	if errVID != nil {
		return InvalidSessionID, ErrInvalidSessionId
	}
	return sesID, nil
}

func GetSessionIDByUuid(sessionUuid uuid.UUID, store Store) (SessionID, error) {
	sesId, err := store.GetSessionId(sessionUuid)
	if err != nil {
		return InvalidSessionID, ErrUnexpected
	}
	return sesId, nil
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
func EndSession(sessionId SessionID, store Store) error {
	errSD := store.Delete(sessionId)
	if errSD != nil {
		return errSD
	}
	return nil
}
