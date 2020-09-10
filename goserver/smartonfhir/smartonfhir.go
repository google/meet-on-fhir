// Package smartonfhir provides functions for SmartOnFhir flow base on
// http://www.hl7.org/fhir/smart-app-launch/.
package smartonfhir

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"golang.org/x/oauth2"
)

const (
	smartConfigPath = "/.well-known/smart-configuration"
	authURLKey      = "authorization_endpoint"
	tokenURLKey     = "token_endpoint"
)

// Config contains configuration information for smartonfhir authentication flow.
type Config struct {
	fhirClientID, fhirRedirectURL string
	fhirScopes                    []string
}

// NewConfig creates and returns a new Config.
func NewConfig(fhirClientID, fhirRedirectURL string, fhirScopes []string) *Config {
	return &Config{fhirClientID: fhirClientID, fhirRedirectURL: fhirRedirectURL, fhirScopes: fhirScopes}
}

// AuthCodeURL returns a URL to the FHIR server's consent page
// that asks for permissions for the scopes specified in Config.
// State is a token to protect the user from CSRF attacks and must
// be provided.
func (c *Config) AuthCodeURL(fhirURL, launchID, state string) (string, error) {
	config, err := c.authConfig(fhirURL)
	if err != nil {
		return "", err
	}

	return config.AuthCodeURL(state, oauth2.SetAuthURLParam("aud", fhirURL), oauth2.SetAuthURLParam("launch", launchID)), nil
}

// Exchange exchanges an authorization code for a token.
func (c *Config) Exchange(ctx context.Context, fhirURL, code string) (*oauth2.Token, error) {
	config, err := c.authConfig(fhirURL)
	if err != nil {
		return nil, err
	}

	return config.Exchange(ctx, code, oauth2.SetAuthURLParam("client_id", c.fhirClientID))
}

func (c *Config) authConfig(fhirURL string) (*oauth2.Config, error) {
	rURL, err := url.Parse(fhirURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse fhirURL %s", fhirURL)
	}
	rURL.Path = smartConfigPath
	resp, err := http.Get(rURL.String())
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("got response code %d when fetching smart configuration", resp.StatusCode)
	}
	var dat map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&dat); err != nil {
		return nil, err
	}

	authURL, ok := dat[authURLKey].(string)
	if !ok {
		return nil, fmt.Errorf("no authorization_endpoint found in smart configuration")
	}
	tokenURL, ok := dat[tokenURLKey].(string)
	if !ok {
		return nil, fmt.Errorf("no token_endpoint found in smart configuration")
	}
	return &oauth2.Config{ClientID: c.fhirClientID, Endpoint: oauth2.Endpoint{AuthURL: authURL, TokenURL: tokenURL}, RedirectURL: c.fhirRedirectURL, Scopes: c.fhirScopes}, nil
}
