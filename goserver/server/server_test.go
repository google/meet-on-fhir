package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"

	"golang.org/x/oauth2"

	"github.com/google/meet-on-fhir/smartonfhir"
	"github.com/google/meet-on-fhir/smartonfhir/smartonfhirtest"

	"github.com/google/meet-on-fhir/session"
	"github.com/google/meet-on-fhir/session/sessiontest"
)

var (
	testLaunchID         = "123"
	testFHIRAuthURL      = "https://auth.com"
	testFHIRTokenURL     = "https://token.com"
	testFHIRClientID     = "fhir_client"
	testFHIRRedirectURL  = "https://redirect.com"
	testScopes           = []string{"launch", "profile"}
	testAuthCode         = "abc"
	testFHIRAccessToken  = "fhir-access-token"
	testFHIRRefreshToken = "fhir-refresh-token"
	testFHIRTokenType    = "Bearer"
	testFHIRTokenJSON    = fmt.Sprintf("{\"access_token\": \"%s\", \"refresh_token\": \"%s\", \"token_type\":\"%s\", \"patient\":\"p123\", \"encounter\": \"e123\"}", testFHIRAccessToken, testFHIRRefreshToken, testFHIRTokenType)
)

func setupBackends() string {
	fts := smartonfhirtest.StartFHIRTokenServer(testAuthCode, testFHIRRedirectURL, testFHIRClientID, []byte(testFHIRTokenJSON))
	sf := smartonfhirtest.StartFHIRServer("/.well-known/smart-configuration", testFHIRAuthURL, fts.URL)
	return sf.URL
}

func defaultServer(fhirURL string, ss session.Store) *Server {
	sm := session.NewManager(ss, 30*time.Minute)
	sc := smartonfhir.NewConfig(testFHIRClientID, testFHIRRedirectURL, testScopes)
	s, _ := NewServer(fhirURL, 0, sm, sc)
	return s
}

func TestNewServerError(t *testing.T) {
	tests := []struct {
		name, authorizedFHIRURL string
		expectedMessage         string
	}{
		{
			name:            "invalid authorized fhir url",
			expectedMessage: authorizedFHIRURLNotProvidedErr,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := NewServer(test.authorizedFHIRURL, 0, nil, nil)
			if err == nil {
				t.Fatal("expecting error, but got nil")
			}
			if !strings.Contains(err.Error(), test.expectedMessage) {
				t.Errorf("expecting error message to contain %s, but got %v", test.expectedMessage, err)
			}
		})
	}
}

func TestLaunchHandlerError(t *testing.T) {
	tests := []struct {
		name, queryParameters string
		store                 session.Store
		expectedHTTPStatus    int
	}{
		{
			name:               "no iss provided",
			queryParameters:    "",
			store:              nil,
			expectedHTTPStatus: http.StatusUnauthorized,
		},
		{
			name:               "empty iss",
			queryParameters:    "?iss=\"\"",
			store:              nil,
			expectedHTTPStatus: http.StatusUnauthorized,
		},
		{
			name:               "unauthorized iss",
			queryParameters:    "?iss=https://unauthorized.fhir.com",
			store:              nil,
			expectedHTTPStatus: http.StatusUnauthorized,
		},
		{
			name:               "new session error",
			queryParameters:    "?iss=https://authorized.fhir.com",
			store:              sessiontest.NewMemoryStore().WithNextStoreErr(fmt.Errorf("new session error")),
			expectedHTTPStatus: http.StatusInternalServerError,
		},
		{
			name:               "save session error",
			queryParameters:    "?iss=https://authorized.fhir.com",
			store:              sessiontest.NewMemoryStore().WithNextStoreExistingErr(fmt.Errorf("save session error")),
			expectedHTTPStatus: http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fhirURL := setupBackends()
			s := defaultServer(fhirURL, sessiontest.NewMemoryStore())
			ts := httptest.NewServer(http.HandlerFunc(s.handleLaunch))
			defer ts.Close()
			res, err := http.Get(ts.URL + test.queryParameters)
			if err != nil {
				t.Fatalf("http.Get() -> %v, nil expected", err)
			}
			if status := res.StatusCode; status != test.expectedHTTPStatus {
				t.Errorf("server.handleLaunch returned wrong status code: got %v want %v",
					status, test.expectedHTTPStatus)
			}
		})
	}
}

