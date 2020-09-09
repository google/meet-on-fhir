package session

import (
	"encoding/json"
	"time"
)

// Session stores necessary information for a telehealth session.
type Session struct {
	ID        string    `json:"id"`
	FHIRURL   string    `json:"fhir_url"`
	ExpiresAt time.Time `json:"expires_at"`
}

// Bytes converts the session to JSON bytes.
func (s *Session) Bytes() ([]byte, error) {
	return json.Marshal(s)
}

// FromBytes constructs a Session with the given JSON bytes.
func FromBytes(data []byte) (*Session, error) {
	s := &Session{}
	if err := json.Unmarshal(data, s); err != nil {
		return nil, err
	}
	return s, nil
}
