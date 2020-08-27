package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLaunchHandlerCheckISSAuthorization(t *testing.T) {
	s := &Server{AuthorizedFhirURL: "https://authorized.fhir.com"}
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
		{
			name:               "authorized iss",
			queryParameters:    "iss=https://authorized.fhir.com",
			expectedHTTPStatus: http.StatusOK,
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
