package report

import (
	"reflect"
	"testing"

	"github.com/artem-smola/dungeon-events-processor/internal/model"
)

func TestFormat(t *testing.T) {
	in := model.Result{
		Logs: []model.LogEntry{
			{TimeText: "12:00:00", PlayerID: 10, Message: "registered"},
			{TimeText: "12:01:00", PlayerID: 10, Message: "entered the dungeon"},
		},
		Reports: []model.PlayerReport{
			{State: "SUCCESS", PlayerID: 10, TotalSeconds: 3661, AvgFloorSeconds: 120, BossSeconds: 41, Health: 77},
		},
	}

	got := Format(in)
	want := []string{
		"[12:00:00] Player [10] registered",
		"[12:01:00] Player [10] entered the dungeon",
		"Final report:",
		"[SUCCESS] 10 [01:01:01, 00:02:00, 00:00:41] HP:77",
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected output\nwant: %v\ngot: %v", want, got)
	}
}
