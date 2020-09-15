// Package session implements server session.
//
// A Manager should be created for creating/updaing sessions. Manager will store/retrieve session
// data using the provided Store.
package session

import (
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
)

const sessionCookieName = "session"

// Store provides functions to store/retrieve keyed binary data.
type Store interface {
	// Store stores a key-value pair.
	Store(key string, val []byte) error
	// Retrieve retrieves the value for the key.
	Retrieve(key string) ([]byte, error)
}

// Manager manages sessions.
type Manager struct {
	store           Store
	sessionDuration time.Duration
}

// NewManager creates a new Manager using the given Store.
func NewManager(ss Store, sessionDuration time.Duration) *Manager {
	return &Manager{store: ss, sessionDuration: sessionDuration}
}

// New creates a new session and set session cookie in both HTTP request and response.
func (m *Manager) New(w http.ResponseWriter, r *http.Request) (*Session, error) {
	s, err := m.create()
	if err != nil {
		return nil, err
	}
	// TODO(Issue #21): Figure out whehter signing the session ID in cookie is needed.
	cookie := &http.Cookie{Name: sessionCookieName, Value: s.ID, Expires: s.ExpiresAt}
	http.SetCookie(w, cookie)
	// Add cookie on the incoming request as well to perform like a middleware which
	// makes it easier to add more request processors relying on sessions.
	r.AddCookie(cookie)
	return s, nil
}

// Retrieve returns the session whose id matches the session id in HTTP request cookie.
func (m *Manager) Retrieve(r *http.Request) (*Session, error) {
	cookie, err := r.Cookie(sessionCookieName)
	if err != nil {
		return nil, err
	}
	sid := cookie.Value
	if sid == "" {
		return nil, fmt.Errorf("session cookie value is empty")
	}
	return m.find(sid)
}

// Save saves the Session by either creating a new one or overriding the existing one.
func (m *Manager) Save(session *Session) error {
	b, err := session.Bytes()
	if err != nil {
		return err
	}
	return m.store.Store(session.ID, b)
}

// create creates a new session and stores in Store.
func (m *Manager) create() (*Session, error) {
	id := uuid.New().String()
	sess := &Session{ID: id, ExpiresAt: time.Now().Add(m.sessionDuration)}
	if err := m.Save(sess); err != nil {
		return nil, err
	}
	return sess, nil
}

// find returns a matching Session in Store based on the id.
func (m *Manager) find(id string) (*Session, error) {
	v, err := m.store.Retrieve(id)
	if err != nil {
		return nil, err
	}
	return FromBytes(v)
}
