package parser

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReadEventsWithExtraWords(t *testing.T) {
	content := "[10:00:00] 1 1\n[10:01:02] 1 9 out of mana\n"
	path := filepath.Join(t.TempDir(), "events.txt")
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write temp file: %v", err)
	}

	events, err := ReadEvents(path)
	if err != nil {
		t.Fatalf("read events: %v", err)
	}
	if len(events) != 2 {
		t.Fatalf("expected 2 events, got %d", len(events))
	}
	if events[1].Extra != "out of mana" {
		t.Fatalf("incorrect extra words %q", events[1].Extra)
	}
}

func TestReadEventsNonMonotonicTime(t *testing.T) {
	content := "[10:00:01] 1 1\n[10:00:00] 1 2\n"
	path := filepath.Join(t.TempDir(), "events.txt")
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write temp file: %v", err)
	}

	_, err := ReadEvents(path)
	if err == nil {
		t.Fatalf("expected monotonicity error")
	}
}
