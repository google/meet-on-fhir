package server

import (
	"context"
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

func TestRunError(t *testing.T) {
	tests := []struct {
		name, authorizedFHIRURL string
		expectedMessage         string
	}{
		{
			name:            "invalid authorized fhir url",
			expectedMessage: authorizedFHIRURLNotProvidedErrorMsg,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s := &Server{AuthorizedFHIRURL: test.authorizedFHIRURL}
			err := s.Run()
			if err == nil {
				t.Fatal("expecting error, but got nil")
			}
			if !strings.Contains(err.Error(), test.expectedMessage) {
				t.Errorf("expecting error message to contain %s, but got %v", test.expectedMessage, err)
				return
			}
		})
	}
}

func TestLaunchHandlerInvalidParameters(t *testing.T) {
	fhirURL := setupFHIRServer("https://auth.com", "https://token.com")
	s := &Server{AuthorizedFHIRURL: fhirURL}
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
	*smartonfhir.FHIRRedirectURL = "https://redirect.com"
	*smartonfhir.FHIRClientID = "fhir_client"

	s := &Server{AuthorizedFHIRURL: fhirURL}
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

func TestHandleFHIRRedirectError(t *testing.T) {
	s := &Server{AuthorizedFHIRURL: "https://fhir.com"}
	tests := []struct {
		name, queryParameters string
		existingSession       map[string]string
		expectedHTTPStatus    int
	}{
		/*{
			name:               "missing fhirURL in session",
			queryParameters:    "",
			expectedHTTPStatus: http.StatusUnauthorized,
		},
		{
			name:               "missing launchID in session",
			existingSession:    map[string]string{"fhirURL": "https://fhir.com"},
			queryParameters:    "",
			expectedHTTPStatus: http.StatusUnauthorized,
		},*/
		{
			name:               "missing code in request",
			existingSession:    map[string]string{"fhirURL": "https://fhir.com", "launchID": "123"},
			queryParameters:    "",
			expectedHTTPStatus: http.StatusUnauthorized,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx := context.Background()
			req, err := http.NewRequest("GET", "?"+test.queryParameters, nil)
			if err != nil {
				t.Fatal(err)
			}
			if test.existingSession != nil {
				existingSess, err := session.Start(ctx, httptest.NewRecorder(), req)
				if err != nil {
					t.Fatal(err)
				}
				for k, v := range test.existingSession {
					existingSess.Set(k, v)
				}
				fmt.Printf("existing session %v", existingSess)

			}

			rr := httptest.NewRecorder()
			s.handleFHIRRedirect(rr, req)
			if status := rr.Code; status != test.expectedHTTPStatus {
				t.Errorf("server.handleFHIRRedirect returned wrong status code: got %v want %v",
					status, test.expectedHTTPStatus)
			}
		})
	}
}
