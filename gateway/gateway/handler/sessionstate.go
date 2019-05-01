package handler

import (
	"time"

	"github.com/uw-thalesians/perceptia-servers/gateway/gateway/user"

	"github.com/uw-thalesians/perceptia-servers/gateway/gateway/session"
)

// SessionState stores the session information for an authenticated user including the time the session started and
// the user's information.
type SessionState struct {
	SessionID session.SessionID `json:"-"`
	StartTime time.Time         `json:"startTime"`
	User      user.User         `json:"user"`
}

// NewSessionState constructs a new SessionState struct using the provided startTime and User.
func NewSessionState(startTime time.Time, user user.User) *SessionState {
	return &SessionState{StartTime: startTime, User: user}
}
