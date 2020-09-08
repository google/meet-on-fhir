package server

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/google/meet-on-fhir/session"
	"github.com/google/meet-on-fhir/smartonfhir"
)

var (
	testLaunchID        = "123"
	testFHIRAuthURL     = "https://auth.com"
	testFHIRTokenURL    = "https://token.com"
	testFHIRClientID    = "fhir_client"
	testFHIRRedirectURL = "https://redirect.com"
	testScopes          = []string{"launch", "profile"}
)

func setupFHIRServer(authURL, tokenURL string) string {
	fhirServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte(fmt.Sprintf("{\"authorization_endpoint\": \"%s\", \"token_endpoint\": \"%s\"}", authURL, tokenURL)))
	}))
	return fhirServer.URL
}

func defaultServer(fhirURL string) *Server {
	sm := session.NewStoreManager(session.NewMemoryStore(), func() string { return "test-session-id" })
	sc := smartonfhir.NewConfig(testFHIRClientID, testFHIRRedirectURL, testScopes)
	s, _ := NewServer(fhirURL, 0, sm, sc)
	return s
}

func TestNewServerError(t *testing.T) {
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
			_, err := NewServer(test.authorizedFHIRURL, 0, nil, nil)
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

func TestLaunchHandler_HTTPError(t *testing.T) {
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
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s := defaultServer("https://authorized.fhir.com")
			req := httptest.NewRequest("GET", "/?"+test.queryParameters, nil)
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
	fhirURL := setupFHIRServer(testFHIRAuthURL, testFHIRTokenURL)
	s := defaultServer(fhirURL)
	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", fmt.Sprintf("/?launch=%s&iss=%s", testLaunchID, fhirURL), nil)
	s.handleLaunch(rr, req)
	if rr.Code != http.StatusFound {
		t.Errorf("server.handleLaunch returned wrong status code: got %v want %v",
			rr.Code, http.StatusFound)
	}

	// Make sure session is created and contains expected values.
	sess, err := session.Find(s.sm, req)
	if err != nil {
		t.Fatalf("cannot find session in either request or store, got err %v", err)
	}
	if v := sess.Get(fhirURLSessionKey).(string); v != fhirURL {
		t.Errorf("invalid fhirURL in session, got %v, exp %s", v, fhirURL)
	}
	if v := sess.Get(launchIDSessionKey).(string); v != testLaunchID {
		t.Errorf("invalid launchID in session, got %v, exp %s", v, testLaunchID)
	}

	redirectURL := rr.Header().Get("Location")
	if !strings.HasPrefix(redirectURL, testFHIRAuthURL) {
		t.Errorf("redirect URL %s does not start with %s", redirectURL, testFHIRAuthURL)
	}
	if !strings.Contains(redirectURL, "response_type=code") {
		t.Errorf("redirect URL %s does not contain response_type=code", redirectURL)
	}
	if !strings.Contains(redirectURL, fmt.Sprintf("client_id=%s", testFHIRClientID)) {
		t.Errorf("redirect URL %s does not contain client_id=%s", redirectURL, testFHIRClientID)
	}
	if !strings.Contains(redirectURL, fmt.Sprintf("redirect_uri=%s", url.QueryEscape(testFHIRRedirectURL))) {
		t.Errorf("redirect URL %s does not contain redirect_uri=%s", redirectURL, url.QueryEscape(testFHIRRedirectURL))
	}
	if !strings.Contains(redirectURL, fmt.Sprintf("launch=%s", testLaunchID)) {
		t.Errorf("redirect URL %s does not contain launch=%s", redirectURL, testLaunchID)
	}
	if !strings.Contains(redirectURL, fmt.Sprintf("scope=%s", strings.Join(testScopes, "+"))) {
		t.Errorf("redirect URL %s does not contain scope=%s", redirectURL, strings.Join(testScopes, "+"))
	}
	if !strings.Contains(redirectURL, fmt.Sprintf("state=%s", sess.ID)) {
		t.Errorf("redirect URL %s does not contain state=%s", redirectURL, sess.ID)
	}
	if !strings.Contains(redirectURL, fmt.Sprintf("aud=%s", url.QueryEscape(fhirURL))) {
		t.Errorf("redirect URL %s does not contain aud=%s", redirectURL, url.QueryEscape(fhirURL))
	}
}

func TestHandleFHIRRedirectError(t *testing.T) {
	s := &Server{authorizedFHIRURL: "https://fhir.com"}
	tests := []struct {
		name, queryParameters string
		existingSession       map[string]string
		expectedHTTPStatus    int
	}{
		{
			name:               "missing fhirURL in session",
			queryParameters:    "",
			expectedHTTPStatus: http.StatusUnauthorized,
		},
		{
			name:               "missing launchID in session",
			existingSession:    map[string]string{"fhirURL": "https://fhir.com"},
			queryParameters:    "",
			expectedHTTPStatus: http.StatusUnauthorized,
		},
		{
			name:               "missing code in request",
			existingSession:    map[string]string{fhirURLSessionKey: "https://fhir.com", launchIDSessionKey: "123"},
			queryParameters:    "",
			expectedHTTPStatus: http.StatusUnauthorized,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s = defaultServer("https://fhir.com")
			req, err := http.NewRequest("GET", "?"+test.queryParameters, nil)
			if err != nil {
				t.Fatal(err)
			}
			if test.existingSession != nil {
				existingSess, err := session.New(s.sm, httptest.NewRecorder(), req)
				if err != nil {
					t.Fatal(err)
				}
				for k, v := range test.existingSession {
					existingSess.Put(k, v)
				}
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
