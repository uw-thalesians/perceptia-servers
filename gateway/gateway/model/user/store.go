package user

import "errors"

// ErrUserNotFound is returned when the user can't be found.
var ErrUserNotFound = errors.New("user not found")

// Store represents a store for Users.
type Store interface {
	// GetByID returns the User with the given UUID.
	GetByUUID(uuid string) (*User, error)

	// GetByEmail returns the User with the given email.
	GetByEmail(email string) (*User, error)

	// GetByUserName returns the User with the given Username.
	GetByUserName(username string) (*User, error)

	// GetCredentialByUserName returns the Credential of the given user.
	GetCredentialByUserName(username string) (*Credential, error)

	// Insert inserts the user into the database, and returns the newly-inserted User.
	Insert(user *User) (*User, error)

	// Update applies UserUpdates to the given user UUID and returns the newly-updated user.
	Update(uuid string, updates *Updates) (*User, error)

	// Delete deletes the user with the given ID.
	Delete(uuid string) error
}
