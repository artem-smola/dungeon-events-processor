package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/artem-smola/dungeon-events-processor/internal/app"
)

func main() {
	configPath := flag.String("config", "testdata/sample/config.json", "path to config json")
	eventsPath := flag.String("events", "testdata/sample/events.txt", "path to events file")
	flag.Parse()

	out, err := app.Run(*configPath, *eventsPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	for _, line := range out {
		fmt.Println(line)
	}
}
