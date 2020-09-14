package session

import (
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/google/meet-on-fhir/session/sessiontest"
)

func TestManager(t *testing.T) {
	m := NewManager(sessiontest.NewMemoryStore(), 30*time.Minute)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest("Get", "https://test.com", nil)
	// Test New
	sess, err := m.New(rr, req)
	if err != nil {
		t.Fatalf("sm.New() -> %v, expect nil", err)
	}
	// Check session cookie set in response
	cookies := rr.HeaderMap["Set-Cookie"]
	if len(cookies) < 1 {
		t.Fatal("\"Set-Cookie\" header missing in response")
	}
	if !strings.Contains(cookies[0], sessionCookieName) {
		t.Fatalf("cookie %s not set in response", sessionCookieName)
	}
	if !strings.Contains(cookies[0], sess.ID) {
		t.Fatal("wrong session id set in response cookie")
	}

	expected := &Session{ID: sess.ID, ExpiresAt: sess.ExpiresAt.Truncate(0)}
	// Test Retrieve
	found, err := m.Retrieve(req)
	if err != nil {
		t.Fatalf("m.Retrieve() -> %v, expect nil", err)
	}
	if !reflect.DeepEqual(found, expected) {
		t.Errorf("found session %v does not equal to expected %v", found, expected)
	}

	// Test save - override the existing one
	sess.FHIRURL = "url"
	if err = m.Save(sess); err != nil {
		t.Fatalf("m.Save() -> %v, expect nil", err)
	}
	expected.FHIRURL = "url"
	found, err = m.Retrieve(req)
	if err != nil {
		t.Fatalf("m.Find() -> %v, expect nil", err)
	}
	if !reflect.DeepEqual(found, expected) {
		t.Errorf("found session %v does not equal to expected %v", found, expected)
	}
}
