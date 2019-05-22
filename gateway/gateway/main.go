/*
Program gateway is an HTTP API Gateway server for the Perceptia application.

Gateway serves a REST api for the Perceptia application, proxying requests to the backend service responsible for a given collection.
*/
package main

//noinspection SpellCheckingInspection
import (
	"context"
	"database/sql"

	"github.com/uw-thalesians/perceptia-servers/gateway/gateway/session"

	"github.com/uw-thalesians/perceptia-servers/gateway/gateway/user"

	"github.com/go-redis/redis"
	"github.com/uw-thalesians/perceptia-servers/gateway/gateway/handler"

	"net/http"
	"os"
	"time"

	_ "github.com/denisenkom/go-mssqldb"
	kitlog "github.com/go-kit/kit/log"
	"github.com/uw-thalesians/perceptia-servers/gateway/gateway/utility"

	"github.com/gorilla/mux"
)

// sessionDuration is the time a session is valid.
const sessionDuration = time.Duration(time.Hour * 48)

// sqlDriverName is the name of the SQL driver to register with the go sql lib
const sqlDriverName = "sqlserver"

// Services to match on
const (
	serviceAqRest  = "anyquiz"
	serviceGateway = "gateway"
)

// gateway provided collections
const (
	colUsers    = "users"
	colSessions = "sessions"
	colHealth   = "health"
)

const (
	uuidV4Regex = "[0-9A-Fa-f]{8}-[0-9A-Fa-f]{4}-4[0-9A-Fa-f]{3}-[89aAbB][0-9A-Fa-f]{3}-[0-9A-Fa-f]{12}"
	apiVXRegex  = "v[0-9]+"
)

//
var gatewayServiceApiVersion, _ = utility.NewSemVer(1, 0, 0)

var gatewayServiceApiVersionsSupported = map[int]*utility.SemVer{gatewayServiceApiVersion.GetMajor(): gatewayServiceApiVersion}

var mssqlRequiredVersion, _ = utility.NewSemVer(1, 0, 0)

