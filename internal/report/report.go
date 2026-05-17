package report

import (
	"fmt"

	"github.com/artem-smola/dungeon-events-processor/internal/model"
)

func Format(result model.Result) []string {
	out := make([]string, 0, len(result.Logs)+len(result.Reports)+1)

	for _, entry := range result.Logs {
		out = append(out, fmt.Sprintf("[%s] Player [%d] %s", entry.TimeText, entry.PlayerID, entry.Message))
	}

	out = append(out, "Final report:")
	for _, r := range result.Reports {
		out = append(out, fmt.Sprintf("[%s] %d [%s, %s, %s] HP:%d",
			r.State,
			r.PlayerID,
			formatDuration(r.TotalSeconds),
			formatDuration(r.AvgFloorSeconds),
			formatDuration(r.BossSeconds),
			r.Health,
		))
	}

	return out
}

func formatDuration(sec int) string {
	if sec < 0 {
		sec = 0
	}
	h := sec / 3600
	m := (sec % 3600) / 60
	s := sec % 60
	return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
}
