package session

import (
	"errors"
	"fmt"
	"net/http"
	"time"
)

const sessionCookieName = "session"

// ErrNotFound is the error returned when something is not found.
var ErrNotFound = errors.New("not found")

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
	sessionID       func() string
	sessionDuration time.Duration
}

// NewManager creates a new Manager using the given Store.
func NewManager(ss Store, sessionID func() string, sessionDuration time.Duration) *Manager {
	return &Manager{store: ss, sessionID: sessionID, sessionDuration: sessionDuration}
}

// New creates a new session and set cookie containning the encoded session id in both HTTP
// request and response.
func (m *Manager) New(w http.ResponseWriter, r *http.Request) (*Session, error) {
	s, err := m.create()
	if err != nil {
		return nil, err
	}
	cookie := &http.Cookie{Name: sessionCookieName, Value: s.ID, Expires: s.ExpiresAt}
	http.SetCookie(w, cookie)
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

// Save saves the Session by overriding the existing one.
func (m *Manager) Save(session *Session) error {
	b, err := session.Bytes()
	if err != nil {
		return err
	}
	return m.store.Store(session.ID, b)
}

// create creates a new session with the given expiration time.
func (m *Manager) create() (*Session, error) {
	id := m.sessionID()
	sess := &Session{ID: id, ExpiresAt: time.Now().Add(m.sessionDuration)}
	if err := m.Save(sess); err != nil {
		return nil, err
	}
	return sess, nil
}

// find finds and returns the Session whose id mathces the given one.
// Returns error if no matching Sessions are found.
func (m *Manager) find(id string) (*Session, error) {
	v, err := m.store.Retrieve(id)
	if err != nil {
		return nil, err
	}
	return FromBytes(v)
}
