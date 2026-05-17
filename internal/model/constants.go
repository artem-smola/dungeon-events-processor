package model

type EventID int

const (
	EventRegister EventID = 1
	EventEnter    EventID = 2
	EventKillMob  EventID = 3
	EventNext     EventID = 4
	EventPrev     EventID = 5
	EventBossIn   EventID = 6
	EventKillBoss EventID = 7
	EventLeave    EventID = 8
	EventCantGo   EventID = 9
	EventHeal     EventID = 10
	EventDamage   EventID = 11
)

type ChallengeState string

const (
	StateSuccess ChallengeState = "SUCCESS"
	StateFail    ChallengeState = "FAIL"
	StateDisqual ChallengeState = "DISQUAL"
)
