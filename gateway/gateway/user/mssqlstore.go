package user

import (
	"database/sql"
	"errors"

	mssql "github.com/denisenkom/go-mssqldb"

	uuid "github.com/satori/go.uuid"
)

type MsSqlStore struct {
	database *sql.DB
}

type userInfo struct {
	Uuid        mssql.UniqueIdentifier
	Username    string
	FullName    string
	DisplayName string
}

func (ms *MsSqlStore) GetByUuid(uuid uuid.UUID) (*User, error) {
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
}

func (ms *MsSqlStore) GetByUsername(username string) (*User, error) {
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
}

func (ms *MsSqlStore) GetEncodedHashByUsername(username string) (string, error) {
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
}

func (ms *MsSqlStore) Insert(newUser *NewUser) (*User, error) {
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
}

func (ms *MsSqlStore) UpdateFullName(uuid uuid.UUID, fullName string) (*User, error) {
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
}

func (ms *MsSqlStore) UpdateDisplayName(uuid uuid.UUID, displayName string) (*User, error) {
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
}

func (ms *MsSqlStore) UpdateEncodedHash(uuid uuid.UUID, encodedHash string) error {
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
}

func (ms *MsSqlStore) Delete(uuid uuid.UUID) error {
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
}

func (ms *MsSqlStore) DeleteEmail(uuid uuid.UUID, email string) error {
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
}
