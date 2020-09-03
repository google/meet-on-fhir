package session

import (
	"encoding/json"
	"errors"
	"time"
)

// ErrNotFound is the error returned when something is not found.
var ErrNotFound = errors.New("not found")

// Store provides functions to store/retrieve keyed binary data.
type Store interface {
	// Store stores a key-value pair with an expiration time.
	Store(key string, val []byte, expiresAt time.Time) error
	// Retrieve retrieves the value and the expiration time for the key.
	Retrieve(key string) ([]byte, time.Time, error)
}

// StoreManager manages sessions.
type StoreManager struct {
	store     Store
	sessionID func() string
}

// NewStoreManager creates a new StoreManager using the given Store.
func NewStoreManager(ss Store, sessionID func() string) *StoreManager {
	return &StoreManager{store: ss, sessionID: sessionID}
}

// Create creates a new session with the given expiration time.
func (m *StoreManager) Create(expiresAt time.Time) (*Session, error) {
	id := m.sessionID()
	sess := &Session{ExpiresAt: expiresAt, ID: id}
	if err := m.store.Store(id, nil, expiresAt); err != nil {
		return nil, err
	}
	return sess, nil
}

// Find finds and returns the Session whose id mathces the given one.
// Returns error if no matching Sessions are found.
func (m *StoreManager) Find(id string) (*Session, error) {
	v, e, err := m.store.Retrieve(id)
	if err != nil {
		return nil, err
	}
	var val map[string]interface{}
	if v != nil {
		if err := json.Unmarshal(v, &val); err != nil {
			return nil, err
		}
	}
	return &Session{ID: id, ExpiresAt: e, Value: val}, nil
}

// Save saves the Session. An existing session with the same id will be overriden.
func (m *StoreManager) Save(session *Session) error {
	js, err := json.Marshal(session.Value)
	if err != nil {
		return err
	}
	return m.store.Store(session.ID, js, session.ExpiresAt)
}
