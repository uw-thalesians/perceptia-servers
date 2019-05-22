package handler

import (
	"time"

	uuid "github.com/satori/go.uuid"

	"github.com/uw-thalesians/perceptia-servers/gateway/gateway/user"

	"github.com/uw-thalesians/perceptia-servers/gateway/gateway/session"
)

// SessionState stores the session information for an authenticated user including the time the session started and
// the user's information.
type SessionState struct {
	SessionID     session.SessionID `json:"-"`
	SessionUuid   uuid.UUID         `json:"sessionUuid"`
	StartTime     time.Time         `json:"startTime"`
	Authenticated bool              `json:"authenticated"`
	User          *user.User        `json:"user"`
}

// NewSessionState constructs a new SessionState struct using the provided startTime and User.
func NewSessionState(startTime time.Time, user *user.User,
	sessionUuid uuid.UUID, sessionId session.SessionID, authenticated bool) *SessionState {
	return &SessionState{StartTime: startTime, User: user,
		SessionUuid: sessionUuid, SessionID: sessionId, Authenticated: authenticated}
}
