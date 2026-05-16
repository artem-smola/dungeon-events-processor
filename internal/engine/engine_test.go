package engine

import (
	"testing"

	"github.com/artem-smola/dungeon-events-processor/internal/model"
	"github.com/artem-smola/dungeon-events-processor/internal/report"
)

func TestCloseTimeEndsChallenge(t *testing.T) {
	cfg := model.Config{Floors: 2, Monsters: 1, OpenAt: "10:00:00", Duration: 1}
	e, err := New(cfg)
	if err != nil {
		t.Fatalf("new engine: %v", err)
	}

	events := []model.Event{
		{TimeSec: 10 * 3600, TimeText: "10:00:00", PlayerID: 1, EventID: model.EventRegister},
		{TimeSec: 10*3600 + 1, TimeText: "10:00:01", PlayerID: 1, EventID: model.EventEnter},
	}

	result := e.Run(events)
	if len(result.Reports) != 1 {
		t.Fatalf("expected one report, got %d", len(result.Reports))
	}
	if result.Reports[0].TotalSeconds != 3599 {
		t.Fatalf("expected total 3599, got %d", result.Reports[0].TotalSeconds)
	}
	if result.Reports[0].State != model.StateFail {
		t.Fatalf("expected FAIL after close, got %s", result.Reports[0].State)
	}
}

func TestRegistrationAfterDisqualification(t *testing.T) {
	cfg := model.Config{Floors: 2, Monsters: 1, OpenAt: "10:00:00", Duration: 1}
	e, err := New(cfg)
	if err != nil {
		t.Fatalf("new engine: %v", err)
	}

	events := []model.Event{
		{TimeSec: 10 * 3600, TimeText: "10:00:00", PlayerID: 77, EventID: model.EventEnter},
		{TimeSec: 10*3600 + 1, TimeText: "10:00:01", PlayerID: 77, EventID: model.EventRegister},
	}

	lines := report.Format(e.Run(events))
	if len(lines) < 2 {
		t.Fatalf("unexpected output: %v", lines)
	}
	if lines[0] != "[10:00:00] Player [77] is disqualified" {
		t.Fatalf("unexpected first line: %s", lines[0])
	}
	for _, line := range lines {
		if line == "[10:00:01] Player [77] registered" {
			t.Fatalf("player must not be registered after disqualification")
		}
	}
}

func TestRunResetState(t *testing.T) {
	cfg := model.Config{Floors: 2, Monsters: 1, OpenAt: "10:00:00", Duration: 1}
	e, err := New(cfg)
	if err != nil {
		t.Fatalf("new engine: %v", err)
	}

	first := []model.Event{{TimeSec: 10 * 3600, TimeText: "10:00:00", PlayerID: 1, EventID: model.EventRegister}}
	second := []model.Event{{TimeSec: 10 * 3600, TimeText: "10:00:00", PlayerID: 2, EventID: model.EventRegister}}

	_ = e.Run(first)
	res := e.Run(second)
	if len(res.Reports) != 1 || res.Reports[0].PlayerID != 2 {
		t.Fatalf("run must reset state, got reports: %#v", res.Reports)
	}
}

func TestNewEndAfterMidnight(t *testing.T) {
	_, err := New(model.Config{Floors: 2, Monsters: 1, OpenAt: "23:00:00", Duration: 2})
	if err == nil {
		t.Fatalf("expected rejection because of end after midnight")
	}
}

func TestAutoCompleteFloorsWithZeroEnemies(t *testing.T) {
	cfg := model.Config{Floors: 3, Monsters: 0, OpenAt: "10:00:00", Duration: 1}
	e, err := New(cfg)
	if err != nil {
		t.Fatalf("new engine: %v", err)
	}

	events := []model.Event{
		{TimeSec: 10 * 3600, TimeText: "10:00:00", PlayerID: 9, EventID: model.EventRegister},
		{TimeSec: 10*3600 + 1, TimeText: "10:00:01", PlayerID: 9, EventID: model.EventEnter},
		{TimeSec: 10*3600 + 2, TimeText: "10:00:02", PlayerID: 9, EventID: model.EventNext},
		{TimeSec: 10*3600 + 3, TimeText: "10:00:03", PlayerID: 9, EventID: model.EventNext},
		{TimeSec: 10*3600 + 4, TimeText: "10:00:04", PlayerID: 9, EventID: model.EventBossIn},
		{TimeSec: 10*3600 + 5, TimeText: "10:00:05", PlayerID: 9, EventID: model.EventKillBoss},
		{TimeSec: 10*3600 + 6, TimeText: "10:00:06", PlayerID: 9, EventID: model.EventLeave},
	}

	res := e.Run(events)
	if len(res.Reports) != 1 {
		t.Fatalf("expected one report, got %d", len(res.Reports))
	}
	r := res.Reports[0]
	if r.State != model.StateSuccess {
		t.Fatalf("expected SUCCESS, got %s", r.State)
	}
	if r.AvgFloorSeconds != 0 {
		t.Fatalf("expected zero avg floor time, got %d", r.AvgFloorSeconds)
	}
}
