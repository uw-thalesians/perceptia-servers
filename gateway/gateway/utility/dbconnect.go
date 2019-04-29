package utility

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/url"
	"time"
)

var ErrInvalidDsn = errors.New("establish: provided dsn invalid")
var ErrUnableToPing = errors.New("establish: unable to ping database")

// Establish creates a connection to the provided dsn using the named sql driver.
// If the dsn is invalid an error will be logged and ErrInvalidDsn will be returned
// If the dsn is valid, but the Database does not respond to a ping,
// the ping error will be printed using the standard logger and the error will be returned along with the valid *sql.DB.
func Establish(driverName, dsn string, ping bool) (*sql.DB, error) {
	sqlDB, errSQLOpen := sql.Open(driverName, dsn)
	if errSQLOpen != nil {
		log.Printf("error: unable to open connection to MySQL database: %s", errSQLOpen.Error())
		return sqlDB, ErrInvalidDsn
	}

	if ping {
		errSQLPing := sqlDB.Ping()
		if errSQLPing != nil {
			for i := 0; i < 5; i++ {
				errSQLPing = sqlDB.Ping()
				if errSQLPing == nil {
					break
				}
				time.Sleep(time.Second * time.Duration(i))
			}
		}

		if errSQLPing != nil {
			log.Printf("Database connection established, but unable to ping database, got error: %s", errSQLPing)
			return sqlDB, ErrUnableToPing
		}
	}

	return sqlDB, nil
}

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
