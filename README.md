# dungeon-events-processor

## What the application does
The app processes dungeon challenge events and produces a final report.

Main capabilities:
- reads dungeon config (`JSON`) and event log (`text`);
- validates player actions against challenge rules;
- tracks player state (registration, dungeon presence, floors, boss, HP);
- emits output event log and `Final report` (`SUCCESS` / `FAIL` / `DISQUAL`).

## Build
```bash
go build -o bin/dungeon ./cmd/dungeon
```

## Run
Default input files:
- `testdata/sample/config.json`
- `testdata/sample/events.txt`

Run with defaults:

```bash
./bin/dungeon
```

Run with custom files:

```bash
./bin/dungeon --config /path/to/config.json --events /path/to/events.txt
```

## Tests
Run all tests:

```bash
go test ./...
```

## Assumptions
- event timeline does not cross midnight due to HH:MM:SS format and non-decreasing time constraint.