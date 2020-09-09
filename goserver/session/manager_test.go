package session

import (
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestManager(t *testing.T) {
	sessionID := "test-id"
	m := NewManager(NewMemoryStore(), "session-secret", func() string { return sessionID }, 30*time.Minute)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest("Get", "https://test.com", nil)
	// test New
	sess, err := m.New(rr, req)
	if err != nil {
		t.Fatalf("sm.New() -> %v, expect nil", err)
		return
	}
	expectedSess := &Session{ID: sessionID}
	if !reflect.DeepEqual(sess, expectedSess) {
		t.Errorf("created session %v does not equal to expected %v", sess, expectedSess)
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

	// Test Retrieve
	sess, err = m.Retrieve(req)
	if err != nil {
		t.Fatalf("m.Retrieve() -> %v, expect nil", err)
		return
	}
	if !reflect.DeepEqual(sess, expectedSess) {
		t.Errorf("found session %v does not equal to expected %v", sess, expectedSess)
	}

	// test save - override the existing one
	sess.Set("key", "val")
	if err = m.Save(sess); err != nil {
		t.Fatalf("m.Save() -> %v, expect nil", err)
		return
	}
	expectedSess.Set("key", "val")
	sess, err = m.Retrieve(req)
	if err != nil {
		t.Fatalf("m.Find() -> %v, expect nil", err)
		return
	}
	if !reflect.DeepEqual(sess, expectedSess) {
		t.Errorf("found session %v does not equal to expected %v", sess, expectedSess)
	}
}

func TestManagerSaveNonexistentSessionError(t *testing.T) {
	m := NewManager(NewMemoryStore(), "session-secret", func() string { return "test-id" }, 30*time.Minute)
	if err := m.Save(&Session{ID: "test-id"}); err == nil {
		t.Fatal("m.Save() -> nil, expect error")
		return
	}
}
