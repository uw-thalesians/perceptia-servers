/*
Program gateway is an HTTP API Gateway server for the Perceptia application.

Gateway serves a REST api for the Perceptia application, proxying requests to the backend service responsible for a given collection.
*/
package main

//noinspection SpellCheckingInspection
import (
	"context"

	"github.com/uw-thalesians/perceptia-servers/gateway/gateway/session"

	"github.com/uw-thalesians/perceptia-servers/gateway/gateway/user"

	"github.com/go-redis/redis"
	"github.com/uw-thalesians/perceptia-servers/gateway/gateway/handler"

	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/denisenkom/go-mssqldb"
	kitlog "github.com/go-kit/kit/log"
	"github.com/uw-thalesians/perceptia-servers/gateway/gateway/utility"

	"github.com/gorilla/mux"
)

// SessionDuration is the time a session is valid. If it has been longer than time.Duration since
// the session was started the session should be treated as no longer valid.
const SessionDuration = time.Duration(time.Hour * 48)

// sqlDriverName is the name of the SQL driver to register with the go sql lib
const sqlDriverName = "sqlserver"

// services and the name of the collection to proxy on
const (
	serviceAqRest = "anyquiz"
)

// gateway provided collections
const (
	colUsers    = "users"
	colSessions = "sessions"
)

const GatewayServiceApiVersion = "0.3.0"

