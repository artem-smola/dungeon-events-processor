package engine

import "github.com/artem-smola/dungeon-events-processor/internal/model"

func (e *Engine) tryEndPlayers(timeSec int) {
	if timeSec < e.closeAtSec {
		return
	}

	for _, p := range e.players {
		if !p.inDungeon || p.ended {
			continue
		}
		e.endPlayer(p, "", e.closeAtSec)
	}
}

func (e *Engine) process(ev model.Event) {
	p := e.player(ev.PlayerID)

	if p.ended {
		return
	}

	if ev.EventID == model.EventRegister {
		e.handleRegistration(p, ev)
		return
	}

	if !p.registered {
		p.state = model.StateDisqual
		p.ended = true
		p.inDungeon = false
		p.endedAt = p.enteredAt
		e.log(ev.TimeText, p.id, "is disqualified")
		return
	}

	switch ev.EventID {
	case model.EventEnter:
		e.handleEnter(p, ev)
	case model.EventKillMob:
		e.handleKillMonster(p, ev)
	case model.EventNext:
		e.handleNextFloor(p, ev)
	case model.EventPrev:
		e.handlePrevFloor(p, ev)
	case model.EventBossIn:
		e.handleEnterBoss(p, ev)
	case model.EventKillBoss:
		e.handleKillBoss(p, ev)
	case model.EventLeave:
		e.handleLeave(p, ev)
	case model.EventCantGo:
		e.handleCannotContinue(p, ev)
	case model.EventHeal:
		e.handleRestore(p, ev)
	case model.EventDamage:
		e.handleDamage(p, ev)
	default:
		e.impossible(p, ev)
	}
}
