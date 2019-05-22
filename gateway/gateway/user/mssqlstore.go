package user

import (
	"database/sql"
	"errors"
	"time"

	"github.com/uw-thalesians/perceptia-servers/gateway/gateway/session"

	mssql "github.com/denisenkom/go-mssqldb"

	uuid "github.com/satori/go.uuid"
)

type MsSqlStore struct {
	database *sql.DB
}

// NewMsSqltore constructs a new MsSqlStore.
// If *sql.DB is nil, function will return an error.
func NewMsSqlStore(db *sql.DB) (*MsSqlStore, error) {
	if db == nil {
		return nil, errors.New("NewMsSqlStore: db cannot be nil")
	}
	return &MsSqlStore{
		db,
	}, nil
}

type userInfo struct {
	Uuid        mssql.UniqueIdentifier
	Username    string
	DisplayName string
}

type sessionInfo struct {
	Uuid      uuid.UUID
	SessionId session.SessionID
	Status    string
	Created   time.Time
}

// CREATE /////////////////////////////////////////////////////////////////////////////////////////////////////////

// CreateSession creates a new session entry with the provided session information.
//TODO: func (ms *MsSqlStore) CreateSession(sessionUuid uuid.UUID, sessionId session.SessionID) error

// CreateUserSession creates a new session entry and associates it with the given user.
//TODO: func (ms *MsSqlStore) CreateUserSession(userUuid uuid.UUID, sessionUuid uuid.UUID, sessionId session.SessionID) error

// CreateUserSessionAssociation associates an existing session entry with the given user.
//TODO: func (ms *MsSqlStore) CreateUserSessionAssociation(userUuid uuid.UUID, sessionUuid uuid.UUID) error

// CreateUser will add the new user to the database
func (ms *MsSqlStore) CreateUser(newUser *NewUser) (*User, error) {
	user := User{}
	userInfo := userInfo{}

	userUuid := uuid.NewV4()
	sqlUuid := mssql.UniqueIdentifier{}
	errSUID := sqlUuid.Scan(userUuid.String())
	if errSUID != nil {
		return &user, errSUID
	}

	stmt, errPS := ms.database.Prepare("USP_CreateUser")
	if errPS != nil {
		return &user, ErrPreparingQuery
	}
	defer stmt.Close()
	errQ := stmt.QueryRow(
		sql.Named("UserUuid", sqlUuid),
		sql.Named("Username", newUser.Username),
		sql.Named("FullName", newUser.FullName),
		sql.Named("DisplayName", newUser.DisplayName),
		sql.Named("EncodedHash", newUser.EncodedHash),
	).Scan(&userInfo.Uuid, &userInfo.Username, &userInfo.DisplayName)
	if errQ != nil {
		if errQ == sql.ErrNoRows {
			return &user, ErrDidNotComplete
		} else if msErr, ok := errQ.(mssql.Error); ok {
			if msErr.Number == 50401 {
				return &user, ErrUserAlreadyExists
			} else if msErr.Number == 50402 {
				return &user, ErrUsernameUnavailable
			}
		}
		return &user, errQ
	}
	errUQ := user.Uuid.Scan(userInfo.Uuid.String())
	if errUQ != nil {
		return &user, errUQ
	}
	user.DisplayName = userInfo.DisplayName
	user.Username = userInfo.Username
	return &user, nil
}

// CreateUserEmail adds the email to the given user's account
//TODO: func (ms *MsSqlStore) CreateUserEmail(userUuid uuid.UUID, email string) error

// READ /////////////////////////////////////////////////////////////////////////////////////////////////////////

// ReadProcedureVersion gets the procedure version implemented in the database.
//TODO: func (ms *MsSqlStore) ReadProcedureVersion() *utility.SemVer

// ReadUserActiveSessions gets the active sessions of the user by uuid.
//TODO: func (ms *MsSqlStore) ReadUserActiveSessions(userUuid uuid.UUID) (TODO: Define Type, error)

// ReadUserSessions gets all the sessions associated with the user.
//TODO: func (ms *MsSqlStore) ReadUserSessions(userUuid uuid.UUID) (TODO: define type, error)

