package server

import (
	"flag"
	"fmt"
	"net/http"
)

const (
	launchPath = "/launch"
)

var authorizedFhirURL = flag.String("authorized_fhir_url", "", "The FHIR base url that is authorized to launch the telehealth app. If not set, launch endpoint will always return HTTP 401.")

// Server handles incoming HTTP requests.
type Server struct {
	Port int
}

// Run starts HTTP server
func (s *Server) Run() {
	http.HandleFunc(launchPath, s.handleLaunch)

	http.ListenAndServe(fmt.Sprintf(":%d", s.Port), http.DefaultServeMux)
}

func (s *Server) handleLaunch(w http.ResponseWriter, r *http.Request) {
	iss := r.URL.Query()["iss"]
	if len(iss) == 0 || len(iss[0]) < 1 || iss[0] != *authorizedFhirURL {
		http.Error(w, "missing iss in URL query parameters", http.StatusUnauthorized)
		return
	}
	w.WriteHeader(http.StatusOK)
}
