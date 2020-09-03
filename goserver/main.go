package main

import (
	"flag"
	"log"

	"github.com/google/meet-on-fhir/server"
	"github.com/google/meet-on-fhir/session"
)

var authorizedFhirURL = flag.String("authorized_fhir_url", "", "The FHIR base url that is authorized to launch the telehealth app. The server will not start if not set.")
var httpServerPort = flag.Int("http_server_port", 8080, "The port to start the server on")

func main() {
	flag.StringVar(&session.SessionCookieSecret, "session_cookie_secret", "", "secret key used to encrypt the session cookie")
	flag.Parse()

	server, err := server.NewServer(*authorizedFhirURL, *httpServerPort, nil)
	if err != nil {
		log.Fatal(err)
		return
	}
	if err := server.Run(); err != nil {
		log.Fatal(err)
		return
	}
}
