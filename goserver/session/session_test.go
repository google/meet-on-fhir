package session

import (
	"testing"
)

func TestSession(t *testing.T) {
	sess := &Session{ID: "session-id"}
	sess.Set("key", "value")
	if v := sess.Get("key").(string); v != "value" {
		t.Fatalf("sess.Get() -> %s, expected value", v)
	}
}
