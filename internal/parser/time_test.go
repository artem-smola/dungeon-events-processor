package parser

import "testing"

func TestParseClock(t *testing.T) {
	sec, err := ParseTime("23:59:58")
	if err != nil {
		t.Fatalf("parse clock: %v", err)
	}
	if sec != 23*3600+59*60+58 {
		t.Fatalf("unexpected seconds: %d", sec)
	}

	if _, err := ParseTime("24:00:00"); err == nil {
		t.Fatalf("expected invalid time error")
	}
}
