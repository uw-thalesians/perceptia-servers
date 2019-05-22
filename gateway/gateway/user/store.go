package user

import (
	"errors"

	uuid "github.com/satori/go.uuid"
)

// ErrUserNotFound is returned when the user can't be found.
var ErrDidNotComplete = errors.New("database operation did not complete")
var ErrUserNotFound = errors.New("user not found")
var ErrUserAlreadyExists = errors.New("user already exists")
var ErrUsernameUnavailable = errors.New("username not available")

var ErrPreparingQuery = errors.New("issue preparing query")

var ErrUnexpected = errors.New("unexpected error occurred")

//var ErrSessionUuidDoesNotExist = errors.New("session does not exist")
//var ErrSessionUuidAlreadyExists = errors.New("session does not exist")

// Store represents a store for Users.
//
// Store abstracts the common actions involving the database for users,
// abstracting the underlying interaction with the database.
type Store interface {
	// CREATE /////////////////////////////////////////////////////////////////////////////////////////////////////////

	// CreateSession creates a new session entry with the provided session information.
	//TODO: CreateSession(sessionUuid uuid.UUID, sessionId session.SessionID) error

	// CreateUserSession creates a new session entry and associates it with the given user.
	//TODO: CreateUserSession(userUuid uuid.UUID, sessionUuid uuid.UUID, sessionId session.SessionID) error

	// CreateUserSessionAssociation associates an existing session entry with the given user.
	//TODO: CreateUserSessionAssociation(userUuid uuid.UUID, sessionUuid uuid.UUID) error

	// CreateUser will add the new user to the database
	CreateUser(newUser *NewUser) (*User, error)

	// CreateUserEmail adds the email to the given user's account
	//TODO: CreateUserEmail(userUuid uuid.UUID, email string) error

	// READ /////////////////////////////////////////////////////////////////////////////////////////////////////////

	// ReadProcedureVersion gets the procedure version implemented in the database.
	//TODO: ReadProcedureVersion() *utility.SemVer

	// ReadUserActiveSessions gets the active sessions of the user by uuid.
	//TODO: ReadUserActiveSessions(userUuid uuid.UUID) (TODO: Define Type, error)

	// ReadUserSessions gets all the sessions associated with the user.
	//TODO: ReadUserSessions(userUuid uuid.UUID) (TODO: define type, error)

	// ReadUserDisplayName gets the display name for the given user.
	//TODO: ReadUserDisplayName(userUuid uuid.UUID) (string, error)

	// ReadUserFullName gets the full name for the given user.
	//TODO: ReadUserFullName(userUuid uuid.UUID)

	// ReadUserEmails gets the list of emails associated with the user.
	//TODO: ReadUserEmails(userUuid uuid.UUID) (TODO: define type, error)

	// ReadUserEncodedHash gets the encoded hash of the users password.
	ReadUserEncodedHash(username string) (string, error)

	// ReadUserInfo gets the basic information about the user.
	ReadUserInfo(userUuid uuid.UUID) (*User, error)

	// ReadUserProfile gets the profile information for the user.
	//TODO: ReadUserProfile(userUuid uuid.UUID) (TODO: define type, error)

	// ReadUserUsername gets the username for the given user.
	//TODO: ReadUserUsername(userUuid uuid.UUID) (string, error)

	// ReadUserUsernamesByEmail gets the usernames associated with a given email.
	//TODO: ReadUserUsernamesByEmail(email string) (TODO: define type, error)

	// ReadUserUuid gets the uuid for the user based on the given username.
	ReadUserUuid(username string) (*uuid.UUID, error)

	// UPDATE /////////////////////////////////////////////////////////////////////////////////////////////////////////

	// UpdateSessionExpired sets the given session's status to "Expired".
	//TODO: UpdateSessionExpired(sessionUuid uuid.UUID) error

	// UpdateUserDisplayName updates the display name of the user.
	//TODO: UpdateUserDisplayName(displayName string) error

	// UpdateUserEncodedHash updates the encoded hash associated with the user.
	//TODO: UpdateUserEncodedHash(userUuid uuid.UUID, encodedHash string) error

	// UpdateUserFullName updates the full name of the user.
	//TODO: UpdateUserFullName(userUuid uuid.UUID, fullName string) error

	// UpdateUserProfileBio updates the bio associated with the user's profile.
	//TODO: UpdateUserProfileBio(userUuid uuid.UUID, bio string) error

	// UpdateUserProfileGravatarUrl updates the gravatar url associated with the user's profile.
	//TODO: UpdateUserProfileGravatarUrl(userUuid uuid.UUID, gravatarUrl string) error

	// UpdateUserProfileSharingBio updates the public sharing preference for the bio
	// associated with the user's profile.
	//TODO: UpdateUserProfileSharingBio(userUuid uuid.UUID, share bool) error

	// UpdateUserProfileSharingDisplayName updates the public sharing preference for the display name
	// associated with the user's profile.
	//TODO: UpdateUserProfileSharingDisplayName(userUuid uuid.UUID, share bool) error

	// UpdateUserProfileSharingGravatarUrl updates the public sharing preference for the gravatar url
	// associated with the user's profile.
	//TODO: UpdateUserProfileSharingGravatarUrl(userUuid uuid.UUID, share bool) error

	// UpdateUserUsername updates the username for the given user.
	//TODO: UpdateUserUsername(userUuid uuid.UUID, username string) error

	// DELETE /////////////////////////////////////////////////////////////////////////////////////////////////////////

	// DeleteUser removes the user from the database.
	DeleteUser(userUuid uuid.UUID) error

	// DeleteUserEmail removes the given email from the users account.
	//TODO: DeleteUserEmail(userUuid uuid.UUID, email string)

	// DeleteSession removes the given session from the list
	//TODO: DeleteSession(sessionUuid uuid.UUID) error

}
