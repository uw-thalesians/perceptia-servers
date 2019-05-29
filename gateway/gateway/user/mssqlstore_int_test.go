// +build all integration

package user

import (
	"database/sql"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/uw-thalesians/perceptia-servers/gateway/gateway/utility"
)

// Basic Tests to run
// All tests rely on an active connection to the database.
func TestMsSqlStore(t *testing.T) {
	tests := []struct {
		name     string
		function func(ms *MsSqlStore) func(t *testing.T)
	}{
		{"TestMsSqlStore_BasicCRUD",
			TestMsSqlStore_BasicCRUD,
		},
	}
	db := setupConnection()
	msSqlStore := MsSqlStore{database: db}
	for _, test := range tests {
		t.Run(test.name, test.function(&msSqlStore))
	}
}

// Basic tests define

// TestMsSqlStore_BasicCRUD runs tests designed to go through the basic CRUD ("Create Read Update Delete" cycle.
var TestMsSqlStore_BasicCRUD = func(ms *MsSqlStore) func(*testing.T) {
	return func(t *testing.T) {
		encodedHash, errCEH := CreateEncodedHash("TestIngPasswordHash")
		if errCEH != nil {
			t.Fatalf("unexpected error has occured when setting up test: error:%s", errCEH)
		}
		newUser := NewUser{
			Username:    fmt.Sprintf("TestUserName%d", time.Now().Unix()),
			FullName:    "The Full User Name",
			DisplayName: "Andrew",
			EncodedHash: encodedHash,
		}
		newUser.PrepNewUser()
		errVNU := newUser.ValidateNewUser()
		if errVNU != nil {
			t.Fatalf("unexpected error has occured when setting up test: error:%s", errVNU)
		}
		user := User{
			Username:    newUser.Username,
			FullName:    newUser.FullName,
			DisplayName: newUser.DisplayName,
		}

		tests := []struct {
			name                string
			newUser             *NewUser
			expectedUser        *User
			updateDisplayName   bool
			updateToDisplayName string
			expectErrorCreate   bool
			expectErrorRead     bool
			expectErrorUpdate   bool
			expectErrorDelete   bool
			detail              string
		}{
			{
				name:                "Basic Test",
				newUser:             &newUser,
				expectedUser:        &user,
				updateDisplayName:   true,
				updateToDisplayName: "Updated NAME",
				expectErrorCreate:   false,
				expectErrorDelete:   false,
				expectErrorRead:     false,
				expectErrorUpdate:   false,
				detail:              "Valid user, no errors should occur",
			},
		}
		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				// Create
				user, errIU := ms.Insert(test.newUser)
				if test.expectErrorCreate && errIU == nil {
					t.Fatalf("Test: %s; error expected creating user, but no error occured; detail: %s", test.name, test.detail)
				} else if !test.expectErrorCreate && errIU != nil {
					t.Fatalf("Test: %s; error not expected creating user, but error occured: %s; detail: %s", test.name, errIU, test.detail)
				}
				if test.expectedUser.DisplayName != user.DisplayName ||
					test.expectedUser.FullName != user.FullName ||
					test.expectedUser.Username != user.Username {
					t.Fatalf("Test: %s; user returned when creating user does not match expected user:\nuser returned:\n%+v\nexpected user:\n%+v\n; detail: %s", test.name, user, test.expectedUser, test.detail)
				}

				// Read
				userReadU, errGUU := ms.GetByUuid(user.Uuid)
				if test.expectErrorRead && errGUU == nil {
					t.Errorf("Test: %s; error expected reading user, but no error occured; detail: %s", test.name, test.detail)
				} else if !test.expectErrorRead && errGUU != nil {
					t.Errorf("Test: %s; error not expected reading user, but error occured: %s; detail: %s", test.name, errGUU, test.detail)
				}
				if test.expectedUser.DisplayName != userReadU.DisplayName ||
					test.expectedUser.FullName != userReadU.FullName ||
					test.expectedUser.Username != userReadU.Username {
					t.Errorf("Test: %s; user returned when reading user does not match expected user:\nuser returned:\n%+v\nexpected user:\n%+v\n; detail: %s", test.name, userReadU, test.expectedUser, test.detail)
				}

				userReadUsr, errGUU := ms.GetByUsername(user.Username)
				if test.expectErrorRead && errGUU == nil {
					t.Errorf("Test: %s; error expected reading user, but no error occured; detail: %s", test.name, test.detail)
				} else if !test.expectErrorRead && errGUU != nil {
					t.Errorf("Test: %s; error not expected reading user, but error occured: %s; detail: %s", test.name, errGUU, test.detail)
				}
				if test.expectedUser.DisplayName != userReadUsr.DisplayName ||
					test.expectedUser.FullName != userReadUsr.FullName ||
					test.expectedUser.Username != userReadUsr.Username {
					t.Errorf("Test: %s; user returned when reading user does not match expected user:\nuser returned:\n%+v\nexpected user:\n%+v\n; detail: %s", test.name, userReadUsr, test.expectedUser, test.detail)
				}

				// Update
				if test.updateDisplayName {
					userReadUpd, errUUD := ms.UpdateDisplayName(user.Uuid, test.updateToDisplayName)
					if test.expectErrorUpdate && errUUD == nil {
						t.Fatalf("Test: %s; error expected updating user, but no error occured; detail: %s", test.name, test.detail)
					} else if !test.expectErrorUpdate && errUUD != nil {
						t.Fatalf("Test: %s; error not expected updating user, but error occured: %s; detail: %s", test.name, errUUD, test.detail)
					}
					if test.updateToDisplayName != userReadUpd.DisplayName ||
						test.expectedUser.FullName != userReadUpd.FullName ||
						test.expectedUser.Username != userReadUpd.Username {
						test.expectedUser.DisplayName = test.updateToDisplayName
						t.Errorf("Test: %s; user returned when updating user does not match expected user:\nuser returned:\n%+v\nexpected user:\n%+v\n; detail: %s", test.name, userReadUpd, test.expectedUser, test.detail)
					}
				}

				// Delete
				errD := ms.Delete(user.Uuid)
				if test.expectErrorDelete && errD == nil {
					t.Fatalf("Test: %s; error expected deleting user, but no error occured; detail: %s", test.name, test.detail)
				} else if !test.expectErrorDelete && errD != nil {
					t.Fatalf("Test: %s; error not expected deleting user, but error occured: %s; detail: %s", test.name, errD, test.detail)
				}

			})
		}
	}
}

