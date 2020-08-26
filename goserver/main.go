package main

import (
	"flag"

	"github.com/google/meet-on-fhir/server"
)

var authorizedFhirEndpoint = flag.String("authorized_fhir_endpoint", "", "The FHIR endpoint that is authorized to laucnh the telehealth app")

func main() {
	flag.Parse()
	server := &server.Server{AuthorizedFHIREndPoint: *authorizedFhirEndpoint}
	server.Run()
}