// ReadUserDisplayName gets the display name for the given user.
//TODO: func (ms *MsSqlStore) ReadUserDisplayName(userUuid uuid.UUID) (string, error)

// ReadUserFullName gets the full name for the given user.
//TODO: func (ms *MsSqlStore) ReadUserFullName(userUuid uuid.UUID)

// ReadUserEmails gets the list of emails associated with the user.
//TODO: func (ms *MsSqlStore) ReadUserEmails(userUuid uuid.UUID) (TODO: define type, error)

// ReadUserEncodedHash gets the encoded hash of the users password.
func (ms *MsSqlStore) ReadUserEncodedHash(username string) (string, error) {
	stmt, errPS := ms.database.Prepare("USP_ReadUserEncodedHash")
	if errPS != nil {
		return InvalidEncodedPasswordHash, ErrPreparingQuery
	}
	defer stmt.Close()
	encodedHash := ""
	errQ := stmt.QueryRow(sql.Named("Username", username)).Scan(&encodedHash)
	if errQ != nil {
		if msErr, ok := errQ.(mssql.Error); ok {
			if msErr.Number == 50101 {
				return InvalidEncodedPasswordHash, ErrUnexpected
			} else if msErr.Number == 50301 {
				return InvalidEncodedPasswordHash, ErrUserNotFound
			}
		} else {
			return InvalidEncodedPasswordHash, ErrUnexpected
		}
	}
	return encodedHash, nil
}

// ReadUserInfo gets the basic information about the user.
func (ms *MsSqlStore) ReadUserInfo(userUuid uuid.UUID) (*User, error) {
	user := User{}
	userInfo := userInfo{}
	sqlUuid := mssql.UniqueIdentifier{}
	errUI := sqlUuid.Scan(userUuid.String())
	if errUI != nil {
		return &user, ErrUnexpected
	}

	stmt, errPS := ms.database.Prepare("USP_ReadUserInfo")
	if errPS != nil {
		return &user, ErrPreparingQuery
	}
	defer stmt.Close()
	errQ := stmt.QueryRow(sql.Named("UserUuid", sqlUuid)).Scan(
		&userInfo.Uuid, &userInfo.Username, &userInfo.DisplayName)
	if errQ != nil {
		if msErr, ok := errQ.(mssql.Error); ok {
			if msErr.Number == 50101 {
				return &user, ErrUnexpected
			} else if msErr.Number == 50301 {
				return &user, ErrUserNotFound
			} else {
				return &user, ErrUnexpected
			}
		}
	}
	errUQ := user.Uuid.Scan(userInfo.Uuid.String())
	if errUQ != nil {
		return &user, ErrUnexpected
	}
	user.DisplayName = userInfo.DisplayName
	user.Username = userInfo.Username
	return &user, nil
}

// ReadUserProfile gets the profile information for the user.
//TODO: func (ms *MsSqlStore) ReadUserProfile(userUuid uuid.UUID) (TODO: define type, error)

// ReadUserUsername gets the username for the given user.
//TODO: func (ms *MsSqlStore) ReadUserUsername(userUuid uuid.UUID) (string, error)

// ReadUserUsernamesByEmail gets the usernames associated with a given email.
//TODO: func (ms *MsSqlStore) ReadUserUsernamesByEmail(email string) (TODO: define type, error)

// ReadUserUuid gets the uuid for the user based on the given username.
func (ms *MsSqlStore) ReadUserUuid(username string) (*uuid.UUID, error) {

	sqlUuid := mssql.UniqueIdentifier{}

	stmt, errPS := ms.database.Prepare("USP_ReadUserUuid")
	if errPS != nil {
		return nil, ErrPreparingQuery
	}
	defer stmt.Close()
	errQ := stmt.QueryRow(sql.Named("Username", username)).Scan(
		&sqlUuid)
	if errQ != nil {
		if msErr, ok := errQ.(mssql.Error); ok {
			if msErr.Number == 50101 {
				return nil, ErrUnexpected
			} else if msErr.Number == 50301 {
				return nil, ErrUserNotFound
			} else {
				return nil, ErrUnexpected
			}
		}
	}
	userUuid := uuid.NewV4()
	errUQ := userUuid.Scan(sqlUuid.String())
	if errUQ != nil {
		return nil, ErrUnexpected
	}
	return &userUuid, nil
}

