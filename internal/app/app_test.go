package app

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunSampleCase(t *testing.T) {
	base := filepath.Join("..", "..", "testdata", "sample")
	got, err := Run(filepath.Join(base, "config.json"), filepath.Join(base, "events.txt"))
	if err != nil {
		t.Fatalf("run: %v", err)
	}

	expectedRaw, err := os.ReadFile(filepath.Join(base, "expected.txt"))
	if err != nil {
		t.Fatalf("read expected: %v", err)
	}

	expected := strings.Split(strings.TrimSpace(string(expectedRaw)), "\n")
	if len(got) != len(expected) {
		t.Fatalf("line count mismatch: got=%d want=%d", len(got), len(expected))
	}

	for i := range expected {
		if got[i] != expected[i] {
			t.Fatalf("line %d mismatch\ngot:  %s\nwant: %s", i+1, got[i], expected[i])
		}
	}
}
