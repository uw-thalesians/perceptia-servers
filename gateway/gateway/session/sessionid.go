package session

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
)

// InvalidSessionID represents an empty, invalid session ID.
const InvalidSessionID SessionID = ""

// idLength is the length of the ID portion.
const idLength = 32

// signedLength is the full length of the signed session ID.
// (ID portion plus signature.)
const signedLength = idLength + sha256.Size

// SessionID represents a valid, digitally-signed session ID.
//
// This is a base64 URL encoded string created from a byte slice where the first `idLength` bytes are
// crytographically random bytes representing the unique session ID, and the remaining bytes are
// an HMAC hash of those ID bytes (i.e., a digital signature).
// The byte slice layout is like so:
// +-----------------------------------------------------+
// |...32 crypto random bytes...|HMAC hash of those bytes|
// +-----------------------------------------------------+
type SessionID string

// ErrInvalidID is returned when an invalid session id is passed to ValidateID().
var ErrInvalidID = errors.New("invalid Session ID")

// NewSessionID creates and returns a new digitally-signed session ID, using `signingKey` as the HMAC signing key.
//
// An error is returned only if there was an error generating random bytes for the session ID, or
// an invalid signingKey was provided.
func NewSessionID(signingKey string) (SessionID, error) {
	if len(signingKey) == 0 {
		return InvalidSessionID, errors.New("NewSessionID: signingKey must have length greater than zero")
	}
	// Creates slice with length of the id
	result := make([]byte, idLength, signedLength)
	// Creates id by filling id portion with crypto rand numbers
	_, errRD := rand.Read(result)
	if errRD != nil {
		return InvalidSessionID, fmt.Errorf("NewSessionID: error generating session id: %s", errRD.Error())
	}
	resultMAC, errCMAC := createMAC(result, []byte(signingKey))
	if errCMAC != nil {
		return InvalidSessionID, fmt.Errorf("NewSessionID: error generating session id: %s", errCMAC.Error())
	}
	result = append(result, resultMAC...)
	resEnc := base64.URLEncoding.EncodeToString(result)
	// Return the result as a SessionID
	return SessionID(resEnc), nil
}

// ValidateID validates the string in the `id` parameter using the `signingKey` as the HMAC signing key
// and returns an error if invalid, or a SessionID if valid.
func ValidateID(id string, signingKey string) (SessionID, error) {
	if len(signingKey) == 0 {
		return InvalidSessionID, errors.New("NewSessionID: signingKey must have length greater than zero")
	}
	idDecoded, err := base64.URLEncoding.DecodeString(id)
	if err != nil {
		return InvalidSessionID, fmt.Errorf("ValidateID: error base64 decoding: %s", err.Error())
	}
	message := idDecoded[:idLength]
	messageMAC := idDecoded[idLength:]

	expectedMAC, errCMAC := createMAC(message, []byte(signingKey))
	if errCMAC != nil {
		return InvalidSessionID, fmt.Errorf("ValidateID: error generating mac of provided message: %s", errCMAC)
	}
	ck := hmac.Equal(messageMAC, expectedMAC)
	if ck {
		return SessionID(id), nil
	}
	return InvalidSessionID, ErrInvalidID
}

// String returns a string representation of the sessionID.
func (sid SessionID) String() string {
	return string(sid)
}

// createMAC creates a MAC from a `message` and a `signingKey`.
func createMAC(message, signingKey []byte) ([]byte, error) {
	mac := hmac.New(sha256.New, signingKey)
	_, err := mac.Write(message)
	if err != nil {
		return nil, fmt.Errorf("createID: issue creating mac of message: %s", err)
	}
	return mac.Sum(nil), nil
}