func setupConnection() *sql.DB {
	mssqlScheme, errMSS := utility.RequireEnv("MSSQL_SCHEME")
	// Fail if the value is not provided
	if errMSS != nil {
		log.Fatal(errMSS)
	}
	mssqlUsername, errMSU := utility.RequireEnv("MSSQL_USERNAME")
	// Fail if the value is not provided
	if errMSU != nil {
		log.Fatal(errMSU)
	}
	mssqlPassword, errMSP := utility.RequireEnv("MSSQL_PASSWORD")
	// Fail if the value is not provided
	if errMSP != nil {
		log.Fatal(errMSP)
	}
	mssqlHost, errMSH := utility.RequireEnv("MSSQL_HOST")
	// Fail if the value is not provided
	if errMSH != nil {
		log.Fatal(errMSH)
	}
	mssqlPort, errMSPO := utility.RequireEnv("MSSQL_PORT")
	// Fail if the value is not provided
	if errMSPO != nil {
		log.Fatal(errMSPO)
	}
	mssqlDatabase, errMSDB := utility.RequireEnv("MSSQL_DATABASE")
	// Fail if the value is not provided
	if errMSDB != nil {
		log.Fatal(errMSDB)
	}

	// Create DSN to use for connection to mssql
	mssqlDsn := utility.BuildDsn(mssqlScheme, mssqlUsername, mssqlPassword, mssqlHost, mssqlPort, mssqlDatabase)

	// Connect to mssql database
	mssqlDb, errEMSD := utility.Establish("sqlserver", mssqlDsn.String(), true)
	if errEMSD == utility.ErrInvalidDsn {
		log.Fatalf("The provided DSN: %s, was invalid: %s", mssqlDsn.String(), errEMSD)
	} else if errEMSD == utility.ErrUnableToPing {
		log.Fatalf("unable to ping the connection using the dsn: %s, error: %s", mssqlDsn.String(), errEMSD)
	} else if errEMSD != nil {
		log.Fatalf("unexpected error connecting to db using dsn: %s; error: %s", mssqlDsn.String(), errEMSD)
	}
	return mssqlDb
}
