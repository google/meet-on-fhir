package server

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/google/meet-on-fhir/session"
	"github.com/google/meet-on-fhir/session/sessiontest"
)

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
			_, err := NewServer(test.authorizedFHIRURL, 0, nil)
			if err == nil {
				t.Fatal("expecting error, but got nil")
			}
			if !strings.Contains(err.Error(), test.expectedMessage) {
				t.Errorf("expecting error message to contain %s, but got %v", test.expectedMessage, err)
			}
		})
	}
}

func TestLaunchHandlerError(t *testing.T) {
	tests := []struct {
		name, queryParameters string
		store                 session.Store
		expectedHTTPStatus    int
	}{
		{
			name:               "no iss provided",
			queryParameters:    "",
			store:              nil,
			expectedHTTPStatus: http.StatusUnauthorized,
		},
		{
			name:               "empty iss",
			queryParameters:    "?iss=\"\"",
			store:              nil,
			expectedHTTPStatus: http.StatusUnauthorized,
		},
		{
			name:               "unauthorized iss",
			queryParameters:    "?iss=https://unauthorized.fhir.com",
			store:              nil,
			expectedHTTPStatus: http.StatusUnauthorized,
		},
		{
			name:               "new session error",
			queryParameters:    "?iss=https://authorized.fhir.com",
			store:              sessiontest.NewMemoryStoreWithError(fmt.Errorf("new session error"), nil),
			expectedHTTPStatus: http.StatusInternalServerError,
		},
		{
			name:               "save session error",
			queryParameters:    "?iss=https://authorized.fhir.com",
			store:              sessiontest.NewMemoryStoreWithError(nil, fmt.Errorf("save session error")),
			expectedHTTPStatus: http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			sm := session.NewManager(test.store, 30*time.Minute)
			s, err := NewServer("https://authorized.fhir.com", 0, sm)
			if err != nil {
				t.Fatalf("NewServer(authorizedFHIRURL, 0, sm) -> %v, nil expected", err)
			}
			ts := httptest.NewServer(http.HandlerFunc(s.handleLaunch))
			defer ts.Close()
			res, err := http.Get(ts.URL + test.queryParameters)
			if err != nil {
				t.Fatalf("http.Get() -> %v, nil expected", err)
			}
			if status := res.StatusCode; status != test.expectedHTTPStatus {
				t.Errorf("server.handleLaunch returned wrong status code: got %v want %v",
					status, test.expectedHTTPStatus)
			}
		})
	}
}

func TestHandleLaunch(t *testing.T) {
	fhirURL := "https://authorized.fhir.com"
	ss := sessiontest.NewMemoryStore()
	sm := session.NewManager(ss, 30*time.Minute)
	s, err := NewServer(fhirURL, 0, sm)
	if err != nil {
		t.Fatalf("NewServer(authorizedFHIRURL, 0, sm) -> %v, nil expected", err)
	}
	ts := httptest.NewServer(http.HandlerFunc(s.handleLaunch))
	defer ts.Close()
	res, err := http.Get(ts.URL + "?iss=" + fhirURL)
	if err != nil {
		t.Fatalf("http.Get() -> %v, nil expected", err)
	}
	if status := res.StatusCode; status != http.StatusOK {
		t.Errorf("server.handleLaunch returned wrong status code, got %v, want %v",
			status, http.StatusOK)
	}

	sess := sessionFromResp(t, ss, res)
	if sess.FHIRURL != fhirURL {
		t.Fatalf("unexpected fhirURL in session: %s, wanted: %s", sess.FHIRURL, fhirURL)
	}
}

func sessionFromResp(t *testing.T, ss session.Store, res *http.Response) *session.Session {
	sessionID := strings.Split(res.Header.Get("Set-Cookie"), ";")[0]
	sessionID = strings.Split(sessionID, "=")[1]
	b, err := ss.Retrieve(sessionID)
	if err != nil {
		t.Fatalf("ss.Retrieve(%s) -> %v, nil expected", sessionID, err)
	}
	session, err := session.FromBytes(b)
	if err != nil {
		t.Fatalf("session.FromBytes() -> %v, nil expected", err)
	}
	return session
}