func TestHandleLaunch(t *testing.T) {
	fhirURL := setupBackends()
	ss := sessiontest.NewMemoryStore()
	s := defaultServer(fhirURL, ss)
	ts := httptest.NewServer(http.HandlerFunc(s.handleLaunch))
	defer ts.Close()

	hc := http.DefaultClient
	var actualRedirectURL *url.URL
	checkRedirectErr := fmt.Errorf("Redirect response received")
	// Prevent HTTP client from issuing a following request on redirect.
	hc.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		actualRedirectURL = req.URL
		return checkRedirectErr
	}
	res, err := http.Get(fmt.Sprintf("%s?iss=%s&launch=%s", ts.URL, fhirURL, testLaunchID))
	if err == nil {
		t.Fatalf("http.Get() -> nil, %v expected", checkRedirectErr)
	}

	sess := sessionFromResp(t, ss, res)
	if sess.FHIRURL != fhirURL {
		t.Fatalf("unexpected fhirURL in session: %s, wanted: %s", sess.FHIRURL, fhirURL)
	}
	if sess.LaunchID != testLaunchID {
		t.Errorf("invalid launchID in session, got %s, exp %s", sess.LaunchID, testLaunchID)
	}

	expected, err := smartonfhirtest.AuthURL(testFHIRAuthURL, testFHIRClientID, testFHIRRedirectURL, testLaunchID, sess.ID, fhirURL, testScopes)
	if err != nil {
		t.Fatalf("smartonfhirtest.AuthURL() -> %v, nil expected", err)
	}
	if diff := smartonfhirtest.DiffAuthURLs(actualRedirectURL, expected); diff != "" {
		t.Errorf("actual authentication URL does not equia to expected, diff %s", diff)
	}
}

func TestHandleFHIRRedirect(t *testing.T) {
	fhirURL := setupBackends()
	ss := sessiontest.NewMemoryStore()
	s := defaultServer(fhirURL, ss)
	ts := httptest.NewServer(http.HandlerFunc(s.handleFHIRRedirect))
	defer ts.Close()

	sessionID := "session-id"
	req, err := http.NewRequest("Get", fmt.Sprintf("%s?state=%s&code=%s", ts.URL, sessionID, testAuthCode), nil)
	if err != nil {
		t.Fatalf("http.NewRequest() -> %v, nil expected", err)
	}

	existingSession := &session.Session{ID: sessionID, FHIRURL: fhirURL, LaunchID: testLaunchID}
	if err := s.sm.Save(existingSession); err != nil {
		t.Fatalf("s.sm.Save() -> %v, nil expected", err)
	}
	addSessionCookie(req, existingSession)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("http.DefaultClient.Do() -> %v, nil expected", err)
	}
	if res.StatusCode != http.StatusOK {
		t.Errorf("got response with status %d, 200 expected", res.StatusCode)
	}

	sess, err := s.sm.Retrieve(req)
	if err != nil {
		t.Fatalf("s.sm.Retrieve() -> %v, nil expected", err)
	}
	expectedFHIRToken := &oauth2.Token{AccessToken: testFHIRAccessToken, RefreshToken: testFHIRRefreshToken, TokenType: "Bearer"}
	var tokenRaw map[string]string
	if err := json.Unmarshal([]byte(testFHIRTokenJSON), &tokenRaw); err != nil {
		t.Fatalf("json.Marshal() -> %v, nil expected", err)
	}
	expectedFHIRToken = expectedFHIRToken.WithExtra(tokenRaw)
	if diff := cmp.Diff(sess.FHIRToken, expectedFHIRToken, cmp.AllowUnexported(oauth2.Token{})); diff != "" {
		t.Errorf("fhir token in session does not equal to expected, diff %s", diff)
	}
}

