package server

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/meet-on-fhir/session"
)

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
			s := &Server{authorizedFHIRURL: test.authorizedFHIRURL}
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

func TestLaunchHandler_HTTPError(t *testing.T) {
	s := &Server{authorizedFHIRURL: "https://authorized.fhir.com"}
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
	sm := session.NewInMemorySessionManager()
	s := &Server{authorizedFHIRURL: fhirURL, sm: sm}
	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/?iss="+fhirURL, nil)
	s.handleLaunch(httptest.NewRecorder(), req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("server.handleLaunch returned wrong status code, got %v, want %v",
			status, http.StatusOK)
	}

	sess, err := session.Find(sm, req)
	if err != nil {
		t.Fatal("cannot find session")
	}
	if sess.FHIRURL != fhirURL {
		t.Fatalf("unexpected FHIRURL session %s, wanted %s", sess.FHIRURL, fhirURL)
	}
}
