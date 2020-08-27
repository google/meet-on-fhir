package server

import (
	"fmt"
	"net/http"
)

const (
	launchPath = "/launch"
)

// Server handles incoming HTTP requests.
type Server struct {
	AuthorizedFhirURL string
	Port              int
}

// Run starts HTTP server
func (s *Server) Run() error {
	if s.AuthorizedFhirURL == "" {
		return fmt.Errorf("AuthorizedFhirURL must be provided")
	}
	http.HandleFunc(launchPath, s.handleLaunch)

	http.ListenAndServe(fmt.Sprintf(":%d", s.Port), http.DefaultServeMux)
	return nil
}

func (s *Server) handleLaunch(w http.ResponseWriter, r *http.Request) {
	iss := r.URL.Query()["iss"]
	if len(iss) == 0 || len(iss[0]) < 1 {
		http.Error(w, "missing iss in URL query parameters", http.StatusUnauthorized)
		return
	}
	if iss[0] != s.AuthorizedFhirURL {
		http.Error(w, fmt.Sprintf("unauthorized iss %s", iss[0]), http.StatusUnauthorized)
		return
	}
	w.WriteHeader(http.StatusOK)
}