func TestHandleFHIRRedirectError(t *testing.T) {
	fhirURL := setupBackends()
	tests := []struct {
		name, queryParameters string
		ss                    session.Store
		existingSession       *session.Session
		expectedHTTPStatus    int
		expectedMessage       string
	}{
		{
			name:               "missing session",
			queryParameters:    fmt.Sprintf("code=%s", testAuthCode),
			ss:                 sessiontest.NewMemoryStore(),
			expectedHTTPStatus: http.StatusUnauthorized,
			expectedMessage:    cookieSessionNotFoundErr,
		},
		{
			name:               "missing fhirURL in session",
			existingSession:    &session.Session{ID: "123", LaunchID: testLaunchID},
			ss:                 sessiontest.NewMemoryStore(),
			queryParameters:    fmt.Sprintf("state=123&code=%s", testAuthCode),
			expectedHTTPStatus: http.StatusUnauthorized,
			expectedMessage:    sessionMissingFHIRURLErr,
		},
		{
			name:               "missing launchID in session",
			existingSession:    &session.Session{ID: "123", FHIRURL: fhirURL},
			ss:                 sessiontest.NewMemoryStore(),
			queryParameters:    fmt.Sprintf("state=123&code=%s", testAuthCode),
			expectedHTTPStatus: http.StatusUnauthorized,
			expectedMessage:    sessionMissingLaunchIDErr,
		},
		{
			name:               "missing code in request",
			existingSession:    &session.Session{ID: "123", LaunchID: testLaunchID, FHIRURL: fhirURL},
			ss:                 sessiontest.NewMemoryStore(),
			queryParameters:    "state=123",
			expectedHTTPStatus: http.StatusBadRequest,
			expectedMessage:    requestQueryMissingCodeErr,
		},
		{
			name:               "missing state in request",
			existingSession:    &session.Session{ID: "123", LaunchID: testLaunchID, FHIRURL: fhirURL},
			ss:                 sessiontest.NewMemoryStore(),
			queryParameters:    fmt.Sprintf("code=%s", testAuthCode),
			expectedHTTPStatus: http.StatusBadRequest,
			expectedMessage:    requestQueryInvalidStateErr,
		},
		{
			name:               "unexpected state in request",
			existingSession:    &session.Session{ID: "123", LaunchID: testLaunchID, FHIRURL: fhirURL},
			ss:                 sessiontest.NewMemoryStore(),
			queryParameters:    fmt.Sprintf("state=122&code=%s", testAuthCode),
			expectedHTTPStatus: http.StatusBadRequest,
			expectedMessage:    requestQueryInvalidStateErr,
		},
		{
			name:               "retrieve session error",
			existingSession:    &session.Session{ID: "123", LaunchID: testLaunchID, FHIRURL: fhirURL},
			ss:                 sessiontest.NewMemoryStore().WithNextRetrieveErr(fmt.Errorf("retrieve session err")),
			queryParameters:    fmt.Sprintf("state=123&code=%s", testAuthCode),
			expectedHTTPStatus: http.StatusInternalServerError,
			expectedMessage:    serverRetrieveSessionErr,
		},
		{
			name:               "save session error",
			existingSession:    &session.Session{ID: "123", LaunchID: testLaunchID, FHIRURL: fhirURL},
			ss:                 sessiontest.NewMemoryStore().WithNextStoreExistingErr(fmt.Errorf("save session err")),
			queryParameters:    fmt.Sprintf("state=123&code=%s", testAuthCode),
			expectedHTTPStatus: http.StatusInternalServerError,
			expectedMessage:    serverSaveSessionErr,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s := defaultServer(fhirURL, test.ss)
			ts := httptest.NewServer(http.HandlerFunc(s.handleFHIRRedirect))
			defer ts.Close()

			req, err := http.NewRequest("Get", fmt.Sprintf("%s?%s", ts.URL, test.queryParameters), nil)
			if err != nil {
				t.Fatalf("http.NewRequest() -> %v, nil expected", err)
			}
			if test.existingSession != nil {
				if err := s.sm.Save(test.existingSession); err != nil {
					t.Fatalf("s.sm.Save() -> %v, nil expected", err)
				}
				addSessionCookie(req, test.existingSession)
			}

			res, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatalf("http.Get() -> %v, nil expected", err)
			}

			if status := res.StatusCode; status != test.expectedHTTPStatus {
				t.Errorf("server.handleFHIRRedirect returned wrong status code: got %v want %v",
					status, test.expectedHTTPStatus)
			}
			b, err := ioutil.ReadAll(res.Body)
			if err != nil {
				t.Fatalf("ioutil.ReadAll() -> %v, nil expected", err)
			}
			if actual := string(b); !strings.Contains(actual, test.expectedMessage) {
				t.Errorf("response message [%s] does not contain expected error message [%s]", actual, test.expectedMessage)
			}
		})
	}
}

func sessionFromResp(t *testing.T, ss session.Store, res *http.Response) *session.Session {
	sessionID := strings.Split(res.Header.Get("Set-Cookie"), ";")[0]
	sessionID = strings.Split(sessionID, "=")[1]
	b, err := ss.Retrieve(sessionID)
	if err != nil {
		t.Fatalf("ss.Retrieve(%s) -> %v, nil expected", sessionID, err)
	}
	session, err := session.FromBytes(b)
	if err != nil {
		t.Fatalf("session.FromBytes() -> %v, nil expected", err)
	}
	return session
}

func addSessionCookie(r *http.Request, s *session.Session) {
	cookie := &http.Cookie{Name: session.SessionCookieName, Value: s.ID, Expires: s.ExpiresAt}
	r.AddCookie(cookie)
}
