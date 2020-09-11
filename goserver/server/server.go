// Package server implements an HTTP server.
// TODO(Issue #22): Handle errors consistently.
package server

import (
	"fmt"
	"net/http"

	"github.com/google/meet-on-fhir/session"
)

const (
	launchPath = "/launch"

	authorizedFHIRURLNotProvidedErrorMsg = "AuthorizedFHIRURL must be provided"

	fhirURLSessionKey = "fhirUrl"
)

// Server handles incoming HTTP requests.
type Server struct {
	// authorizedFHIRURL is the FHIR URL authorized to launch this app. The value will be validated
	// by launch endpoint to match the iss passed as the query parameter.
	authorizedFHIRURL string
	// The port the HTTP server runs on.
	port int
	// sm is the session manager of the server.
	sm *session.Manager
}

// NewServer creates and returns a new server.
func NewServer(authorizedFHIRURL string, port int, sm *session.Manager) (*Server, error) {
	if authorizedFHIRURL == "" {
		return nil, fmt.Errorf(authorizedFHIRURLNotProvidedErrorMsg)
	}
	return &Server{authorizedFHIRURL: authorizedFHIRURL, port: port, sm: sm}, nil
}

// Run starts HTTP server
func (s *Server) Run() error {
	http.HandleFunc(launchPath, s.handleLaunch)

	http.ListenAndServe(fmt.Sprintf(":%d", s.port), http.DefaultServeMux)
	return nil
}

func (s *Server) handleLaunch(w http.ResponseWriter, r *http.Request) {
	iss := r.URL.Query()["iss"]
	if len(iss) == 0 || len(iss[0]) < 1 {
		http.Error(w, "missing iss in URL query parameters", http.StatusUnauthorized)
		return
	}
	if iss[0] != s.authorizedFHIRURL {
		http.Error(w, fmt.Sprintf("unauthorized iss %s", iss[0]), http.StatusUnauthorized)
		return
	}

	sess, err := s.sm.New(w, r)
	if err != nil {
		http.Error(w, "cannot create session", http.StatusInternalServerError)
		return
	}
	sess.FHIRURL = iss[0]
	s.sm.Save(sess)

	w.WriteHeader(http.StatusOK)
}
