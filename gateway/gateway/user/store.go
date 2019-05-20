package user

import (
	"errors"

	"github.com/uw-thalesians/perceptia-servers/gateway/gateway/session"

	uuid "github.com/satori/go.uuid"
)

// ErrUserNotFound is returned when the user can't be found.
var ErrUserNotFound = errors.New("user not found")

// Store represents a store for Users.
//
// Store abstracts the common actions involving the database for users,
// abstracting the underlying interaction with the database.
type Store interface {
	// GetByUuid returns the User with the given Uuid.
	GetByUuid(uuid uuid.UUID) (*User, error)

	// GetByUsername returns the User with the given Username.
	GetByUsername(username string) (*User, error)

	// GetEncodedHashByUsername returns the Encoded Hash for the given user.
	GetEncodedHashByUsername(username string) (string, error)

	// GetSessionsByUuid returns the list of sessions the user has started.
	GetActiveSessionsByUuid(uuid uuid.UUID) (*session.Sessions, error)

	// Insert inserts the user into the database, and returns the newly-inserted User.
	Insert(newUser *NewUser) (*User, error)

	// InsertEmail adds the email to the given user's account
	InsertEmail(uuid uuid.UUID, email string) error

	// UpdateFullName updates the full name for the given user
	UpdateFullName(uuid uuid.UUID, fullName string) (*User, error)

	// UpdateDisplayName updates the display name for the given user
	UpdateDisplayName(uuid uuid.UUID, displayName string) (*User, error)

	// UpdateEncodedHash updates the encoded hash stored for the user
	UpdateEncodedHash(uuid uuid.UUID, encodedHash string) error

	// Delete deletes the user with the given ID.
	Delete(uuid uuid.UUID) error

	// DeleteEmail deletes the email from the given user's account
	DeleteEmail(uuid uuid.UUID, email string) error

	// DeleteSession deletes the session entry
	DeleteSession(uuid uuid.UUID, session uuid.UUID) error
}
