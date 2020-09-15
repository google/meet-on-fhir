package smartonfhir

import (
	"context"
	"fmt"
	"net/url"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	"golang.org/x/oauth2"

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
	testRefreshToken    = "test-refresh-token"
	testTokenType       = "Bearer"
	testPatientID       = "p123"
	testEncounterID     = "e123"
	testScope           = "launch+profile"
	testTokenJSON       = fmt.Sprintf("{\"access_token\":\"%s\", \"refresh_token\":\"%s\", \"token_type\":\"%s\", \"patient\":\"%s\", \"encounter\":\"%s\", \"scope\":\"%s\"}", testToken, testRefreshToken, testTokenType, testPatientID, testEncounterID, testScope)
)

func TestAuthCodeURL(t *testing.T) {
	s := smartonfhirtest.StartFHIRServer(smartConfigPath, testFHIRAuthURL, testFHIRTokenURL)
	defer s.Close()
	config := NewConfig(testFHIRClientID, testFHIRRedirectURL, testScopes)
	rawurl, err := config.AuthCodeURL(s.URL, testLaunchID, testState)
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
	ts := smartonfhirtest.StartFHIRTokenServer(testAuthCode, testFHIRRedirectURL, testFHIRClientID, []byte(testTokenJSON))
	fs := smartonfhirtest.StartFHIRServer(smartConfigPath, testFHIRAuthURL, ts.URL)
	defer func() {
		ts.Close()
		fs.Close()
	}()

	config := NewConfig(testFHIRClientID, testFHIRRedirectURL, testScopes)
	fc, err := config.Exchange(context.Background(), fs.URL, testAuthCode)
	if err != nil {
		t.Fatalf("config.Exchange() -> %v, nil expected", err)
	}
	expectedContext := &FHIRContext{Token: &oauth2.Token{AccessToken: testToken, RefreshToken: testRefreshToken, TokenType: testTokenType}, EncounterID: testEncounterID, PatientID: testPatientID, Scope: testScope}
	if diff := cmp.Diff(fc, expectedContext, cmpopts.IgnoreUnexported(oauth2.Token{})); diff != "" {
		t.Errorf("got fc does not equal to expected one, diff %s", diff)
	}
}

func TestExchangeError(t *testing.T) {
	ts := smartonfhirtest.StartFHIRTokenServer(testAuthCode, testFHIRRedirectURL, testFHIRClientID, []byte(testTokenJSON))
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
			config := NewConfig(test.clientID, test.redirectURL, test.scopes)
			_, err := config.Exchange(context.Background(), fs.URL, testAuthCode)
			if err == nil {
				t.Fatalf("config.Exchange() -> nil, error expected")
			}
		})
	}
}
