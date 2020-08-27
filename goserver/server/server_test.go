package server

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/go-session/session"
	"github.com/google/meet-on-fhir/smartonfhir"
)

func setupFHIRServer(authURL, tokenURL string) string {
	fhirServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte(fmt.Sprintf("{\"authorization_endpoint\": \"%s\", \"token_endpoint\": \"%s\"}", authURL, tokenURL)))
	}))
	return fhirServer.URL
}

func TestLaunchHandlerInvalidParameters(t *testing.T) {
	fhirURL := setupFHIRServer("https://auth.com", "https://token.com")
	*authorizedFhirURL = fhirURL
	tests := []struct {
		name, queryParameters string
		expectedHTTPStatus    int
	}{
		{
			name:               "no iss",
			queryParameters:    "launch=123",
			expectedHTTPStatus: http.StatusUnauthorized,
		},
		{
			name:               "empty iss",
			queryParameters:    "iss=\"\"&launch=123",
			expectedHTTPStatus: http.StatusUnauthorized,
		},
		{
			name:               "unauthorized iss",
			queryParameters:    "iss=https://unauthorized.fhir.com&launch=123",
			expectedHTTPStatus: http.StatusUnauthorized,
		},
		{
			name:               "empty launch id",
			queryParameters:    "launch=\"\"",
			expectedHTTPStatus: http.StatusUnauthorized,
		},
		{
			name:               "no launch id",
			queryParameters:    fmt.Sprintf("iss=%s", fhirURL),
			expectedHTTPStatus: http.StatusUnauthorized,
		},
		{
			name:               "with authorized iss and launch id",
			queryParameters:    fmt.Sprintf("iss=%s&launch=123", fhirURL),
			expectedHTTPStatus: http.StatusFound,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s := &Server{}
			req, err := http.NewRequest("GET", "?"+test.queryParameters, nil)
			if err != nil {
				t.Fatal(err)
			}
			rr := httptest.NewRecorder()
			s.handleLaunch(rr, req)
			if status := rr.Code; status != test.expectedHTTPStatus {
				t.Errorf("server.handleLaunch returned wrong status code: got %v want %v",
					status, test.expectedHTTPStatus)
			}
		})
	}
}

func TestLaunchHandler(t *testing.T) {
	fhirURL := setupFHIRServer("https://auth.com", "https://token.com")
	*authorizedFhirURL = fhirURL
	*smartonfhir.FHIRRedirectURL = "https://redirect.com"
	*smartonfhir.FHIRClientID = "fhir_client"

	s := &Server{}
	req, err := http.NewRequest("GET", fmt.Sprintf("?launch=123&iss=%s", fhirURL), nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	s.handleLaunch(rr, req)
	if rr.Code != http.StatusFound {
		t.Errorf("server.handleLaunch returned wrong status code: got %v want %v",
			rr.Code, http.StatusFound)
	}
	if len(req.Cookies()) == 0 {
		t.Errorf("cookies not set in request")
	}
	sess, err := session.Start(req.Context(), nil, req)
	if err != nil {
		t.Fatal(err)
	}
	if v, ok := sess.Get("fhirURL"); !ok || v.(string) != fhirURL {
		t.Errorf("invalid fhirURL in session, got %v, exp %s", v, fhirURL)
	}
	if v, ok := sess.Get("launchID"); !ok || v.(string) != "123" {
		t.Errorf("invalid launchID in session, got %v, exp 123", v)
	}
	redirectURL := rr.Header().Get("Location")
	if !strings.HasPrefix(redirectURL, "https://auth.com") {
		t.Errorf("redirect URL %s does not start with https://auth.com", redirectURL)
	}
	if !strings.Contains(redirectURL, "response_type=code") {
		t.Errorf("redirect URL %s does not contain response_type=code", redirectURL)
	}
	if !strings.Contains(redirectURL, fmt.Sprintf("client_id=%s", *smartonfhir.FHIRClientID)) {
		t.Errorf("redirect URL %s does not contain client_id=%s", redirectURL, *smartonfhir.FHIRClientID)
	}
	if !strings.Contains(redirectURL, fmt.Sprintf("redirect_uri=%s", url.QueryEscape(*smartonfhir.FHIRRedirectURL))) {
		t.Errorf("redirect URL %s does not contain redirect_uri=%s", redirectURL, url.QueryEscape(*smartonfhir.FHIRRedirectURL))
	}
	if !strings.Contains(redirectURL, "launch=123") {
		t.Errorf("redirect URL %s does not contain launch=123", redirectURL)
	}
	if !strings.Contains(redirectURL, "scope=") {
		t.Errorf("redirect URL %s does not contain scope=", redirectURL)
	}
	if !strings.Contains(redirectURL, fmt.Sprintf("state=%s", sess.SessionID())) {
		t.Errorf("redirect URL %s does not contain state=%s", redirectURL, sess.SessionID())
	}
	if !strings.Contains(redirectURL, fmt.Sprintf("aud=%s", url.QueryEscape(fhirURL))) {
		t.Errorf("redirect URL %s does not contain aud=%s", redirectURL, url.QueryEscape(fhirURL))
	}
}
