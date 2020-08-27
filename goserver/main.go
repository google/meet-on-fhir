package main

import (
	"flag"

	"github.com/google/meet-on-fhir/server"
)

var httpServerPort = flag.Int("http_server_port", 8080, "The port to start the server on")

func main() {
	flag.Parse()
	server := &server.Server{Port: *httpServerPort}
	server.Run()
}
