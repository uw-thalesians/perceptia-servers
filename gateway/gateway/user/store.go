package user

import (
	"errors"

	uuid "github.com/satori/go.uuid"
)

// ErrUserNotFound is returned when the user can't be found.
var ErrUserNotFound = errors.New("user not found")

// Store represents a store for Users.
type Store interface {
	// GetByUUID returns the User with the given UUID.
	GetByUUID(uuid uuid.UUID) (*User, error)

	// GetByUsername returns the User with the given Username.
	GetByUsername(username string) (*User, error)

	// GetEncodedHashByUsername returns the Encoded Hash for the given user.
	GetEncodedHashByUsername(username string) (string, error)

	// Insert inserts the user into the database, and returns the newly-inserted User.
	Insert(user *User) (*User, error)

	// Delete deletes the user with the given ID.
	Delete(uuid uuid.UUID) error
}
