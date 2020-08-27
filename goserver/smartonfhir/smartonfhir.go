package smartonfhir

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"

	"golang.org/x/oauth2"
)

const (
	smartConfigPath = "/.well-known/smart-configuration"
	authURLKey      = "authorization_endpoint"
	tokenURLKey     = "token_endpoint"
)

var fhirScopes = []string{"openid", "fhirUser", "profile", "launch", "launch/patient", "launch/encounter"}

var FHIRClientID = flag.String("fhir_client_id", "", "Smart on FHIR client id")
var FHIRRedirectURL = flag.String("fhir_redirect_url", "", "Smart on FHIR redirect URL")

func GetFHIRAuthURL(fhirURL, launchID, state string) (string, error) {
	config, err := gerFHIROAuthConfig(fhirURL)
	if err != nil {
		return "", err
	}

	return config.AuthCodeURL(state, oauth2.SetAuthURLParam("aud", fhirURL), oauth2.SetAuthURLParam("launch", launchID)), nil
}

func GetFHIRAuthToken(ctx context.Context, fhirURL, code string) (*oauth2.Token, error) {
	config, err := gerFHIROAuthConfig(fhirURL)
	if err != nil {
		return nil, err
	}

	return config.Exchange(ctx, code, oauth2.SetAuthURLParam("client_id", *FHIRClientID))
}

func gerFHIROAuthConfig(fhirURL string) (*oauth2.Config, error) {
	resp, err := http.Get(fhirURL + smartConfigPath)
	if err != nil {
		return nil, fmt.Errorf("error")
	}
	var dat map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&dat); err != nil {
		return nil, err
	}

	return &oauth2.Config{ClientID: *FHIRClientID, Endpoint: oauth2.Endpoint{AuthURL: dat[authURLKey].(string), TokenURL: dat[tokenURLKey].(string)}, RedirectURL: *FHIRRedirectURL, Scopes: fhirScopes}, nil
}
