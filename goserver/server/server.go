// Package server implements an HTTP server.
// TODO(Issue #22): Handle errors consistently.
package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/meet-on-fhir/session"
	"github.com/google/meet-on-fhir/smartonfhir"
)

const (
	launchPath       = "/launch"
	fhirRedirectPath = "/fhir_redirect"
	issKey           = "iss"
	launchIDKey      = "launch"
	codeKey          = "code"
	stateKey         = "state"

	authorizedFHIRURLNotProvidedErr    = "AuthorizedFHIRURL must be provided"
	cookieSessionNotFoundErr           = "cookie session not found"
	sessionMissingFHIRURLErr           = "invalid session: missing fhirURL"
	sessionMissingLaunchIDErr          = "invalid session: missing launchID"
	requestQueryMissingCodeErr         = "missing code in URL query parameters"
	requestQueryInvalidStateErr        = "invalid state in query parameters"
	severFailedForFHIRTokenExchangeErr = "server failed to exchange for FHIR access token"
	serverRetrieveSessionErr           = "server failed to retrieve session"
	serverSaveSessionErr               = "server failed to save session"

	fhirURLSessionKey  = "fhirUrl"
	launchIDSessionKey = "launchId"
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
	// sc is the configuration for SmartOnFhir.
	sc *smartonfhir.Config
}

// NewServer creates and returns a new server.
func NewServer(authorizedFHIRURL string, port int, sm *session.Manager, sc *smartonfhir.Config) (*Server, error) {
	if authorizedFHIRURL == "" {
		return nil, fmt.Errorf(authorizedFHIRURLNotProvidedErr)
	}
	return &Server{authorizedFHIRURL: authorizedFHIRURL, port: port, sm: sm, sc: sc}, nil
}

// Run starts HTTP server
func (s *Server) Run() error {
	http.HandleFunc(launchPath, s.handleLaunch)
	http.HandleFunc(launchPath, s.handleLaunch)

	http.ListenAndServe(fmt.Sprintf(":%d", s.port), http.DefaultServeMux)
	return nil
}

func (s *Server) handleLaunch(w http.ResponseWriter, r *http.Request) {
	fhirURL := getFirstParamOrEmpty(r, issKey)
	if fhirURL == "" {
		http.Error(w, "missing iss in URL query parameters", http.StatusUnauthorized)
		return
	}
	if fhirURL != s.authorizedFHIRURL {
		http.Error(w, fmt.Sprintf("unauthorized iss %s", fhirURL), http.StatusUnauthorized)
		return
	}

	launchID := getFirstParamOrEmpty(r, launchIDKey)
	if launchID == "" {
		http.Error(w, "missing launch in URL query parameters", http.StatusUnauthorized)
		return
	}

	sess, err := s.sm.New(w, r)
	if err != nil {
		http.Error(w, "cannot create session", http.StatusInternalServerError)
		return
	}

	// use session ID as the state to prevent CSRF attacks
	redirectURL, err := s.sc.AuthCodeURL(fhirURL, launchID, sess.ID)
	if err != nil {
		http.Error(w, "cannot get FHIR authentication URL", http.StatusBadRequest)
		return
	}

	sess.FHIRURL = fhirURL
	sess.LaunchID = launchID
	if err := s.sm.Save(sess); err != nil {
		http.Error(w, "cannot create session", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, redirectURL, http.StatusFound)
}

func (s *Server) handleFHIRRedirect(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	sess, err := s.sm.Retrieve(r)
	if err == session.ErrNotFound || err == http.ErrNoCookie {
		http.Error(w, cookieSessionNotFoundErr, http.StatusUnauthorized)
		return
	}
	if err != nil {
		http.Error(w, serverRetrieveSessionErr, http.StatusInternalServerError)
		return
	}
	fhirURL := sess.FHIRURL
	if fhirURL == "" {
		http.Error(w, sessionMissingFHIRURLErr, http.StatusUnauthorized)
		return
	}
	launchID := sess.LaunchID
	if launchID == "" {
		http.Error(w, sessionMissingLaunchIDErr, http.StatusUnauthorized)
		return
	}
	code := getFirstParamOrEmpty(r, codeKey)
	if code == "" {
		http.Error(w, requestQueryMissingCodeErr, http.StatusBadRequest)
		return
	}
	state := getFirstParamOrEmpty(r, stateKey)
	if state != sess.ID {
		http.Error(w, requestQueryInvalidStateErr, http.StatusBadRequest)
		return
	}

	fhirContext, err := s.sc.Exchange(ctx, fhirURL, code)
	if err != nil {
		http.Error(w, severFailedForFHIRTokenExchangeErr, http.StatusInternalServerError)
		return
	}

	sess.FHIRContext = fhirContext
	if err := s.sm.Save(sess); err != nil {
		http.Error(w, serverSaveSessionErr, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Successfully authenticated with FHIR")
	// TODO (Issue #24): Return FE contents.
}

func getFirstParamOrEmpty(r *http.Request, key string) string {
	params := r.URL.Query()[key]
	if len(params) == 0 {
		return ""
	}
	return params[0]
}
