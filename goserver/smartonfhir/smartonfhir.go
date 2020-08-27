package smartonfhir

import (
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

var fhirClientID = flag.String("fhir_client_id", "", "Smart on FHIR client id")
var fhirScopes = flag.String("fhir_scopes", "", "Smart on FHIR scopes")
var fhirRedirectURL = flag.String("fhir_redirect_url", "", "Smart on FHIR redirect URL")

type SmartConfig struct {
	authURL, tokenURL string
}

func GetSmartOAuth2Config(fhirURL string) (*oauth2.Config, error) {
	sc, err := gerSmartConfiguration(fhirURL)
	if err != nil {
		return nil, err
	}

	return &oauth2.Config{ClientID: *fhirClientID, Endpoint: oauth2.Endpoint{AuthURL: sc.authURL, TokenURL: sc.tokenURL}, RedirectURL: *fhirRedirectURL}, nil
}

func gerSmartConfiguration(fhirURL string) (*SmartConfig, error) {
	resp, err := http.Get(fhirURL + smartConfigPath)
	if err != nil {
		return nil, fmt.Errorf("error")
	}
	var dat map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&dat); err != nil {
		return nil, err
	}

	return &SmartConfig{authURL: dat[authURLKey].(string), tokenURL: dat[tokenURLKey].(string)}, nil
}