// UPDATE /////////////////////////////////////////////////////////////////////////////////////////////////////////

// UpdateSessionExpired sets the given session's status to "Expired".
//TODO: func (ms *MsSqlStore) UpdateSessionExpired(sessionUuid uuid.UUID) error

// UpdateUserDisplayName updates the display name of the user.
//TODO: func (ms *MsSqlStore) UpdateUserDisplayName(displayName string) error

// UpdateUserEncodedHash updates the encoded hash associated with the user.
//TODO: func (ms *MsSqlStore) UpdateUserEncodedHash(userUuid uuid.UUID, encodedHash string) error

// UpdateUserFullName updates the full name of the user.
//TODO: func (ms *MsSqlStore) UpdateUserFullName(userUuid uuid.UUID, fullName string) error

// UpdateUserProfileBio updates the bio associated with the user's profile.
//TODO: func (ms *MsSqlStore) UpdateUserProfileBio(userUuid uuid.UUID, bio string) error

// UpdateUserProfileGravatarUrl updates the gravatar url associated with the user's profile.
//TODO: func (ms *MsSqlStore) UpdateUserProfileGravatarUrl(userUuid uuid.UUID, gravatarUrl string) error

// UpdateUserProfileSharingBio updates the public sharing preference for the bio
// associated with the user's profile.
//TODO: func (ms *MsSqlStore) UpdateUserProfileSharingBio(userUuid uuid.UUID, share bool) error

// UpdateUserProfileSharingDisplayName updates the public sharing preference for the display name
// associated with the user's profile.
//TODO: func (ms *MsSqlStore) UpdateUserProfileSharingDisplayName(userUuid uuid.UUID, share bool) error

// UpdateUserProfileSharingGravatarUrl updates the public sharing preference for the gravatar url
// associated with the user's profile.
//TODO: func (ms *MsSqlStore) UpdateUserProfileSharingGravatarUrl(userUuid uuid.UUID, share bool) error

// UpdateUserUsername updates the username for the given user.
//TODO: func (ms *MsSqlStore) UpdateUserUsername(userUuid uuid.UUID, username string) error

// DELETE /////////////////////////////////////////////////////////////////////////////////////////////////////////

// DeleteUser removes the user from the database.
func (ms *MsSqlStore) DeleteUser(userUuid uuid.UUID) error {
	sqlUuid := mssql.UniqueIdentifier{}
	errSUID := sqlUuid.Scan(userUuid.String())
	if errSUID != nil {
		return ErrUnexpected
	}

	stmt, errPS := ms.database.Prepare("USP_DeleteUser")
	if errPS != nil {
		return ErrPreparingQuery
	}
	defer stmt.Close()
	rs, errQ := stmt.Exec(sql.Named("UserUuid", sqlUuid))
	if errQ != nil {
		if mssqlerr, ok := errQ.(mssql.Error); ok {
			if mssqlerr.Number == 50101 {
				return ErrUnexpected
			} else if mssqlerr.Number == 50301 {
				return ErrUserNotFound
			}
		} else {
			return ErrUnexpected
		}
	}
	if ra, errRA := rs.RowsAffected(); errRA != nil && ra < 1 {
		return errors.New("unexpected error: no rows deleted")
	}
	return nil
}

// DeleteUserEmail removes the given email from the users account.
//TODO: func (ms *MsSqlStore) DeleteUserEmail(userUuid uuid.UUID, email string)

// DeleteSession removes the given session from the list
//TODO: func (ms *MsSqlStore) DeleteSession(sessionUuid uuid.UUID) error

