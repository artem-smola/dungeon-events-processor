package engine

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/artem-smola/dungeon-events-processor/internal/model"
)

func (e *Engine) handleRegistration(p *playerState, ev model.Event) {
	if p.registered {
		e.impossible(p, ev)
		return
	}

	p.registered = true
	e.log(ev.TimeText, p.id, "registered")
}

func (e *Engine) handleEnter(p *playerState, ev model.Event) {
	if p.inDungeon {
		e.impossible(p, ev)
		return
	}
	if ev.TimeSec < e.openAtSec || ev.TimeSec >= e.closeAtSec {
		e.impossible(p, ev)
		return
	}

	p.inDungeon = true
	p.enteredAt = ev.TimeSec
	p.endedAt = ev.TimeSec
	p.currentFloor = 1
	p.onBoss = false
	p.bossKilled = false
	p.bossElapsed = 0
	p.bossStartedAt = -1
	p.bossTimerIsActive = false

	if e.regularFloors > 0 {
		e.tryAutoCompleteCurrentFloor(p)
		e.startFloorTimer(p, ev.TimeSec)
	}

	e.log(ev.TimeText, p.id, "entered the dungeon")
}

func (e *Engine) handleKillMonster(p *playerState, ev model.Event) {
	if !p.inDungeon || p.onBoss || p.currentFloor < 1 || p.currentFloor > e.regularFloors {
		e.impossible(p, ev)
		return
	}
	if p.floorCompleted[p.currentFloor] {
		e.impossible(p, ev)
		return
	}

	if p.floorKills[p.currentFloor] >= e.cfg.Monsters {
		e.impossible(p, ev)
		return
	}

	p.floorKills[p.currentFloor]++
	e.log(ev.TimeText, p.id, "killed the monster")

	if p.floorKills[p.currentFloor] == e.cfg.Monsters {
		e.stopFloorTimer(p, ev.TimeSec)
		p.floorCompleted[p.currentFloor] = true
		p.completedFloors++
		p.sumCompletedFloorTime += p.floorTimes[p.currentFloor]
	}
}

func (e *Engine) handleNextFloor(p *playerState, ev model.Event) {
	if !p.inDungeon || p.onBoss || p.currentFloor >= e.cfg.Floors {
		e.impossible(p, ev)
		return
	}

	e.stopFloorTimer(p, ev.TimeSec)
	p.currentFloor++
	e.tryAutoCompleteCurrentFloor(p)
	e.startFloorTimer(p, ev.TimeSec)
	e.log(ev.TimeText, p.id, "went to the next floor")
}

func (e *Engine) handlePrevFloor(p *playerState, ev model.Event) {
	if !p.inDungeon || p.currentFloor <= 1 {
		e.impossible(p, ev)
		return
	}

	e.stopFloorTimer(p, ev.TimeSec)
	e.stopBossTimer(p, ev.TimeSec)
	p.currentFloor--
	p.onBoss = false
	e.tryAutoCompleteCurrentFloor(p)
	e.startFloorTimer(p, ev.TimeSec)
	e.log(ev.TimeText, p.id, "went to the previous floor")
}

func (e *Engine) handleEnterBoss(p *playerState, ev model.Event) {
	if !p.inDungeon || p.onBoss {
		e.impossible(p, ev)
		return
	}
	if p.currentFloor != e.cfg.Floors {
		e.impossible(p, ev)
		return
	}

	e.stopFloorTimer(p, ev.TimeSec)
	p.onBoss = true

	if !p.bossKilled {
		e.startBossTimer(p, ev.TimeSec)
	}

	e.log(ev.TimeText, p.id, "entered the boss's floor")
}

func (e *Engine) handleKillBoss(p *playerState, ev model.Event) {
	if !p.inDungeon || !p.onBoss || p.bossKilled {
		e.impossible(p, ev)
		return
	}

	p.bossKilled = true
	e.stopBossTimer(p, ev.TimeSec)
	e.log(ev.TimeText, p.id, "killed the boss")
}

func (e *Engine) handleLeave(p *playerState, ev model.Event) {
	if !p.inDungeon {
		e.impossible(p, ev)
		return
	}

	e.endPlayer(p, "", ev.TimeSec)
	e.log(ev.TimeText, p.id, "left the dungeon")
}

func (e *Engine) handleCannotContinue(p *playerState, ev model.Event) {
	if !p.inDungeon {
		e.impossible(p, ev)
		return
	}

	reason := strings.TrimSpace(ev.Extra)
	if reason == "" {
		e.impossible(p, ev)
		return
	}

	e.endPlayer(p, model.StateDisqual, ev.TimeSec)
	e.log(ev.TimeText, p.id, fmt.Sprintf("cannot continue due to [%s]", reason))
}

func (e *Engine) handleRestore(p *playerState, ev model.Event) {
	if !p.inDungeon {
		e.impossible(p, ev)
		return
	}

	health, err := strconv.Atoi(strings.TrimSpace(ev.Extra))
	if err != nil || health < 0 {
		e.impossible(p, ev)
		return
	}

	p.health += health
	if p.health > maxHealth {
		p.health = maxHealth
	}
	e.log(ev.TimeText, p.id, fmt.Sprintf("has restored [%d] of health", health))
}

func (e *Engine) handleDamage(p *playerState, ev model.Event) {
	if !p.inDungeon {
		e.impossible(p, ev)
		return
	}

	damage, err := strconv.Atoi(strings.TrimSpace(ev.Extra))
	if err != nil || damage < 0 {
		e.impossible(p, ev)
		return
	}

	p.health -= damage
	if p.health < 0 {
		p.health = 0
	}
	e.log(ev.TimeText, p.id, fmt.Sprintf("recieved [%d] of damage", damage))

	if p.health == 0 {
		e.endPlayer(p, model.StateFail, ev.TimeSec)
		e.log(ev.TimeText, p.id, "is dead")
	}
}
