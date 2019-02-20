/*
 * Program gateway is an HTTP API Gateway server for the Perceptia application.
 */
package main

import (
	"fmt"
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

func main() {
	fmt.Print("Hello World!")

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

	// Create new mux router
	gmux := mux.NewRouter()

	// Fake testHealth route
	gmux.HandleFunc("/v1/testHealth", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "text/plain")
		writer.WriteHeader(http.StatusOK)
		_, err := io.WriteString(writer, "Gateway is working")
		if err != nil {
			log.Printf("%s", "Error writing response")
		}
	})

	log.Printf("tlscert: %s\ntlskey: %s", tlsCertPath, tlsKeyPath)

	//Starts listening at the address set, and passes requests at that address
	//to the mux. Exits if ListenAndServerTLS fails
	log.Printf("server is listening at https://%s...", listenAddr)
	log.Fatal(http.ListenAndServeTLS(listenAddr, tlsCertPath, tlsKeyPath, gmux))
}