/*func (ms *MsSqlStore) GetByUuid(uuid uuid.UUID) (*User, error) {
	user := User{}
	userInfo := userInfo{}
	sqlUuid := mssql.UniqueIdentifier{}
	errUI := sqlUuid.Scan(uuid.String())
	if errUI != nil {
		return &user, errUI
	}

	stmt, errPS := ms.database.Prepare("USP_GetUserInfoByUuid")
	if errPS != nil {
		return &user, errPS
	}
	defer stmt.Close()
	errQ := stmt.QueryRow(sql.Named("Uuid", sqlUuid)).Scan(
		&userInfo.Uuid, &userInfo.Username, &userInfo.FullName, &userInfo.DisplayName)
	if errQ != nil {
		if errQ == sql.ErrNoRows {
			return &user, ErrUserNotFound
		} else {
			return &user, errQ
		}
	}
	errUQ := user.Uuid.Scan(userInfo.Uuid.String())
	if errUQ != nil {
		return &user, errUQ
	}
	user.DisplayName = userInfo.DisplayName
	user.FullName = userInfo.FullName
	user.Username = userInfo.Username
	return &user, nil
}*/

/*func (ms *MsSqlStore) GetActiveSessionsByUuid(uuid uuid.UUID) (*session.Sessions, error) {
	var sessions session.Sessions
	sqlUuid := mssql.UniqueIdentifier{}
	errUI := sqlUuid.Scan(uuid.String())
	if errUI != nil {
		return &sessions, errUI
	}

	stmt, errPS := ms.database.Prepare("USP_GetUserSessionsByUuid")
	if errPS != nil {
		return &sessions, errPS
	}
	defer stmt.Close()
	// TODO:
	//errQ := stmt.QueryRow(sql.Named("Uuid", sqlUuid)).Scan(
	//	&userInfo.Uuid, &userInfo.Username, &userInfo.FullName, &userInfo.DisplayName)
	//if errQ != nil {
	//	if errQ == sql.ErrNoRows {
	//		return &user, ErrUserNotFound
	//	} else {
	//		return &user, errQ
	//	}
	//}
	//errUQ := user.Uuid.Scan(userInfo.Uuid.String())
	//if errUQ != nil {
	//	return &user, errUQ
	//}
	//user.DisplayName = userInfo.DisplayName
	//user.FullName = userInfo.FullName
	//user.Username = userInfo.Username
	return &sessions, nil
}*/

/*func (ms *MsSqlStore) GetByUsername(username string) (*User, error) {
	user := User{}
	userInfo := userInfo{}

	stmt, errPS := ms.database.Prepare("USP_GetUserInfoByUsername")
	if errPS != nil {
		return &user, errPS
	}
	defer stmt.Close()
	errQ := stmt.QueryRow(sql.Named("Username", username)).Scan(
		&userInfo.Uuid, &userInfo.Username, &userInfo.FullName, &userInfo.DisplayName)
	if errQ != nil {
		if errQ == sql.ErrNoRows {
			return &user, ErrUserNotFound
		} else {
			return &user, errQ
		}
	}
	errUQ := user.Uuid.Scan(userInfo.Uuid.String())
	if errUQ != nil {
		return &user, errUQ
	}
	user.DisplayName = userInfo.DisplayName
	user.FullName = userInfo.FullName
	user.Username = userInfo.Username
	return &user, nil
}*/

/*func (ms *MsSqlStore) GetEncodedHashByUsername(username string) (string, error) {
	stmt, errPS := ms.database.Prepare("USP_GetUserEncodedHashByUsername")
	if errPS != nil {
		return InvalidEncodedPasswordHash, errPS
	}
	defer stmt.Close()
	encodedHash := ""
	errQ := stmt.QueryRow(sql.Named("Username", username)).Scan(&encodedHash)
	if errQ != nil {
		if errQ == sql.ErrNoRows {
			return InvalidEncodedPasswordHash, ErrUserNotFound
		} else {
			return InvalidEncodedPasswordHash, errQ
		}
	}
	return encodedHash, nil
}*/

