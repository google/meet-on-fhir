package session

import (
	"time"
)

// MemoryStore is an in-memory implementaion of Store.
// It is not thread-safe and should be used by tests only.
type MemoryStore struct {
	items map[string]interface{}
}

type item struct {
	expiresAt time.Time
	val       []byte
}

// NewMemoryStore creates a new MemoryStore.
func NewMemoryStore() *MemoryStore {
	items := make(map[string]interface{})
	return &MemoryStore{items: items}
}

// Store implements Store of Store interface.
func (s *MemoryStore) Store(key string, val []byte, expiresAt time.Time) error {
	s.items[key] = &item{expiresAt: expiresAt, val: val}
	return nil
}

// Retrieve implements Retrieve of Store interface.
func (s *MemoryStore) Retrieve(key string) ([]byte, time.Time, error) {
	i, ok := s.items[key]
	if !ok {
		return nil, time.Time{}, ErrNotFound
	}
	return i.(*item).val, i.(*item).expiresAt, nil
}
