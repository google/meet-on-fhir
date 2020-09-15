// Package smartonfhirtest provides untilities for smartonfhir testing.
package smartonfhirtest

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"

	"github.com/google/go-cmp/cmp"
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
// returns the given token JSON in the response body. The handler will return status 400 if
// the request contains form values that do not match the given ones.
func StartFHIRTokenServer(code, redirectURI, clientID string, tokenJSON []byte) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		if r.FormValue("grant_type") != "authorization_code" || r.FormValue("code") != code || r.FormValue("redirect_uri") != redirectURI || r.FormValue("client_id") != clientID {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.Header()["Content-Type"] = []string{"application/json"}
		w.WriteHeader(http.StatusOK)
		w.Write(tokenJSON)
	}))
}

// AuthURL returns a FHIR authentication URL with all required query parameters.
func AuthURL(host, clientID, redirectURI, launch, state, aud string, scopes []string) (*url.URL, error) {
	aURL, err := url.Parse(host)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL from host %s", host)
	}
	aURL.RawQuery = fmt.Sprintf("response_type=code&client_id=%s&redirect_uri=%s&launch=%s&state=%s&aud=%s&scope=%s", clientID, redirectURI, launch, state, aud, strings.Join(scopes, "+"))
	return aURL, nil
}

// DiffAuthURLs compares two URLs and return their diff. Returns empty string if equivalent.
func DiffAuthURLs(actual, expected *url.URL) string {
	op := cmp.Comparer(func(a, b *url.URL) bool {
		if a.Scheme != b.Scheme {
			return false
		}
		if a.Host != b.Host {
			return false
		}
		if !cmp.Equal(a.Query(), b.Query()) {
			return false
		}
		return true
	})
	return cmp.Diff(actual, expected, op)
}
