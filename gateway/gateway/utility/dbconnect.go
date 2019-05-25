package utility

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"time"

	kitlog "github.com/go-kit/kit/log"
)

// BuildDsn uses the provided values to build a URL based DSN to connect to a database.
//
// Parameters
// 		scheme is the connection scheme, such as "sqlserver"
// 		username: to authenticate to the server with, such as "sa"
// 		password: for the account given by the username
// 		hostname: for the server hosting the database, such as "localhost"
// 		port: for the port the server is listening for a connection on
// 		database: for the database to use for all requests using this connection
//
// Outputs
// 		dsn: database connection string in format: "<scheme>://<username>:<password>@<hostname>:<port>?params=p1"
func BuildDsn(scheme, username, password, hostname, port, database string) (dsn *url.URL) {
	query := url.Values{}
	query.Add("app name", "gateway")
	query.Add("database", database)
	return &url.URL{
		Scheme:   scheme,
		User:     url.UserPassword(username, password),
		Host:     fmt.Sprintf("%s:%s", hostname, port),
		RawQuery: query.Encode(),
	}
}

// PingDatabase will periodically check to see if the database connection is still open and the database is accessible.
//
//
// Parameters
//
//
//		ctx: basic context
//
//		db: the connection for the database that should be pinged
//
//		sleepFailTime: the amount of time PingDatabase should wait to ping again after a failed ping
//
//		mssqlRequiredVersion: checks the version of the stored procedures the database is exposing
//			Will log an error if a version other than the one specified is exposed
//
//		logger: a Logger to log any issues/errors
//
//
// Outputs none
func PingDatabase(ctx context.Context, db *sql.DB, sleepFailTime time.Duration,
	sleepTestTime time.Duration, mssqlRequiredVersion *SemVer, logger kitlog.Logger,
	statusNotOkay chan bool) {
	for {
		select {
		case <-ctx.Done():
			_ = logger.Log("PingDatabase: ping check canceled")

			return
		default:

			if err := db.Ping(); err != nil {
				_ = logger.Log("func", "utility.PingDatabase", "pingError", err.Error(), "note", "will retry in "+sleepTestTime.String())
				select {
				case _, _ = <-statusNotOkay:
					break
				default:
					break
				}
				statusNotOkay <- true
				time.Sleep(sleepFailTime)
			} else {
				row := db.QueryRow("USP_ReadProcedureVersion")
				type versionVal struct {
					Version string `json:"version"`
				}
				var vv versionVal
				errS := row.Scan(&vv.Version)
				if errS == nil {
					semVerProc, errSVS := SemVerFromString(vv.Version)
					if errSVS != nil {
						_ = logger.Log("utility.PingDatabase", "invalid version string provided")
						select {
						case _, _ = <-statusNotOkay:
							break
						default:
							break
						}
						statusNotOkay <- true
					} else if semVerProc.Compare(mssqlRequiredVersion) < 0 {
						_ = logger.Log("utility.PingDatabase", "unsupported database version",
							"minimumVersionRequired", mssqlRequiredVersion.String(),
							"versionConnected", semVerProc.String())
						select {
						case _, _ = <-statusNotOkay:
							break
						default:
							break
						}
						statusNotOkay <- true
					}
				} else {
					_ = logger.Log("utility.PingDatabase", "unable to get proc version from database", "error", errS.Error())
					select {
					case _, _ = <-statusNotOkay:
						break
					default:
						break
					}
					statusNotOkay <- true
				}
			}
			select {
			case _, _ = <-statusNotOkay:
				break
			default:
				break
			}
			statusNotOkay <- false
			time.Sleep(sleepTestTime)
		}
	}
}
