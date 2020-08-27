package main

import (
	"flag"
	"log"

	"github.com/google/meet-on-fhir/server"
)

var authorizedFhirURL = flag.String("authorized_fhir_url", "", "The FHIR base url that is authorized to launch the telehealth app. The server will not start if not set.")
var httpServerPort = flag.Int("http_server_port", 8080, "The port to start the server on")

func main() {
	flag.Parse()
	server := &server.Server{AuthorizedFhirURL: *authorizedFhirURL, Port: *httpServerPort}
	if err := server.Run(); err != nil {
		log.Fatal(err)
		return
	}
}
