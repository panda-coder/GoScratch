package console

import (
	"testing"
	"time"
)

func TestPublishAndSnapshot(t *testing.T) {
	c := New()
	e := Event{Type: "stdout", Payload: "hello", Timestamp: time.Now()}
	c.Publish(e)
	snap := c.Snapshot()
	if len(snap) != 1 {
		t.Fatalf("expected 1 event, got %d", len(snap))
	}
	if snap[0].Payload != "hello" {
		t.Fatalf("unexpected payload: %s", snap[0].Payload)
	}
}

func TestFromStdout(t *testing.T) {
	c := New()
	go FromStdout(c, "line1\n\nline2\n")
	time.Sleep(10 * time.Millisecond)
	snap := c.Snapshot()
	if len(snap) != 2 {
		t.Fatalf("expected 2 events, got %d", len(snap))
	}
	if snap[0].Payload != "line1" || snap[1].Payload != "line2" {
		t.Fatalf("unexpected payloads: %v %v", snap[0].Payload, snap[1].Payload)
	}
}
