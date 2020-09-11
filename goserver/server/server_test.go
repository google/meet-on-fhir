package server

import (
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
			name:               "no iss provided",
			queryParameters:    "",
			expectedHTTPStatus: http.StatusUnauthorized,
		},
		{
			name:               "empty iss",
			queryParameters:    "iss=\"\"",
			expectedHTTPStatus: http.StatusUnauthorized,
		},
		{
			name:               "unauthorized iss",
			queryParameters:    "iss=https://unauthorized.fhir.com",
			expectedHTTPStatus: http.StatusUnauthorized,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s, err := NewServer("https://authorized.fhir.com", 0, nil)
			if err != nil {
				t.Fatalf("NewServer(authorizedFHIRURL, 0, sm) -> %v, nil expected", err)
				return
			}
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

func TestHandleLaunch(t *testing.T) {
	fhirURL := "https://authorized.fhir.com"
	sm := session.NewManager(sessiontest.NewMemoryStore(), 30*time.Minute)
	s, err := NewServer(fhirURL, 0, sm)
	if err != nil {
		t.Fatalf("NewServer(authorizedFHIRURL, 0, sm) -> %v, nil expected", err)
		return
	}
	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/?iss="+fhirURL, nil)
	s.handleLaunch(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("server.handleLaunch returned wrong status code, got %v, want %v",
			status, http.StatusOK)
	}

	sess, err := sm.Retrieve(req)
	if err != nil {
		t.Fatalf("cannot find session, got err %v", err)
	}
	if sess.FHIRURL != fhirURL {
		t.Fatalf("unexpected fhirURL in session: %s, wanted: %s", sess.FHIRURL, fhirURL)
	}
}
