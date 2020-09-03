package session

import (
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

func TestNew(t *testing.T) {
	sm := NewStoreManager(NewMemoryStore(), func() string { return "test-id" })
	rr := httptest.NewRecorder()
	r := httptest.NewRequest("method", "https://test.com", nil)
	sess, err := New(sm, rr, r)
	if err != nil {
		t.Errorf("cannot create new session, got err: %v", err)
	}

	// Make sure session is created in manager.
	sess, err = sm.Find(sess.ID)
	if err != nil {
		t.Fatalf("cannot find session in session manager, got err: %v", err)
	}
	if sess == nil {
		t.Fatalf("cannot find session in session manager, got nil")
	}

	// Make sure cookie containing encoded session id is set in response.
	cookies := rr.HeaderMap["Set-Cookie"]
	if len(cookies) < 1 {
		t.Fatal("cannot find cookie response")
	}
	esid := encodeSessionID(sess.ID)
	if !strings.Contains(cookies[0], cookieName) {
		t.Fatalf("cookie %s set in response", cookieName)
	}
	if !strings.Contains(cookies[0], esid) {
		t.Fatal("no session id set in response cookie")
	}

	// Make sure cookie containing encoded session id is set in request.
	cookie, err := r.Cookie(cookieName)
	if err != nil {
		t.Fatalf("cookie %v not set request, got err %v", cookieName, err)
	}
	if cookie == nil {
		t.Fatalf("cookie %s not set request, got nil", cookieName)
	}
	if cookie.Value != esid {
		t.Fatalf("value of cookie %s (%s) does not match expected %s", cookieName, cookie.Value, esid)
	}
}

func TestFind(t *testing.T) {
	sm := NewStoreManager(NewMemoryStore(), func() string { return "test-id" })
	r := httptest.NewRequest("method", "https://test.com", nil)
	sess, err := New(sm, httptest.NewRecorder(), r)
	if err != nil {
		t.Errorf("cannot create new session, got err: %v", err)
	}

	foundSess, err := Find(sm, r)
	if err != nil {
		t.Fatalf("cannot find session, got err: %v", err)
	}

	if !reflect.DeepEqual(sess, foundSess) {
		t.Fatal("The found session does not equal to the created one")
	}
}
