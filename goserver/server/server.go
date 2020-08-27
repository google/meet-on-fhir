package server

import (
	"context"
	"flag"
	"fmt"
	"net/http"

	"github.com/go-session/session"
	"github.com/google/meet-on-fhir/smartonfhir"
)

const (
	launchPath  = "/launch"
	issKey      = "iss"
	launchIDKey = "launch"
	codeKey     = "code"
	stateKey    = "state"
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
	fhirURL := getFirstParamOrEmpty(r, issKey)
	if fhirURL == "" {
		http.Error(w, "missing iss in URL query parameters", http.StatusUnauthorized)
		return
	}
	if fhirURL != *authorizedFhirURL {
		http.Error(w, fmt.Sprintf("unauthorized iss %s", fhirURL), http.StatusUnauthorized)
		return
	}

	launchID := getFirstParamOrEmpty(r, launchIDKey)
	if launchID == "" {
		http.Error(w, "missing launch in URL query parameters", http.StatusUnauthorized)
		return
	}

	sess, err := session.Start(context.Background(), w, r)
	if err != nil {
		http.Error(w, "cannot create session", http.StatusBadRequest)
		return
	}

	// use session ID as the state to prevent CSRF attacks
	redirectURL, err := smartonfhir.GetFHIRAuthURL(fhirURL, launchID, sess.SessionID())
	if err != nil {
		http.Error(w, "cannot get FHIR authentication URL", http.StatusBadRequest)
		return
	}

	sess.Set("fhirURL", fhirURL)
	sess.Set("launchID", launchID)
	if err := sess.Save(); err != nil {
		http.Error(w, "cannot create session", http.StatusBadRequest)
		return
	}
	http.Redirect(w, r, redirectURL, http.StatusFound)
}

func (s *Server) handleFHIRRedirect(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	sess, err := session.Start(ctx, w, r)
	if err != nil {
		http.Error(w, "invalid session", http.StatusUnauthorized)
		return
	}

	fhirURL := getSessionStringOrEmpty(sess, "fhirURL")
	if fhirURL == "" {
		http.Error(w, "invalid session: missing fhirURL", http.StatusUnauthorized)
		return
	}
	launchID := getSessionStringOrEmpty(sess, "launchID")
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
	if state == "" || state != sess.SessionID() {
		http.Error(w, "missing or invalid state", http.StatusBadRequest)
		return
	}

	token, err := smartonfhir.GetFHIRAuthToken(ctx, fhirURL, code)
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

func getSessionStringOrEmpty(sess session.Store, key string) string {
	v, exists := sess.Get(key)
	if !exists {
		return ""
	}
	return v.(string)
}
