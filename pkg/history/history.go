package history

// Entry represents a single SQL history entry.
type Entry struct {
	Query     string `json:"query"`
	Timestamp int64  `json:"timestamp"`
	Source    string `json:"source"`
}

// Store is a simple in-memory store implementation used for tests and local usage.
type Store interface {
	Append(entry Entry) error
	List() ([]Entry, error)
	Clear() error
}

// inMemoryStore is a minimal thread-safe in-memory implementation.
type inMemoryStore struct {
	entries []Entry
}

// NewInMemoryStore returns a fresh in-memory store.
func NewInMemoryStore() Store {
	return &inMemoryStore{entries: make([]Entry, 0)}
}

func (s *inMemoryStore) Append(entry Entry) error {
	s.entries = append(s.entries, entry)
	return nil
}

func (s *inMemoryStore) List() ([]Entry, error) {
	// return a copy to avoid external mutation
	out := make([]Entry, len(s.entries))
	copy(out, s.entries)
	return out, nil
}

func (s *inMemoryStore) Clear() error {
	s.entries = s.entries[:0]
	return nil
}
