package server

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/google/meet-on-fhir/session"
	"github.com/google/meet-on-fhir/session/sessiontest"
	"github.com/google/meet-on-fhir/smartonfhir"
	"github.com/google/meet-on-fhir/smartonfhir/smartonfhirtest"
)

var (
	testLaunchID        = "123"
	testFHIRAuthURL     = "https://auth.com"
	testFHIRTokenURL    = "https://token.com"
	testFHIRClientID    = "fhir_client"
	testFHIRRedirectURL = "https://redirect.com"
	testScopes          = []string{"launch", "profile"}
)

func defaultServer(fhirURL string) *Server {
	sm := session.NewManager(sessiontest.NewMemoryStore(), func() string { return "test-session-id" }, 30*time.Minute)
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
	sf := smartonfhirtest.StartFHIRServer("/config", testFHIRAuthURL, testFHIRTokenURL)
	fhirURL := sf.URL
	s := defaultServer(fhirURL)
	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", fmt.Sprintf("/?launch=%s&iss=%s", testLaunchID, fhirURL), nil)
	s.handleLaunch(rr, req)
	if rr.Code != http.StatusFound {
		t.Errorf("server.handleLaunch returned wrong status code: got %v want %v",
			rr.Code, http.StatusFound)
	}

	// Make sure session is created and contains expected values.
	sess, err := s.sm.Retrieve(req)
	if err != nil {
		t.Fatalf("cannot find session in either request or store, got err %v", err)
	}
	if sess.FHIRURL != fhirURL {
		t.Errorf("invalid fhirURL in session, got %s, exp %s", sess.FHIRURL, fhirURL)
	}
	if sess.LaunchID != testLaunchID {
		t.Errorf("invalid launchID in session, got %s, exp %s", sess.LaunchID, testLaunchID)
	}

	rawurl := rr.Header().Get("Location")
	authURL, err := url.Parse(rawurl)
	if err != nil {
		t.Fatalf("url.Parse() -> %v, nil expected", err)
	}
	smartonfhirtest.ValidateAuthURL(t, authURL, testFHIRAuthURL, testFHIRClientID, testFHIRRedirectURL, testLaunchID, sess.ID, fhirURL, testScopes)
}

func TestHandleFHIRRedirectError(t *testing.T) {
	fhirURL := "https://fhir.com"
	s := &Server{authorizedFHIRURL: fhirURL}
	tests := []struct {
		name, queryParameters           string
		sessionFHIRURL, sessionLaunchID string
		expectedHTTPStatus              int
	}{
		{
			name:               "missing session",
			queryParameters:    "code=456",
			expectedHTTPStatus: http.StatusUnauthorized,
		},
		{
			name:               "missing fhirURL in session",
			sessionLaunchID:    "123",
			queryParameters:    "code=456",
			expectedHTTPStatus: http.StatusUnauthorized,
		},
		{
			name:               "missing launchID in session",
			sessionFHIRURL:     fhirURL,
			queryParameters:    "code=456",
			expectedHTTPStatus: http.StatusUnauthorized,
		},
		{
			name:               "missing code in request",
			sessionFHIRURL:     fhirURL,
			sessionLaunchID:    "123",
			expectedHTTPStatus: http.StatusBadRequest,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s = defaultServer(fhirURL)
			req, err := http.NewRequest("GET", "?"+test.queryParameters, nil)
			if err != nil {
				t.Fatalf("http.NewRequest() -> %v, nil expected", err)
			}
			if test.sessionFHIRURL != "" || test.sessionLaunchID != "" {
				sess, err := s.sm.New(httptest.NewRecorder(), req)
				if err != nil {
					t.Fatal(err)
				}
				sess.FHIRURL = test.sessionFHIRURL
				sess.LaunchID = test.sessionLaunchID
				if err = s.sm.Save(sess); err != nil {
					t.Fatalf("s.sm.Save() -> %v, nil expected", err)
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
