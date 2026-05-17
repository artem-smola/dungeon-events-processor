package parser

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/artem-smola/dungeon-events-processor/internal/model"
)

func ReadConfig(path string) (model.Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return model.Config{}, fmt.Errorf("read config: %w", err)
	}

	var cfg model.Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return model.Config{}, fmt.Errorf("parse config: %w", err)
	}

	return cfg, nil
}
