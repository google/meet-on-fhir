package server

import (
	"context"
	"flag"
	"fmt"
	"net/http"

	"github.com/go-session/session"
	"github.com/google/meet-on-fhir/smartonfhir"
	"golang.org/x/oauth2"
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

	authConfig, err := smartonfhir.GetSmartOAuth2Config(fhirURL)
	if err != nil {
		fmt.Fprint(w, err)
		return
	}

	sess, err := session.Start(context.Background(), w, r)
	if err != nil {
		fmt.Fprint(w, err)
		return
	}
	sess.Set("fhirURL", fhirURL)
	sess.Set("launchID", launchID)
	redirectURL := authConfig.AuthCodeURL(sess.SessionID(), oauth2.SetAuthURLParam("aud", fhirURL))
	http.Redirect(w, r, redirectURL, http.StatusFound)
}

func (s *Server) handleFHIRRedirect(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	sess, err := session.Start(ctx, w, r)
	if err != nil {
		fmt.Fprint(w, err)
		return
	}

	fhirURL := getSessionStringOrEmpty(sess, "fhirURL")
	if fhirURL == "" {
		// error
	}
	launchID := getSessionStringOrEmpty(sess, "launchID")
	if launchID == "" {
		// error
	}

	code := getFirstParamOrEmpty(r, codeKey)
	if code == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	state := getFirstParamOrEmpty(r, stateKey)
	if state == "" || state != sess.SessionID() {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	authConfig, err := smartonfhir.GetSmartOAuth2Config(fhirURL)
	if err != nil {
		fmt.Fprint(w, err)
		return
	}
	// TODO: Figure out whether client ID needs to be passed.
	token, err := authConfig.Exchange(ctx, code)
	sess.Set("fhirAccessToken", value)
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