func main() {
	// Setup Logger
	logger := kitlog.NewLogfmtLogger(kitlog.NewSyncWriter(os.Stdout))
	logger = kitlog.With(logger, "ts", kitlog.DefaultTimestampUTC, "caller", kitlog.DefaultCaller)

	// Get address for server to listen for requests on
	listenAddr := utility.DefaultEnv("GATEWAY_LISTEN_ADDR", ":443")

	// Get the directory path to the TLS key and cert
	tlsCertPath, errTLSCertPath := utility.RequireEnv("GATEWAY_TLSCERTPATH")
	// Fail if the path to the cert is not provided
	if errTLSCertPath != nil {
		_ = logger.Log("error", errTLSCertPath, "result", "exit")
		os.Exit(1)
	}
	tlsKeyPath, errTLSKeyPath := utility.RequireEnv("GATEWAY_TLSKEYPATH")
	// Fail if the path to the key is not provided
	if errTLSKeyPath != nil {
		_ = logger.Log("error", errTLSKeyPath, "result", "exit")
		os.Exit(1)
	}

	sessionSigningKey, errSSK := utility.RequireEnv("GATEWAY_SESSION_KEY")
	// Fail if the key data is not provided
	if errSSK != nil {
		_ = logger.Log("error", errSSK, "result", "exit")
		os.Exit(1)
	}

	mssqlScheme, errMSS := utility.RequireEnv("MSSQL_SCHEME")
	// Fail if the value is not provided
	if errMSS != nil {
		_ = logger.Log("error", errMSS, "result", "exit")
		os.Exit(1)
	}
	mssqlUsername, errMSU := utility.RequireEnv("MSSQL_USERNAME")
	// Fail if the value is not provided
	if errMSU != nil {
		_ = logger.Log("error", errMSU, "result", "exit")
		os.Exit(1)
	}
	mssqlPassword, errMSP := utility.RequireEnv("MSSQL_PASSWORD")
	// Fail if the value is not provided
	if errMSP != nil {
		_ = logger.Log("error", errMSP, "result", "exit")
		os.Exit(1)
	}
	mssqlHost, errMSH := utility.RequireEnv("MSSQL_HOST")
	// Fail if the value is not provided
	if errMSH != nil {
		_ = logger.Log("error", errMSH, "result", "exit")
		os.Exit(1)
	}
	mssqlPort, errMSPO := utility.RequireEnv("MSSQL_PORT")
	// Fail if the value is not provided
	if errMSPO != nil {
		_ = logger.Log("error", errMSPO, "result", "exit")
		os.Exit(1)
	}
	mssqlDatabase, errMSDB := utility.RequireEnv("MSSQL_DATABASE")
	// Fail if the value is not provided
	if errMSDB != nil {
		_ = logger.Log("error", errMSDB, "result", "exit")
		os.Exit(1)
	}

	// Read in Redis connection string
	redisAddress, errRDAD := utility.RequireEnv("REDIS_ADDRESS")
	// Fail if the value is not provided
	if errRDAD != nil {
		_ = logger.Log("error", errRDAD, "result", "exit")
		os.Exit(1)
	}

	// Read in service names
	aqRestHostname, errAQHN := utility.RequireEnv("AQREST_HOSTNAME")
	// Fail if the value is not provided
	if errAQHN != nil {
		_ = logger.Log("error", errAQHN, "result", "exit")
		os.Exit(1)
	}
	aqRestPort, errAQPN := utility.RequireEnv("AQREST_PORT")
	// Fail if the value is not provided
	if errAQPN != nil {
		_ = logger.Log("error", errAQPN, "result", "exit")
		os.Exit(1)
	}

	// Create DSN to use for connection to mssql
	mssqlDsn := utility.BuildDsn(mssqlScheme, mssqlUsername, mssqlPassword, mssqlHost, mssqlPort, mssqlDatabase)

	// Connect to mssql database
	perceptiaDb, errEMSD := utility.Establish(sqlDriverName, mssqlDsn.String(), true)
	if errEMSD == utility.ErrInvalidDsn {
		_ = logger.Log("msg", "provided dsn invalid", "dsn", mssqlDsn.String(), "error", errEMSD, "result", "exit")
		os.Exit(1)
	} else if errEMSD == utility.ErrUnableToPing {
		_ = logger.Log("msg", "unable to ping database", "dsn", mssqlDsn.String(), "error", errEMSD)
	}

	// Periodically check status of mssql database connection
	pingDbCtx := context.TODO()
	go utility.PingDatabase(pingDbCtx, perceptiaDb, time.Second*10, time.Minute)

	//Create a new Redis client.
	rc := redis.NewClient(&redis.Options{Addr: redisAddress, Password: "", DB: 0})

	// Setup Stores
	userStore, errNMSDB := user.NewMsSqlStore(perceptiaDb)
	if errNMSDB != nil {
		_ = logger.Log("error", errNMSDB, "result", "exit")
		os.Exit(1)
	}

	sessionStore := session.NewRedisStore(rc, SessionDuration)

	// Create Handler Context
	hcx := handler.NewContext(sessionStore, userStore, sessionSigningKey, GatewayServiceApiVersion, logger)

	// Create new mux router
	gmux := mux.NewRouter()

	gmuxApi := gmux.PathPrefix("/api/").Subrouter()

	// "/api/v1/"
	gmuxApiV := gmuxApi.PathPrefix("/{" + handler.ReqVarMajorVersion + ":v[0-9]+}/").Subrouter()

	//// Service Routes

	// "/api/vX/anyquiz/"
	gmuxApiV.PathPrefix("/" + serviceAqRest + "/").Handler(hcx.NewServiceProxy(aqRestHostname, aqRestPort))

	//// Gateway routes
	gmuxApiVGateway := gmuxApiV.PathPrefix("/gateway/").Subrouter()

	// Health check route
	gmuxApiVGateway.HandleFunc("/health", hcx.HealthHandler)

	// Users route
	gmuxApiVGateway.HandleFunc("/"+colUsers, hcx.UsersDefaultHandler)

	// Sessions route

	gmuxApiVGateway.HandleFunc("/"+colSessions, hcx.SessionsDefaultHandler)

	// Users Subroutes
	gmuxApiVGatewayUsers := gmuxApiVGateway.PathPrefix("/" + colUsers + "/").Subrouter()

	// Users Specific routes
	gmuxApiVGatewayUsersSpecific := gmuxApiVGatewayUsers.PathPrefix("/{" + handler.ReqVarUserUuid +
		":[0-9A-Fa-f]{8}-[0-9A-Fa-f]{4}-4[0-9A-Fa-f]{3}-[89aAbB][0-9A-Fa-f]{3}-[0-9A-Fa-f]{12}}").Subrouter()

	gmuxApiVGatewayUsersSpecific.PathPrefix("").HandlerFunc(hcx.UsersSpecificHandler)

	// Sessions Subroutes
	gmuxApiVGatewaySessions := gmuxApiVGateway.PathPrefix("/" + colSessions + "/").Subrouter()

	// Sessions Specific routes
	gmuxApiVGatewaySessionsSpecific := gmuxApiVGatewaySessions.PathPrefix(
		"/{" + handler.ReqVarSession + ":" + handler.SpecificSessionHandlerDeleteUserAlias + "}").Subrouter()

	gmuxApiVGatewaySessionsSpecific.PathPrefix("").HandlerFunc(hcx.SessionsSpecificHandler)

	// Add Middleware to "/api"
	gmuxApi.Use(handler.NewCors)
	gmuxApi.Use(hcx.NewAuthenticator)

	// Add Middleware to "/api/{majorVersion}"

	// Add Middleware to "/api/{majorVersion}/gateway"

	// Add Middleware to "/api/{majorVersion}/gateway/users/{uuid}"
	gmuxApiVGatewayUsersSpecific.Use(hcx.NewEnsureAuth)

	// Add Middleware to "/api/{majorVersion}/anyquiz"

	//Starts listening at the address set, and passes requests at that address
	//to the mux. Exits if ListenAndServerTLS fails
	log.Printf("server is listening at https://%s...", listenAddr)
	log.Fatal(http.ListenAndServeTLS(listenAddr, tlsCertPath, tlsKeyPath, gmux))
}
