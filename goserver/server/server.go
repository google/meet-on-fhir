package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-session/session"
	"github.com/google/meet-on-fhir/smartonfhir"
	"golang.org/x/oauth2"
)

const (
	launchPath  = "/launch"
	issKey      = "iss"
	launchIDKey = "launch_id"
	codeKey     = "code"
	stateKey    = "state"
)

// Server handles incoming HTTP requests.
type Server struct {
	AuthorizedFHIREndPoint string
}

// Run starts HTTP server on port 8080
func (s *Server) Run() {
	http.HandleFunc(launchPath, s.handleLaunch)

	http.ListenAndServe(":8080", nil)
}

func (s *Server) handleLaunch(w http.ResponseWriter, r *http.Request) {
	fhirURL := getFirstParamOrEmpty(r, issKey)
	if fhirURL == "" != s.AuthorizedFHIREndPoint {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	launchID := getFirstParamOrEmpty(r, launchIDKey)
	if launchID == "" {
		w.WriteHeader(http.StatusUnauthorized)
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
