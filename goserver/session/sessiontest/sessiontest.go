// Package sessiontest provides testing utilities for session package.
package sessiontest

// MemoryStore is an in-memory implementaion of Store.
// It is not thread-safe and should be used by tests only.
type MemoryStore struct {
	items                map[string]interface{}
	nextStoreErr         error
	nextStoreExistingErr error
	nextRetrieveErr      error
}

// NewMemoryStore creates a new MemoryStore.
func NewMemoryStore() *MemoryStore {
	items := make(map[string]interface{})
	return &MemoryStore{items: items}
}

// WithNextStoreErr returns the same Store with a nextStoreErr.
func (s *MemoryStore) WithNextStoreErr(err error) *MemoryStore {
	s.nextStoreErr = err
	return s
}

// WithNextStoreExistingErr returns the same Store with a nextStoreExistingErr.
func (s *MemoryStore) WithNextStoreExistingErr(err error) *MemoryStore {
	s.nextStoreExistingErr = err
	return s
}

// WithNextRetrieveErr returns the same Store with a nextRetrieveErr.
func (s *MemoryStore) WithNextRetrieveErr(err error) *MemoryStore {
	s.nextRetrieveErr = err
	return s
}

// Store implements Store of Store interface.
func (s *MemoryStore) Store(key string, val []byte) error {
	if s.nextStoreErr != nil {
		return s.nextStoreErr
	}

	if _, ok := s.items[key]; ok && s.nextStoreExistingErr != nil {
		return s.nextStoreExistingErr
	}

	s.items[key] = val
	return nil
}

// Retrieve implements Retrieve of Store interface.
func (s *MemoryStore) Retrieve(key string) ([]byte, error) {
	if s.nextRetrieveErr != nil {
		return nil, s.nextRetrieveErr
	}
	i := s.items[key]
	return i.([]byte), nil
}
