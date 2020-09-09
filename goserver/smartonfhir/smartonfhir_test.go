package smartonfhir

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	testLaunchID        = "123"
	testFHIRAuthURL     = "https://auth.com"
	testFHIRTokenURL    = "https://token.com"
	testFHIRClientID    = "fhir_client"
	testFHIRRedirectURL = "https://redirect.com"
	testScopes          = []string{"launch", "profile"}
	testState           = "test-state"
)

func setupFHIRServer() string {
	fhirServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte(fmt.Sprintf("{\"authorization_endpoint\": \"%s\", \"token_endpoint\": \"%s\"}", testFHIRAuthURL, testFHIRTokenURL)))
	}))
	return fhirServer.URL
}

func TestAuthCodeURL(t *testing.T) {
	fhirURL := setupFHIRServer()
	config := NewConfig(testFHIRClientID, testFHIRRedirectURL, testScopes)
	url := config.AuthCodeURL(fhirURL, testLaunchID, testState)

}
