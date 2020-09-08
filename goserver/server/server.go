package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/meet-on-fhir/session"
	"github.com/google/meet-on-fhir/smartonfhir"
)

const (
	launchPath  = "/launch"
	issKey      = "iss"
	launchIDKey = "launch"
	codeKey     = "code"
	stateKey    = "state"

	authorizedFHIRURLNotProvidedErrorMsg = "AuthorizedFHIRURL must be provided"

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
	sm *session.StoreManager
	// sc is the configuration for SmartOnFhir.
	sc *smartonfhir.Config
}

// NewServer creates and returns a new server.
func NewServer(authorizedFHIRURL string, port int, sm *session.StoreManager, sc *smartonfhir.Config) (*Server, error) {
	if authorizedFHIRURL == "" {
		return nil, fmt.Errorf(authorizedFHIRURLNotProvidedErrorMsg)
	}
	return &Server{authorizedFHIRURL: authorizedFHIRURL, port: port, sm: sm, sc: sc}, nil
}

// Run starts HTTP server
func (s *Server) Run() error {
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

	sess, err := session.New(s.sm, w, r)
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

	sess.Put(fhirURLSessionKey, fhirURL)
	sess.Put(launchIDSessionKey, launchID)
	if err := s.sm.Save(sess); err != nil {
		http.Error(w, "cannot create session", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, redirectURL, http.StatusFound)
}

func (s *Server) handleFHIRRedirect(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	sess, err := session.Find(s.sm, r)
	if err != nil {
		http.Error(w, "invalid session", http.StatusUnauthorized)
		return
	}
	fhirURL := sess.GetStringOrEmpty(fhirURLSessionKey)
	if fhirURL == "" {
		http.Error(w, "invalid session: missing fhirURL", http.StatusUnauthorized)
		return
	}
	launchID := sess.GetStringOrEmpty(launchIDSessionKey)
	if launchID == "" {
		http.Error(w, "invalid session: missing launchID", http.StatusUnauthorized)
		return
	}
	code := getFirstParamOrEmpty(r, codeKey)
	if code == "" {
		http.Error(w, "missing code in URL query parameters", http.StatusBadRequest)
		return
	}
	state := getFirstParamOrEmpty(r, stateKey)
	if state == "" || state != sess.ID {
		http.Error(w, "missing or invalid state", http.StatusBadRequest)
		return
	}
	token, err := s.sc.Exchange(ctx, fhirURL, code)
	if err != nil {
		http.Error(w, "cannot exchange for FHIR access token", http.StatusBadRequest)
	}

	fmt.Fprintf(w, "Successfully exchange for FHIR access token %s", token.AccessToken)
	// TODO: Store token in session.
}

func getFirstParamOrEmpty(r *http.Request, key string) string {
	params := r.URL.Query()[key]
	if len(params) == 0 {
		return ""
	}
	return params[0]
}