/*func (ms *MsSqlStore) Insert(newUser *NewUser) (*User, error) {
	user := User{}
	userInfo := userInfo{}

	userUuid := uuid.NewV4()
	sqlUuid := mssql.UniqueIdentifier{}
	errSUID := sqlUuid.Scan(userUuid.String())
	if errSUID != nil {
		return &user, errSUID
	}

	stmt, errPS := ms.database.Prepare("USP_CreateUser")
	if errPS != nil {
		return &user, errPS
	}
	defer stmt.Close()
	errQ := stmt.QueryRow(
		sql.Named("Uuid", sqlUuid),
		sql.Named("Username", newUser.Username),
		sql.Named("FullName", newUser.FullName),
		sql.Named("DisplayName", newUser.DisplayName),
		sql.Named("EncodedHash", newUser.EncodedHash),
	).Scan(&userInfo.Uuid, &userInfo.Username, &userInfo.FullName, &userInfo.DisplayName)
	if errQ != nil {
		if errQ == sql.ErrNoRows {
			return &user, ErrUserNotFound
		}
		return &user, errQ
	}
	errUQ := user.Uuid.Scan(userInfo.Uuid.String())
	if errUQ != nil {
		return &user, errUQ
	}
	user.DisplayName = userInfo.DisplayName
	user.FullName = userInfo.FullName
	user.Username = userInfo.Username
	return &user, nil
}

func (ms *MsSqlStore) InsertEmail(uuid uuid.UUID, email string) error {
	sqlUuid := mssql.UniqueIdentifier{}
	errUI := sqlUuid.Scan(uuid.String())
	if errUI != nil {
		return errUI
	}
	stmt, errPS := ms.database.Prepare("USP_AddUserEmail")
	if errPS != nil {
		return errPS
	}
	defer stmt.Close()
	rs, errQ := stmt.Exec(sql.Named("Uuid", sqlUuid), sql.Named("Email", email))
	if errQ != nil {
		if mssqlerr, ok := errQ.(mssql.Error); ok {
			if mssqlerr.Number == 50201 {
				return ErrUserNotFound
			}
		} else {
			return errQ
		}
	}
	if ra, errRA := rs.RowsAffected(); errRA != nil && ra < 1 {
		return errors.New("unexpected error: no rows updated")
	}
	return nil
}*/

/*func (ms *MsSqlStore) UpdateFullName(uuid uuid.UUID, fullName string) (*User, error) {
	user := User{}
	sqlUuid := mssql.UniqueIdentifier{}
	errSUID := sqlUuid.Scan(uuid.String())
	if errSUID != nil {
		return &user, errSUID
	}

	stmt, errPS := ms.database.Prepare("USP_UpdateUserFullName")
	if errPS != nil {
		return &user, errPS
	}
	defer stmt.Close()
	rs, errQ := stmt.Exec(
		sql.Named("Uuid", sqlUuid),
		sql.Named("FullName", fullName),
	)
	if errQ != nil {
		if mssqlerr, ok := errQ.(mssql.Error); ok {
			if mssqlerr.Number == 50201 {
				return &user, ErrUserNotFound
			}
		}
		return &user, errQ
	}
	if ra, errRA := rs.RowsAffected(); errRA != nil && ra < 1 {
		return &user, errors.New("unexpected error: no rows updated")
	}
	userRet, errGU := ms.GetByUuid(uuid)
	if errGU != nil {
		return &user, errGU
	}
	return userRet, nil
}*/

/*func (ms *MsSqlStore) UpdateDisplayName(uuid uuid.UUID, displayName string) (*User, error) {
	user := User{}

	sqlUuid := mssql.UniqueIdentifier{}
	errSUID := sqlUuid.Scan(uuid.String())
	if errSUID != nil {
		return &user, errSUID
	}

	stmt, errPS := ms.database.Prepare("USP_UpdateUserDisplayName")
	if errPS != nil {
		return &user, errPS
	}
	defer stmt.Close()
	rs, errQ := stmt.Exec(
		sql.Named("Uuid", sqlUuid),
		sql.Named("DisplayName", displayName),
	)
	if errQ != nil {
		if mssqlerr, ok := errQ.(mssql.Error); ok {
			if mssqlerr.Number == 50201 {
				return &user, ErrUserNotFound
			}
		}
		return &user, errQ
	}
	if ra, errRA := rs.RowsAffected(); errRA != nil && ra < 1 {
		return &user, errors.New("unexpected error: no rows updated")
	}
	userRet, errGU := ms.GetByUuid(uuid)
	if errGU != nil {
		return &user, errGU
	}
	return userRet, nil
}*/

