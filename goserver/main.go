package main

import (
	"flag"
	"log"
	"time"

	"github.com/google/meet-on-fhir/smartonfhir"

	"github.com/google/meet-on-fhir/server"
	"github.com/google/meet-on-fhir/session"
	"github.com/google/meet-on-fhir/session/sessiontest"
)

var fhirScopes = []string{"openid", "fhirUser", "profile", "launch", "launch/patient", "launch/encounter"}
var authorizedFhirURL = flag.String("authorized_fhir_url", "", "The FHIR base url that is authorized to launch the telehealth app. The server will not start if not set.")
var fhirClientID = flag.String("fhir_client_id", "", "The client id for FHIR antuenticaion")
var fhirRedirectURL = flag.String("fhir_redirect_url", "", "The redirect URL for FHIR antuenticaion where the user will be redirected to after a successful FHIR authenticaion.")

var sessionDuration = flag.Duration("session_duration", 60*time.Minute, "The max duration of a session")

var httpServerPort = flag.Int("http_server_port", 8080, "The port to start the server on")

func main() {
	flag.Parse()
	// TODO(Issue #20): Use a Cloud SQL based session manage for producation.
	sm := session.NewManager(sessiontest.NewMemoryStore(), *sessionDuration)
	sc := smartonfhir.NewConfig(*fhirClientID, *fhirRedirectURL, fhirScopes)
	server, err := server.NewServer(*authorizedFhirURL, *httpServerPort, sm, sc)
	if err != nil {
		log.Fatal(err)
	}
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