func main() {
	logger := kitlog.NewLogfmtLogger(kitlog.NewSyncWriter(os.Stdout))
	// Get environment setting
	gwEnv, _ := logEnvVar(logger, "GATEWAY_ENVIRONMENT", "development", false)

	// Setup Logger
	logger = kitlog.With(logger, "ts", kitlog.DefaultTimestampUTC, "caller", kitlog.DefaultCaller, "env", gwEnv)

	// Get address for server to listen for requests on
	listenAddr, _ := logEnvVar(logger, "GATEWAY_LISTEN_ADDR", ":443", false)

	// Get the directory path to the TLS key and cert
	tlsCertPath := exitOnEnvError(logger, "GATEWAY_TLSCERTPATH")
	tlsKeyPath := exitOnEnvError(logger, "GATEWAY_TLSKEYPATH")

	sessionSigningKey := exitOnEnvError(logger, "GATEWAY_SESSION_KEY")

	mssqlScheme := exitOnEnvError(logger, "MSSQL_SCHEME")

	mssqlUsername := exitOnEnvError(logger, "MSSQL_USERNAME")

	mssqlPassword := exitOnEnvError(logger, "MSSQL_PASSWORD")

	mssqlHost := exitOnEnvError(logger, "MSSQL_HOST")

	mssqlPort := exitOnEnvError(logger, "MSSQL_PORT")

	mssqlDatabase := exitOnEnvError(logger, "MSSQL_DATABASE")

	redisAddress := exitOnEnvError(logger, "REDIS_ADDRESS")

	aqRestHostname := exitOnEnvError(logger, "AQREST_HOSTNAME")

	aqRestPort := exitOnEnvError(logger, "AQREST_PORT")

	// Create DSN to use for connection to mssql
	mssqlDsn := utility.BuildDsn(mssqlScheme, mssqlUsername, mssqlPassword, mssqlHost, mssqlPort, mssqlDatabase)

	// Connect to mssql database
	perceptiaDb, errEMSD := sql.Open(sqlDriverName, mssqlDsn.String())
	if errEMSD != nil {
		if gwEnv == "development" {
			_ = logger.Log("msg", "unable to connect to db", "dsn", mssqlDsn.String(), "error", errEMSD, "result", "exit")
		} else {
			_ = logger.Log("msg", "unable to connect to db", "dsn", "hidden", "error", errEMSD, "result", "exit")
		}
		os.Exit(1)
	}

	// Periodically check status of mssql database connection
	pingDbCtx := context.TODO()
	go utility.PingDatabase(pingDbCtx, perceptiaDb, time.Second*10, time.Minute, mssqlRequiredVersion, logger)

	//Create a new Redis client.
	rc := redis.NewClient(&redis.Options{Addr: redisAddress, Password: "", DB: 0})

	// Setup Stores
	userStore, errNMSDB := user.NewMsSqlStore(perceptiaDb)
	if errNMSDB != nil {
		_ = logger.Log("error", errNMSDB, "result", "exit")
		os.Exit(1)
	}

	sessionStore := session.NewRedisStore(rc, sessionDuration)

	// Create Handler Context
	hcx := handler.NewContext(sessionStore, userStore, sessionSigningKey,
		gatewayServiceApiVersion, gatewayServiceApiVersionsSupported, logger)

	// Create new mux router
	gmux := mux.NewRouter()

	gmuxApi := gmux.PathPrefix("/api/").Subrouter()

	// "/api/vX/"
	gmuxApiV := gmuxApi.PathPrefix("/{" + handler.ReqVarMajorVersion + ":" + apiVXRegex + "}/").Subrouter()

	//// Service Routes

	// "/api/vX/anyquiz/"
	gmuxApiV.PathPrefix("/" + serviceAqRest + "/").Handler(hcx.NewServiceProxy(aqRestHostname, aqRestPort))

	//// Gateway routes /api/vX/gateway/
	gmuxApiVGateway := gmuxApiV.PathPrefix("/" + serviceGateway + "/").Subrouter()

	// Health check route
	gmuxApiVGateway.HandleFunc("/"+colHealth, hcx.HealthHandler)

	// Users route
	gmuxApiVGateway.HandleFunc("/"+colUsers, hcx.UsersDefaultHandler)

	// Sessions route

	gmuxApiVGateway.HandleFunc("/"+colSessions, hcx.SessionsDefaultHandler)

	// Users Subroutes
	gmuxApiVGatewayUsers := gmuxApiVGateway.PathPrefix("/" + colUsers + "/").Subrouter()

	// Users Specific routes
	gmuxApiVGatewayUsersSpecific := gmuxApiVGatewayUsers.PathPrefix("/{" + handler.ReqVarUserUuid +
		":" + uuidV4Regex + "}").Subrouter()

	gmuxApiVGatewayUsersSpecific.PathPrefix("").HandlerFunc(hcx.UsersSpecificHandler)

	// Sessions Subroutes
	gmuxApiVGatewaySessions := gmuxApiVGateway.PathPrefix("/" + colSessions + "/").Subrouter()

	// Sessions Specific routes

	// Matches for: /api/vX/gateway/sessions/{sessionIdentifier} which is either "this" or session uuid
	gmuxApiVGatewaySessionsSpecific := gmuxApiVGatewaySessions.PathPrefix(
		"/{" + handler.ReqVarSession + ":(?:" + handler.SpecificSessionHandlerDeleteUserAlias + "|(?:" + uuidV4Regex + "))}").Subrouter()

	gmuxApiVGatewaySessionsSpecific.PathPrefix("").HandlerFunc(hcx.SessionsSpecificHandler)

	// Add Middleware to "/api"
	//gmuxApi.Use
	gmuxApi.Use(handler.NewCors)
	gmuxApi.Use(hcx.NewAuthenticator)
	gmuxApi.Use(hcx.NewRequestLogger)

	// Add Middleware to "/api/{majorVersion}"
	// gmuxApiV.Use

	// Add Middleware to "/api/{majorVersion}/gateway"
	// gmuxApiVGateway.Use
	gmuxApiVGateway.Use(hcx.NewGatewayVersion)
	gmuxApiVGateway.Use(hcx.NewEnsureGatewayVersionSupported)

	// Add Middleware to "/api/{majorVersion}/gateway/users/{uuid}"
	gmuxApiVGatewayUsersSpecific.Use(hcx.NewEnsureAuth)

	// Add Middleware to "/api/{majorVersion}/gateway/sessions/{matchVar}"

	// Add Middleware to "/api/{majorVersion}/anyquiz"

	//Starts listening at the address set, and passes requests at that address
	//to the mux. Exits if ListenAndServerTLS fails
	_ = logger.Log("listenAddress", listenAddr)
	errLS := http.ListenAndServeTLS(listenAddr, tlsCertPath, tlsKeyPath, gmux)
	if errLS != nil {
		_ = logger.Log("http.ListenAndServerTLS", "an error occurred while serving", "error", errLS.Error())
	}
}

func logEnvVar(logger kitlog.Logger, envVar, defaultVal string, required bool) (envVal string, err error) {
	err = nil
	envVal = ""
	if required {
		envVal, err = utility.RequireEnv(envVar)
		if err != nil {
			return
		}
	} else {
		var varSet bool
		envVal, varSet = utility.DefaultEnv(envVar, defaultVal)
		if !varSet {
			_ = logger.Log("notice", "environment variable not set, using default value", "var", envVar, "default", defaultVal)
			return
		}
	}
	_ = logger.Log("logEnvVar", "environment variable set", "var", envVar, "val", envVal)
	return
}

func exitOnEnvError(logger kitlog.Logger, envVar string) (envVal string) {
	val, errL := logEnvVar(logger, envVar, "", true)
	if errL != nil {
		_ = logger.Log("exitOnEnvError", "environment variable must be set", "error", errL, "var", envVar, "result", "exit")
		os.Exit(1)
	}
	return val
}
