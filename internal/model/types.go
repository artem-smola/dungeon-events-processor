package model

type Config struct {
	Floors   int    `json:"Floors"`
	Monsters int    `json:"Monsters"`
	OpenAt   string `json:"OpenAt"`
	Duration int    `json:"Duration"`
}

type Event struct {
	TimeSec  int
	TimeText string
	PlayerID int
	EventID  EventID
	Extra    string
}

type LogEntry struct {
	TimeText string
	PlayerID int
	Message  string
}

type PlayerReport struct {
	State           ChallengeState
	PlayerID        int
	TotalSeconds    int
	AvgFloorSeconds int
	BossSeconds     int
	Health          int
}

type Result struct {
	Logs    []LogEntry
	Reports []PlayerReport
}