/*func (ms *MsSqlStore) UpdateEncodedHash(uuid uuid.UUID, encodedHash string) error {
	sqlUuid := mssql.UniqueIdentifier{}
	errSUID := sqlUuid.Scan(uuid.String())
	if errSUID != nil {
		return errSUID
	}

	stmt, errPS := ms.database.Prepare("USP_UpdateUserEncodedHash")
	if errPS != nil {
		return errPS
	}
	defer stmt.Close()
	rs, errQ := stmt.Exec(
		sql.Named("Uuid", sqlUuid),
		sql.Named("EncodedHash", encodedHash),
	)
	if errQ != nil {
		if mssqlerr, ok := errQ.(mssql.Error); ok {
			if mssqlerr.Number == 50201 {
				return ErrUserNotFound
			}
		}
		return errQ
	}
	if ra, errRA := rs.RowsAffected(); errRA != nil && ra < 1 {
		return errors.New("unexpected error: no rows updated")
	}
	return nil
}*/

/*func (ms *MsSqlStore) Delete(uuid uuid.UUID) error {
	sqlUuid := mssql.UniqueIdentifier{}
	errSUID := sqlUuid.Scan(uuid.String())
	if errSUID != nil {
		return errSUID
	}

	stmt, errPS := ms.database.Prepare("USP_DeleteUser")
	if errPS != nil {
		return errPS
	}
	defer stmt.Close()
	rs, errQ := stmt.Exec(sql.Named("Uuid", sqlUuid))
	if errQ != nil {
		if mssqlerr, ok := errQ.(mssql.Error); ok {
			if mssqlerr.Number == 50201 {
				return ErrUserNotFound
			}
		} else {
			return errQ
		}
	}
	if ra, errRA := rs.RowsAffected(); errRA != nil && ra < 1 {
		return errors.New("unexpected error: no rows deleted")
	}
	return nil
}*/

/*func (ms *MsSqlStore) DeleteEmail(uuid uuid.UUID, email string) error {
	sqlUuid := mssql.UniqueIdentifier{}
	errSUID := sqlUuid.Scan(uuid.String())
	if errSUID != nil {
		return errSUID
	}

	stmt, errPS := ms.database.Prepare("USP_DeleteUserEmail")
	if errPS != nil {
		return errPS
	}
	defer stmt.Close()
	rs, errQ := stmt.Exec(sql.Named("Uuid", sqlUuid), sql.Named("Email", email))
	if errQ != nil {
		if mssqlerr, ok := errQ.(mssql.Error); ok {
			if mssqlerr.Number == 50201 {
				return ErrUserNotFound
			} else if mssqlerr.Number == 50202 {
				return errors.New("provided email is not associated with users account")
			}
		} else {
			return errQ
		}
	}
	if ra, errRA := rs.RowsAffected(); errRA != nil && ra < 1 {
		return errors.New("unexpected error: no rows deleted")
	}
	return nil
}*/

/*func (ms *MsSqlStore) DeleteSession(uuid uuid.UUID, session uuid.UUID) error {
	sqlUuid := mssql.UniqueIdentifier{}
	errSUID := sqlUuid.Scan(uuid.String())
	if errSUID != nil {
		return errSUID
	}
	sqlSessionUuid := mssql.UniqueIdentifier{}
	errSSID := sqlUuid.Scan(session.String())
	if errSSID != nil {
		return errSSID
	}

	stmt, errPS := ms.database.Prepare("USP_DeleteUserSession")
	if errPS != nil {
		return errPS
	}
	defer stmt.Close()
	rs, errQ := stmt.Exec(sql.Named("Uuid", sqlUuid), sql.Named("SessionUuid", sqlSessionUuid))
	if errQ != nil {
		if mssqlerr, ok := errQ.(mssql.Error); ok {
			if mssqlerr.Number == 50201 {
				return ErrUserNotFound
			} else if mssqlerr.Number == 50202 {
				return errors.New("provided session uuid is not associated with user's account")
			}
		} else {
			return errQ
		}
	}
	if ra, errRA := rs.RowsAffected(); errRA != nil && ra < 1 {
		return errors.New("unexpected error: no rows deleted")
	}
	return nil
}*/
