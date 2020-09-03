package session

import (
	"reflect"
	"testing"
	"time"
)

func TestStoreManager(t *testing.T) {
	sessionID := "test-id"
	sm := NewStoreManager(NewMemoryStore(), func() string { return sessionID })
	expiresAt := time.Now().Add(30 * time.Minute)

	// test Create
	sess, err := sm.Create(expiresAt)
	if err != nil {
		t.Fatalf("sm.Create() -> %v, expect nil", err)
		return
	}
	expectedSess := &Session{ID: sessionID, ExpiresAt: expiresAt}
	if !reflect.DeepEqual(sess, expectedSess) {
		t.Errorf("created session %v does not equal to expected %v", sess, expectedSess)
	}

	// test Find
	sess, err = sm.Find(sessionID)
	if err != nil {
		t.Fatalf("sm.Find() -> %v, expect nil", err)
		return
	}
	if !reflect.DeepEqual(sess, expectedSess) {
		t.Errorf("found session %v does not equal to expected %v", sess, expectedSess)
	}

	// test save - override the existing one
	sess.Put("key", "val")
	if err = sm.Save(sess); err != nil {
		t.Fatalf("sm.Save() -> %v, expect nil", err)
		return
	}
	expectedSess.Put("key", "val")
	sess, err = sm.Find(sessionID)
	if err != nil {
		t.Fatalf("sm.Find() -> %v, expect nil", err)
		return
	}
	if !reflect.DeepEqual(sess, expectedSess) {
		t.Errorf("found session %v does not equal to expected %v", sess, expectedSess)
	}

	// test save - create new one
	expectedSess = &Session{ID: "test-id-2", ExpiresAt: expiresAt}
	if err = sm.Save(expectedSess); err != nil {
		t.Fatalf("sm.Save() -> %v, expect nil", err)
		return
	}
	sess, err = sm.Find("test-id-2")
	if err != nil {
		t.Fatalf("sm.Find() -> %v, expect nil", err)
		return
	}
	if !reflect.DeepEqual(sess, expectedSess) {
		t.Errorf("found session %v does not equal to expected %v", sess, expectedSess)
	}
}
