package smartonfhir

import (
	"context"
	"net/url"
	"testing"

	"github.com/google/meet-on-fhir/smartonfhir/smartonfhirtest"
)

var (
	testLaunchID        = "123"
	testFHIRAuthURL     = "https://auth.com"
	testFHIRTokenURL    = "https://token.com"
	testFHIRClientID    = "fhir_client"
	testFHIRRedirectURL = "https://redirect.com"
	testScopes          = []string{"launch", "profile"}
	testState           = "test-state"
	testAuthCode        = "auth-code"
	testToken           = "test-token"
)

func TestAuthCodeURL(t *testing.T) {
	s := smartonfhirtest.StartFHIRServer(smartConfigPath, testFHIRAuthURL, testFHIRTokenURL)
	defer s.Close()
	config := NewConfig(testFHIRClientID, s.URL, testFHIRRedirectURL, testScopes)
	rawurl, err := config.AuthCodeURL(testLaunchID, testState)
	if err != nil {
		t.Fatalf("config.AuthCodeURL() -> %v, nil expected", err)
	}
	actual, err := url.Parse(rawurl)
	if err != nil {
		t.Fatalf("url.Parse() -> %v, nil expected", err)
	}
	expected, err := smartonfhirtest.AuthURL(testFHIRAuthURL, testFHIRClientID, testFHIRRedirectURL, testLaunchID, testState, s.URL, testScopes)
	if err != nil {
		t.Fatalf("smartonfhirtest.AuthURL() -> %v, nil expected", err)
	}
	if diff := smartonfhirtest.DiffAuthURLs(actual, expected); diff != "" {
		t.Errorf("actual authentication URL does not equia to expected, diff %s", diff)
	}
}

func TestExchange(t *testing.T) {
	ts := smartonfhirtest.StartFHIRTokenServer(testAuthCode, testFHIRRedirectURL, testFHIRClientID, testToken)
	fs := smartonfhirtest.StartFHIRServer(smartConfigPath, testFHIRAuthURL, ts.URL)
	defer func() {
		ts.Close()
		fs.Close()
	}()

	config := NewConfig(testFHIRClientID, fs.URL, testFHIRRedirectURL, testScopes)
	token, err := config.Exchange(context.Background(), testAuthCode)
	if err != nil {
		t.Fatalf("config.Exchange() -> %v, nil expected", err)
	}
	if token.AccessToken != testToken {
		t.Errorf("returned token %s does not equal to expected %s", token.AccessToken, testToken)
	}
}

func TestExchangeError(t *testing.T) {
	ts := smartonfhirtest.StartFHIRTokenServer(testAuthCode, testFHIRRedirectURL, testFHIRClientID, testToken)
	fs := smartonfhirtest.StartFHIRServer(smartConfigPath, testFHIRAuthURL, ts.URL)
	defer func() {
		ts.Close()
		fs.Close()
	}()

	tests := []struct {
		name, clientID, redirectURL string
		scopes                      []string
	}{
		{
			name:        "client id does not match",
			clientID:    "wrong id",
			redirectURL: testFHIRRedirectURL,
			scopes:      testScopes,
		},
		{
			name:        "redirectURL does not match",
			clientID:    testFHIRClientID,
			redirectURL: "wrong redirect URL",
			scopes:      testScopes,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			config := NewConfig(test.clientID, fs.URL, test.redirectURL, test.scopes)
			_, err := config.Exchange(context.Background(), testAuthCode)
			if err == nil {
				t.Fatalf("config.Exchange() -> nil, error expected")
			}
		})
	}
}
