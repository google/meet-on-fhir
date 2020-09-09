package smartonfhir

import (
	"context"
	"net/url"
	"testing"

	"github.com/google/meet-on-fhir/smartonfhir/smartonfhirtest"
)

var (
	testLaunchID            = "123"
	testFHIRAuthURL         = "https://auth.com"
	testFHIRAuthURLNoSchema = "auth.com"
	testFHIRTokenURL        = "https://token.com"
	testFHIRClientID        = "fhir_client"
	testFHIRRedirectURL     = "https://redirect.com"
	testScopes              = []string{"launch", "profile"}
	testState               = "test-state"
	testAuthCode            = "auth-code"
	testToken               = "test-token"
)

func TestAuthCodeURL(t *testing.T) {
	s := smartonfhirtest.StartFHIRServer(smartConfigPath, testFHIRAuthURL, testFHIRTokenURL)
	defer s.Close()
	config := NewConfig(testFHIRClientID, testFHIRRedirectURL, testScopes)
	rawurl, err := config.AuthCodeURL(s.URL, testLaunchID, testState)
	if err != nil {
		t.Fatalf("config.AuthCodeURL() -> %v, nil expected", err)
	}
	url, err := url.Parse(rawurl)
	if err != nil {
		t.Fatalf("url.Parse() -> %v, nil expected", err)
	}
	smartonfhirtest.ValidateAuthURL(t, url, testFHIRAuthURLNoSchema, testFHIRClientID, testFHIRRedirectURL, testLaunchID, testState, s.URL, testScopes)
}

func TestExchange(t *testing.T) {
	ts := smartonfhirtest.StartFHIRTokenServer(testAuthCode, testFHIRRedirectURL, testFHIRClientID, testToken)
	fs := smartonfhirtest.StartFHIRServer(smartConfigPath, testFHIRAuthURL, ts.URL)
	defer func() {
		ts.Close()
		fs.Close()
	}()

	config := NewConfig(testFHIRClientID, testFHIRRedirectURL, testScopes)
	token, err := config.Exchange(context.Background(), fs.URL, testAuthCode)
	if err != nil {
		t.Fatalf("config.Exchange() -> %v, nil expected", err)
	}
	if token.AccessToken != testToken {
		t.Errorf("returned token %s does not equal to expected %s", token.AccessToken, testToken)
	}
}
