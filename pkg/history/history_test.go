package history

import (
	"testing"
	"time"
)

func TestInMemoryStoreAppendListClear(t *testing.T) {
	s := NewInMemoryStore()
	if entries, _ := s.List(); len(entries) != 0 {
		t.Fatalf("expected empty store initially")
	}

	now := time.Now().Unix()
	entry := Entry{Query: "SELECT 1", Timestamp: now, Source: "test"}
	if err := s.Append(entry); err != nil {
		t.Fatalf("unexpected append error: %v", err)
	}

	entries, _ := s.List()
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Query != "SELECT 1" {
		t.Fatalf("unexpected query in entry: %s", entries[0].Query)
	}

	if err := s.Clear(); err != nil {
		t.Fatalf("unexpected clear error: %v", err)
	}
	if entries, _ := s.List(); len(entries) != 0 {
		t.Fatalf("expected empty store after clear")
	}
}
