package main

import (
	"flag"
	"log"
	"time"

	"github.com/google/meet-on-fhir/server"
	"github.com/google/meet-on-fhir/session"
	"github.com/google/meet-on-fhir/session/sessiontest"
)

var authorizedFhirURL = flag.String("authorized_fhir_url", "", "The FHIR base url that is authorized to launch the telehealth app. The server will not start if not set.")
var sessionDuration = flag.Duration("session_duration", 60*time.Minute, "The max duration of a session")

var httpServerPort = flag.Int("http_server_port", 8080, "The port to start the server on")

func main() {
	flag.Parse()
	// TODO(Issue #20): Use a Cloud SQL based session manage for producation.
	sm := session.NewManager(sessiontest.NewMemoryStore(), *sessionDuration)
	server, err := server.NewServer(*authorizedFhirURL, *httpServerPort, sm)
	if err != nil {
		log.Fatal(err)
	}
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
