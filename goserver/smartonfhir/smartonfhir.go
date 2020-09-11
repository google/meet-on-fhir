/*
Package smartonfhir implements functions for SmartOnFhir protocol base on
http://www.hl7.org/fhir/smart-app-launch/.

To support SmartOnFhir using this package, the server should have two HTTP handlers.
Usage example:
func (s *Server) handleLaunch(w http.ResponseWriter, r *http.Request) {
    sc := smartonfhir.NewConfig(*fhirClientID, *fhirURL, *fhirRedirectURL, fhirScopes)
    redirectURL, err := s.sc.AuthCodeURL(launchID, state)
    if err != nil {
        // Handle error
    }
    // Sends the user to be authenticated by the FHIR server after which they will be
    // redirected back to redirectURL
    http.Redirect(w,  r, redirectURL, http.StatusFound)
}

// Handles requests sent to fhirRedirectURL.
func (s *Server) handleRedirect(w http.ResponseWriter, r *http.Request) {
	code := r.Body["code"]
	rs := r.Body["state"]
	if state != rs {
		// Return error since to prevent CSRF attacks.
	}
    sc := smartonfhir.NewConfig(*fhirClientID, *fhirURL, *fhirRedirectURL, fhirScopes)
    token, err := s.sc.Exchange(ctx, code)
    if err != nil {
        // Handle error
    }
    // Store the token for future use.
}
*/
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
	fhirClientID, fhirURL, fhirRedirectURL string
	fhirScopes                             []string
}

// NewConfig creates and returns a new Config.
func NewConfig(fhirClientID, fhirURL, fhirRedirectURL string, fhirScopes []string) *Config {
	return &Config{fhirClientID: fhirClientID, fhirURL: fhirURL, fhirRedirectURL: fhirRedirectURL, fhirScopes: fhirScopes}
}

// AuthCodeURL returns a URL to the FHIR server's consent page that asks for permissions for the
// scopes specified in Config.
// state is a token to protect the user from CSRF attacks and must be provided. Once a request
// is received in fhirRedirectURL, the server should ensure the state in the request always equals
// to the state passed here.
func (c *Config) AuthCodeURL(launchID, state string) (string, error) {
	config, err := c.authConfig()
	if err != nil {
		return "", err
	}

	return config.AuthCodeURL(state, oauth2.SetAuthURLParam("aud", c.fhirURL), oauth2.SetAuthURLParam("launch", launchID)), nil
}

// Exchange exchanges an authorization code for a token.
func (c *Config) Exchange(ctx context.Context, code string) (*oauth2.Token, error) {
	config, err := c.authConfig()
	if err != nil {
		return nil, err
	}

	return config.Exchange(ctx, code, oauth2.SetAuthURLParam("client_id", c.fhirClientID))
}

// authConfig fetches the FHIR authentication configuration and returns an oauth2.Config based on
// the authURL and tokenURL. authConfig will not check supported_scopes in the FHIR authentication
// configuration since it is not a required field.
func (c *Config) authConfig() (*oauth2.Config, error) {
	rURL, err := url.Parse(c.fhirURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse fhirURL %s", c.fhirURL)
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
