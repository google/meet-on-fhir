package session

import (
	"testing"
	"time"

	"golang.org/x/oauth2"

	"github.com/google/go-cmp/cmp"
	"github.com/google/meet-on-fhir/smartonfhir"
)

func TestSession(t *testing.T) {
	tests := []struct {
		name    string
		session *Session
	}{
		{
			name:    "empty session",
			session: new(Session),
		},
		{
			name:    "session with all fields",
			session: &Session{ID: "session-id", FHIRURL: "fhir-url", ExpiresAt: time.Now(), FHIRContext: &smartonfhir.FHIRContext{Token: &oauth2.Token{AccessToken: "access-token", RefreshToken: "refresh-token", TokenType: "Bearer", Expiry: time.Now()}, Scope: "scope", EncounterID: "e123", PatientID: "p123"}},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			b, err := test.session.Bytes()
			if err != nil {
				t.Fatalf("test.session.Bytes() -> %v, nil expected", err)
			}
			decoded, err := FromBytes(b)
			if err != nil {
				t.Fatalf("FromBytes(b)-> %v, nil expected", err)
			}
			if diff := cmp.Diff(test.session, decoded, cmp.AllowUnexported(oauth2.Token{})); diff != "" {
				t.Errorf("decoded session does not equal to the original one, diff %s", diff)
			}
		})
	}
}
