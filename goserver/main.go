package main

import (
	"flag"
	"log"

	"github.com/google/meet-on-fhir/server"
)

var authorizedFhirURL = flag.String("authorized_fhir_url", "", "The FHIR base url that is authorized to launch the telehealth app. If not set, launch endpoint will always return HTTP 401.")
var httpServerPort = flag.Int("http_server_port", 8080, "The port to start the server on")

func main() {
	flag.Parse()
	server := &server.Server{AuthorizedFhirURL: *authorizedFhirURL, Port: *httpServerPort}
	if err := server.Run(); err != nil {
		log.Fatal(err)
		return
	}
}
