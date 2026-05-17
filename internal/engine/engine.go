package engine

import (
	"errors"
	"fmt"
	"sort"

	"github.com/artem-smola/dungeon-events-processor/internal/model"
	"github.com/artem-smola/dungeon-events-processor/internal/parser"
)

const (
	maxHealth  = 100
	daySeconds = 24 * 60 * 60
)

type Engine struct {
	cfg           model.Config
	openAtSec     int
	closeAtSec    int
	regularFloors int

	players map[int]*playerState
	logs    []model.LogEntry
}

type playerState struct {
	id int

	registered bool
	inDungeon  bool
	ended      bool
	state      model.ChallengeState

	health int

	enteredAt int
	endedAt   int

	currentFloor int
	onBoss       bool
	bossKilled   bool
	bossElapsed  int

	bossStartedAt     int
	bossTimerIsActive bool

	floorKills         []int
	floorCompleted     []bool
	floorTimes         []int
	floorStartedAt     int
	floorTimerIsActive bool

	completedFloors       int
	sumCompletedFloorTime int
}

func New(cfg model.Config) (*Engine, error) {
	if cfg.Floors < 1 {
		return nil, errors.New("Floors must be greater than zero")
	}
	if cfg.Monsters < 0 {
		return nil, errors.New("Monsters must be non-negative")
	}
	if cfg.Duration < 1 {
		return nil, errors.New("Duration must be greater than zero")
	}

	openAtSec, err := parser.ParseTime(cfg.OpenAt)
	if err != nil {
		return nil, fmt.Errorf("invalid OpenAt: %w", err)
	}

	closeAtSec := openAtSec + cfg.Duration*3600
	if closeAtSec > daySeconds {
		return nil, errors.New("OpenAt + Duration crosses midnight, which is unsupported")
	}

	return &Engine{
		cfg:           cfg,
		openAtSec:     openAtSec,
		closeAtSec:    closeAtSec,
		regularFloors: cfg.Floors - 1,
	}, nil
}

func (e *Engine) Run(events []model.Event) model.Result {
	e.resetRuntimeState()

	for _, ev := range events {
		e.tryEndPlayers(ev.TimeSec)
		e.process(ev)
	}
	e.tryEndPlayers(e.closeAtSec)

	ids := make([]int, 0, len(e.players))
	for id := range e.players {
		ids = append(ids, id)
	}
	sort.Ints(ids)

	reports := make([]model.PlayerReport, 0, len(ids))
	for _, id := range ids {
		p := e.players[id]
		reports = append(reports, model.PlayerReport{
			State:           e.resolveState(p),
			PlayerID:        p.id,
			TotalSeconds:    e.totalTime(p),
			AvgFloorSeconds: e.avgFloorTime(p),
			BossSeconds:     e.bossTime(p),
			Health:          p.health,
		})
	}

	return model.Result{Logs: e.logs, Reports: reports}
}

func (e *Engine) resetRuntimeState() {
	e.players = make(map[int]*playerState)
	e.logs = nil
}

func (e *Engine) resolveState(p *playerState) model.ChallengeState {
	if p.state != "" {
		return p.state
	}
	if p.bossKilled && p.completedFloors == e.regularFloors {
		return model.StateSuccess
	}
	return model.StateFail
}

func (e *Engine) totalTime(p *playerState) int {
	if p.enteredAt < 0 || p.endedAt < p.enteredAt {
		return 0
	}
	return p.endedAt - p.enteredAt
}

func (e *Engine) avgFloorTime(p *playerState) int {
	if p.completedFloors == 0 {
		return 0
	}
	return p.sumCompletedFloorTime / p.completedFloors
}

func (e *Engine) bossTime(p *playerState) int {
	if !p.bossKilled {
		return 0
	}
	return p.bossElapsed
}
