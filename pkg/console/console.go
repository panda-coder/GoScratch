package console

import (
	"strings"
	"sync"
	"time"
)

type Event struct {
	Type      string    `json:"type"`
	Payload   string    `json:"payload"`
	Timestamp time.Time `json:"timestamp"`
}

type Console struct {
	mu     sync.Mutex
	events []Event
	ch     chan Event
}

func New() *Console {
	return &Console{ch: make(chan Event, 256)}
}

func (c *Console) Publish(evt Event) {
	c.mu.Lock()
	c.events = append(c.events, evt)
	c.mu.Unlock()
	select {
	case c.ch <- evt:
	default:
	}
}

func (c *Console) Events() <-chan Event { return c.ch }

func (c *Console) Snapshot() []Event {
	c.mu.Lock()
	defer c.mu.Unlock()
	out := make([]Event, len(c.events))
	copy(out, c.events)
	return out
}

// FromStdout publishes each non-empty line of stdout as a stdout event.
func FromStdout(c *Console, stdout string) {
	if c == nil {
		return
	}
	lines := strings.Split(stdout, "\n")
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		evt := Event{Type: "stdout", Payload: line, Timestamp: time.Now()}
		c.Publish(evt)
	}
}
