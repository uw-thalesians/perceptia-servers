/*
Program gateway is an HTTP API Gateway server for the Perceptia application.

Gateway serves a REST api for the Perceptia application, proxying requests to the backend service responsible for a given collection.
*/
package main

//noinspection SpellCheckingInspection
import (
	"context"
	"encoding/json"

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
const SessionDuration = time.Duration(time.Hour * 24)

// sqlDriverName is the name of the SQL driver to register with the go sql lib
const sqlDriverName = "sqlserver"

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
		logger.Log("error", errTLSCertPath, "result", "exit")
		os.Exit(1)
	}
	tlsKeyPath, errTLSKeyPath := utility.RequireEnv("GATEWAY_TLSKEYPATH")
	// Fail if the path to the key is not provided
	if errTLSKeyPath != nil {
		logger.Log("error", errTLSKeyPath, "result", "exit")
		os.Exit(1)
	}

	mssqlScheme, errMSS := utility.RequireEnv("MSSQL_SCHEME")
	// Fail if the value is not provided
	if errMSS != nil {
		logger.Log("error", errMSS, "result", "exit")
		os.Exit(1)
	}
	mssqlUsername, errMSU := utility.RequireEnv("MSSQL_USERNAME")
	// Fail if the value is not provided
	if errMSU != nil {
		logger.Log("error", errMSU, "result", "exit")
		os.Exit(1)
	}
	mssqlPassword, errMSP := utility.RequireEnv("MSSQL_PASSWORD")
	// Fail if the value is not provided
	if errMSP != nil {
		logger.Log("error", errMSP, "result", "exit")
		os.Exit(1)
	}
	mssqlHost, errMSH := utility.RequireEnv("MSSQL_HOST")
	// Fail if the value is not provided
	if errMSH != nil {
		logger.Log("error", errMSH, "result", "exit")
		os.Exit(1)
	}
	mssqlPort, errMSPO := utility.RequireEnv("MSSQL_PORT")
	// Fail if the value is not provided
	if errMSPO != nil {
		logger.Log("error", errMSPO, "result", "exit")
		os.Exit(1)
	}
	mssqlDatabase, errMSDB := utility.RequireEnv("MSSQL_DATABASE")
	// Fail if the value is not provided
	if errMSDB != nil {
		logger.Log("error", errMSDB, "result", "exit")
		os.Exit(1)
	}

	// Create DSN to use for connection to mssql
	mssqlDsn := utility.BuildDsn(mssqlScheme, mssqlUsername, mssqlPassword, mssqlHost, mssqlPort, mssqlDatabase)

	// Connect to mssql database
	perceptiaDb, errEMSD := utility.Establish(sqlDriverName, mssqlDsn.String(), true)
	if errEMSD == utility.ErrInvalidDsn {
		logger.Log("msg", "provided dsn invalid", "dsn", mssqlDsn.String(), "error", errEMSD, "result", "exit")
		os.Exit(1)
	} else if errEMSD == utility.ErrUnableToPing {
		logger.Log("msg", "unable to ping database", "dsn", mssqlDsn.String(), "error", errEMSD)
	}

	// Periodically check status of mssql database connection
	pingDbCtx := context.TODO()
	go utility.PingDatabase(pingDbCtx, perceptiaDb, time.Second*10, time.Minute)

	// Create new mux router
	gmux := mux.NewRouter()

	// Health check route
	gmux.HandleFunc("/v1/health", func(w http.ResponseWriter, request *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		type healthObj struct {
			Name    string `json:"name"`
			Version string `json:"version"`
			Status  string `json:"status"`
		}

		healthStatus := healthObj{
			Name:    "Perceptia API Health Report",
			Version: "0.1.1",
			Status:  "ready",
		}
		w.WriteHeader(http.StatusOK)
		errWJ := json.NewEncoder(w).Encode(healthStatus)
		if errWJ != nil {
			log.Printf("%s", "Error writing response")
		}
		return
	})

	//Starts listening at the address set, and passes requests at that address
	//to the mux. Exits if ListenAndServerTLS fails
	log.Printf("server is listening at https://%s...", listenAddr)
	log.Fatal(http.ListenAndServeTLS(listenAddr, tlsCertPath, tlsKeyPath, gmux))
}
