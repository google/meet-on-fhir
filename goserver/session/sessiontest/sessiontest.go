// Package sessiontest provides testing utilities for session package.
package sessiontest

import (
	"errors"
)

// ErrNotFound is the error returned when something is not found.
var ErrNotFound = errors.New("not found")

// MemoryStore is an in-memory implementaion of Store.
// It is not thread-safe and should be used by tests only.
type MemoryStore struct {
	items            map[string]interface{}
	storeErr         error
	storeExistingErr error
}

// NewMemoryStore creates a new MemoryStore.
func NewMemoryStore() *MemoryStore {
	items := make(map[string]interface{})
	return &MemoryStore{items: items}
}

// NewMemoryStoreWithError creates a new MemoryStore that returns the provided error
// under centain cases.
func NewMemoryStoreWithError(storeErr, storeExistingErr error) *MemoryStore {
	items := make(map[string]interface{})
	return &MemoryStore{items: items, storeErr: storeErr, storeExistingErr: storeExistingErr}
}

// Store implements Store of Store interface.
func (s *MemoryStore) Store(key string, val []byte) error {
	if s.storeErr != nil {
		return s.storeErr
	}

	if _, ok := s.items[key]; ok && s.storeExistingErr != nil {
		return s.storeExistingErr
	}
	s.items[key] = val
	return nil
}

// Retrieve implements Retrieve of Store interface.
func (s *MemoryStore) Retrieve(key string) ([]byte, error) {
	if s.storeErr != nil {
		return nil, s.storeErr
	}
	i, ok := s.items[key]
	if !ok {
		return nil, ErrNotFound
	}
	return i.([]byte), nil
}
