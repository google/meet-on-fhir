package smartonfhirtest

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

// StartFHIRServer starts and returns a server that handles requests sent on
// the given configPath and returns the given authBaseURL and tokenURL as JSON
// in the response body.
func StartFHIRServer(configPath, authBaseURL, tokenURL string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == configPath {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(fmt.Sprintf("{\"authorization_endpoint\": \"%s\", \"token_endpoint\": \"%s\"}", authBaseURL, tokenURL)))
			return
		}
		w.WriteHeader(http.StatusNotFound)
		return
	}))
}

// StartFHIRTokenServer starts and returns a server that handles token exchange requests and
// returns the given token as JSON in the response body. The handler will return status 400 if
// the request contains form values that do not match the given ones.
func StartFHIRTokenServer(code, redirectURI, clientID, token string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.FormValue("grant_type") != "authorization_code" || r.FormValue("code") != code || r.FormValue("redirect_uri") != redirectURI || r.FormValue("client_id") != clientID {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.Header()["Content-Type"] = []string{"application/json"}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("{\"access_token\": \"%s\"}", token)))
	}))
}

// ValidateAuthURL validate an authURL by making sure its host and query parameters match the
// given ones.
func ValidateAuthURL(t *testing.T, authURL *url.URL, host, clientID, redirectURL, launchID, state, aud string, scopes []string) {
	if authURL.Host != host {
		t.Errorf("host does not match, got %s, expected %s", authURL.Host, host)
	}
	if authURL.Query().Get("response_type") != "code" {
		t.Errorf("response_type does not match, got %s, expected code", authURL.Query().Get("response_type"))
	}
	if authURL.Query().Get("client_id") != clientID {
		t.Errorf("response_type does not match, got %s, expected %s", authURL.Query().Get("client_id"), clientID)
	}
	if authURL.Query().Get("redirect_uri") != redirectURL {
		t.Errorf("redirect_uri does not match, got %s, expected %s", authURL.Query().Get("redirect_uri"), redirectURL)
	}
	if authURL.Query().Get("launch") != launchID {
		t.Errorf("launch does not match, got %s, expected %s", authURL.Query().Get("launch"), launchID)
	}
	if authURL.Query().Get("scope") != strings.Join(scopes, " ") {
		t.Errorf("scope does not match, got %s, expected %s", authURL.Query().Get("scope"), strings.Join(scopes, "+"))
	}
	if authURL.Query().Get("state") != state {
		t.Errorf("state does not match, got %s, expected %s", authURL.Query().Get("state"), state)
	}
	if authURL.Query().Get("aud") != aud {
		t.Errorf("aud does not match, got %s, expected %s", authURL.Query().Get("aud"), aud)
	}
}
