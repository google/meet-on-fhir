package sessiontest

import (
	"errors"
)

// ErrNotFound is the error returned when something is not found.
var ErrNotFound = errors.New("not found")

// MemoryStore is an in-memory implementaion of Store.
// It is not thread-safe and should be used by tests only.
type MemoryStore struct {
	items map[string]interface{}
}

// NewMemoryStore creates a new MemoryStore.
func NewMemoryStore() *MemoryStore {
	items := make(map[string]interface{})
	return &MemoryStore{items: items}
}

// Store implements Store of Store interface.
func (s *MemoryStore) Store(key string, val []byte) error {
	s.items[key] = val
	return nil
}

// Retrieve implements Retrieve of Store interface.
func (s *MemoryStore) Retrieve(key string) ([]byte, error) {
	i, ok := s.items[key]
	if !ok {
		return nil, ErrNotFound
	}
	return i.([]byte), nil
}
