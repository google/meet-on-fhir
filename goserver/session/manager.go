package session

import (
	"fmt"
	"time"
)

// Manager manages sessions for the server.
type Manager interface {
	// Create creates a new session with the given expire time.
	Create(expireAt time.Time) (*Session, error)
	// Find finds the session for the given id.
	Find(sid string) (*Session, error)
	// Save saves the given session by override the existing one.
	// It will return an error if no existing one is found.
	Save(session *Session) error
}

// InMemorySessionManager is an in-memory implementation of Manager.
// It's not thread-safe and should be used for testing only.
type InMemorySessionManager struct {
	nextID   int
	sessions map[string]*Session
}

// NewInMemorySessionManager creates a new InMemorySessionManager.
func NewInMemorySessionManager() *InMemorySessionManager {
	return &InMemorySessionManager{sessions: make(map[string]*Session)}
}

// Create implements Create of Manager.
func (m *InMemorySessionManager) Create(expireAt time.Time) (*Session, error) {
	s := &Session{
		expireAt: expireAt,
		sid:      string(m.nextID),
	}
	m.nextID++
	m.sessions[s.SessionID()] = s
	return s, nil
}

// Find implements Find of Manager.
func (m *InMemorySessionManager) Find(sid string) (*Session, error) {
	sess := m.sessions[sid]
	if sess == nil {
		return nil, fmt.Errorf("no session found for sid %s", sid)
	}
	return sess, nil
}

// Save implements Save of Manager.
func (m *InMemorySessionManager) Save(s *Session) error {
	oldSession, err := m.Find(s.SessionID())
	if err != nil || oldSession == nil {
		return fmt.Errorf("session %s does not exist", s.SessionID())
	}
	m.sessions[s.SessionID()] = s
	return nil
}
