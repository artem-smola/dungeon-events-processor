package app

import (
	"github.com/artem-smola/dungeon-events-processor/internal/engine"
	"github.com/artem-smola/dungeon-events-processor/internal/parser"
	"github.com/artem-smola/dungeon-events-processor/internal/report"
)

func Run(configPath, eventsPath string) ([]string, error) {
	cfg, err := parser.ReadConfig(configPath)
	if err != nil {
		return nil, err
	}

	events, err := parser.ReadEvents(eventsPath)
	if err != nil {
		return nil, err
	}

	e, err := engine.New(cfg)
	if err != nil {
		return nil, err
	}

	return report.Format(e.Run(events)), nil
}
