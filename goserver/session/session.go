package session

// Session stores necessary information for a telehealth session.
type Session struct {
	ID    string
	Value map[string]interface{}
}

// Set sets the value for a key.
func (s *Session) Set(key string, val interface{}) {
	if s.Value == nil {
		s.Value = make(map[string]interface{})
	}
	s.Value[key] = val
}

// Get returns the value for the given key.
func (s *Session) Get(key string) interface{} {
	if s.Value == nil {
		return nil
	}
	v, ok := s.Value[key]
	if !ok {
		return nil
	}
	return v
}
