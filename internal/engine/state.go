package engine

import (
	"fmt"

	"github.com/artem-smola/dungeon-events-processor/internal/model"
)

func (e *Engine) startFloorTimer(p *playerState, timeSec int) {
	if p.currentFloor < 1 || p.currentFloor > e.regularFloors {
		return
	}
	if p.floorCompleted[p.currentFloor] {
		return
	}

	p.floorStartedAt = timeSec
	p.floorTimerIsActive = true
}

func (e *Engine) tryAutoCompleteCurrentFloor(p *playerState) {
	if e.cfg.Monsters != 0 {
		return
	}
	if p.currentFloor < 1 || p.currentFloor > e.regularFloors {
		return
	}
	if p.floorCompleted[p.currentFloor] {
		return
	}

	p.floorCompleted[p.currentFloor] = true
	p.completedFloors++
}

func (e *Engine) stopFloorTimer(p *playerState, timeSec int) {
	if !p.floorTimerIsActive {
		return
	}
	if timeSec < p.floorStartedAt {
		return
	}

	p.floorTimes[p.currentFloor] += timeSec - p.floorStartedAt
	p.floorTimerIsActive = false
}

func (e *Engine) startBossTimer(p *playerState, timeSec int) {
	if !p.onBoss || p.bossTimerIsActive {
		return
	}
	p.bossStartedAt = timeSec
	p.bossTimerIsActive = true
}

func (e *Engine) stopBossTimer(p *playerState, timeSec int) {
	if !p.bossTimerIsActive {
		return
	}
	if timeSec < p.bossStartedAt {
		return
	}
	p.bossElapsed += timeSec - p.bossStartedAt
	p.bossTimerIsActive = false
}

func (e *Engine) endPlayer(p *playerState, forcedState model.ChallengeState, timeSec int) {
	if p.ended {
		return
	}
	if timeSec < p.enteredAt {
		timeSec = p.enteredAt
	}
	e.stopBossTimer(p, timeSec)

	p.inDungeon = false
	p.ended = true
	p.endedAt = timeSec

	if forcedState != "" {
		p.state = forcedState
		return
	}

	if p.bossKilled && p.completedFloors == e.regularFloors {
		p.state = model.StateSuccess
	} else {
		p.state = model.StateFail
	}
}

func (e *Engine) impossible(p *playerState, ev model.Event) {
	e.log(ev.TimeText, p.id, fmt.Sprintf("makes imposible move [%d]", ev.EventID))
}

func (e *Engine) log(timeText string, playerID int, message string) {
	e.logs = append(e.logs, model.LogEntry{
		TimeText: timeText,
		PlayerID: playerID,
		Message:  message,
	})
}

func (e *Engine) player(id int) *playerState {
	if p, ok := e.players[id]; ok {
		return p
	}

	p := &playerState{
		id:             id,
		health:         maxHealth,
		enteredAt:      -1,
		endedAt:        -1,
		bossStartedAt:  -1,
		floorKills:     make([]int, e.regularFloors+1),
		floorCompleted: make([]bool, e.regularFloors+1),
		floorTimes:     make([]int, e.regularFloors+1),
	}
	e.players[id] = p
	return p
}
