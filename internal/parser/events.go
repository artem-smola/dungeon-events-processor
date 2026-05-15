package parser

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/artem-smola/dungeon-events-processor/internal/model"
)

var eventRe = regexp.MustCompile(`^\[(\d{2}:\d{2}:\d{2})\]\s+(\d+)\s+(\d+)(?:\s+(.*))?$`)

func ReadEvents(path string) ([]model.Event, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open events: %w", err)
	}
	defer f.Close()

	s := bufio.NewScanner(f)
	var events []model.Event

	lineCount := 0
	prevTime := -1
	for s.Scan() {
		lineCount++
		line := strings.TrimSpace(s.Text())
		if line == "" {
			continue
		}

		match := eventRe.FindStringSubmatch(line)
		if len(match) == 0 {
			return nil, fmt.Errorf("invalid event format on line %d", lineCount)
		}

		timeText := match[1]
		timeSec, err := ParseTime(timeText)
		if err != nil {
			return nil, fmt.Errorf("invalid time on line %d: %w", lineCount, err)
		}
		if prevTime > timeSec {
			return nil, fmt.Errorf("non-monotonic event time on line %d", lineCount)
		}
		prevTime = timeSec

		playerID, _ := strconv.Atoi(match[2])
		eventID, _ := strconv.Atoi(match[3])

		extra := ""
		if len(match) >= 5 {
			extra = strings.TrimSpace(match[4])
		}

		events = append(events, model.Event{
			TimeSec:  timeSec,
			TimeText: timeText,
			PlayerID: playerID,
			EventID:  model.EventID(eventID),
			Extra:    extra,
		})
	}

	if err := s.Err(); err != nil {
		return nil, fmt.Errorf("scan events: %w", err)
	}

	return events, nil
}
