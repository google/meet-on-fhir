package main

import (
	"flag"
	"log"
	"time"

	"github.com/google/meet-on-fhir/smartonfhir"

	"github.com/google/meet-on-fhir/server"
	"github.com/google/meet-on-fhir/session"
	"github.com/google/meet-on-fhir/session/sessiontest"
	"github.com/google/uuid"
)

var authorizedFhirURL = flag.String("authorized_fhir_url", "", "The FHIR base url that is authorized to launch the telehealth app. The server will not start if not set.")
var sessionDuration = flag.Duration("session_duration", 60*time.Minute, "The max duration of a session")

var httpServerPort = flag.Int("http_server_port", 8080, "The port to start the server on")
var fhirClientID = flag.String("fhir_client_id", "", "Smart on FHIR client id")
var fhirRedirectURL = flag.String("fhir_redirect_url", "", "Smart on FHIR redirect URL")

const (
	fhirScopes = []string{"openid", "fhirUser", "profile", "launch", "launch/patient", "launch/encounter"}
)

func main() {
	flag.Parse()
<<<<<<< HEAD

	sc := smartonfhir.NewConfig(*fhirClientID, *fhirRedirectURL, fhirScopes)
	server, err := server.NewServer(*authorizedFhirURL, *httpServerPort, nil, sc)
=======
	// TODO: Use a session manage for producation use.
	sm := session.NewManager(sessiontest.NewMemoryStore(), func() string { return uuid.New().String() }, *sessionDuration)
	server, err := server.NewServer(*authorizedFhirURL, *httpServerPort, sm)
>>>>>>> server-session
	if err != nil {
		log.Fatal(err)
		return
	}
	if err := server.Run(); err != nil {
		log.Fatal(err)
		return
	}
}
