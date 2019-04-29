/*
Program gateway is an HTTP API Gateway server for the Perceptia application.
*/
package main

import (
	"io"
	"log"
	"net/http"
	"time"

	"github.com/uw-thalesians/perceptia-servers/gateway/gateway/utility"

	"github.com/gorilla/mux"
)

// SessionDuration is the time a session is valid. If it has been longer than time.Duration since
// the session was started the session should be treated as no longer valid.
const SessionDuration = time.Duration(time.Hour * 24)

const sqlDriverName = "sqlserver"

func main() {
	// Get address for server to listen for requests on
	listenAddr := utility.DefaultEnv("GATEWAY_LISTEN_ADDR", ":443")

	// Get the directory path to the TLS key and cert
	tlsCertPath, errTLSCertPath := utility.RequireEnv("GATEWAY_TLSCERTPATH")
	// Fail if the path to the cert is not provided
	if errTLSCertPath != nil {
		log.Fatal(errTLSCertPath)
	}
	tlsKeyPath, errTLSKeyPath := utility.RequireEnv("GATEWAY_TLSKEYPATH")
	// Fail if the path to the key is not provided
	if errTLSKeyPath != nil {
		log.Fatal(errTLSKeyPath)
	}

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
	mssqlDb, errEMSD := utility.Establish(sqlDriverName, mssqlDsn.String(), true)
	if errEMSD == utility.ErrInvalidDsn {
		log.Fatalf("The provided DSN: %s, was invalid: %s", mssqlDsn.String(), errEMSD)
	}

	// Create new mux router
	gmux := mux.NewRouter()

	// Health check route
	gmux.HandleFunc("/v1/health", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "text/plain")
		writer.WriteHeader(http.StatusOK)
		_, err := io.WriteString(writer, "Gateway is working as expected.")
		if err != nil {
			log.Printf("%s", "Error writing response")
		}
	})

	//Starts listening at the address set, and passes requests at that address
	//to the mux. Exits if ListenAndServerTLS fails
	log.Printf("server is listening at https://%s...", listenAddr)
	log.Fatal(http.ListenAndServeTLS(listenAddr, tlsCertPath, tlsKeyPath, gmux))
}
