package processor

import (
	"io"
	"time"

	"dungeon_game/internal/model"
)

type Processor struct {
	h        EventHandler
	handlers map[model.EventID]func(model.Event, io.Writer)
	closeAt  time.Time
	closed   bool
}

func New(h EventHandler, closeAt time.Time) *Processor {
	return &Processor{
		h:       h,
		closeAt: closeAt,
		handlers: map[model.EventID]func(model.Event, io.Writer){
			model.EventRegister:     h.Register,
			model.EventEnterDungeon: h.EnterDungeon,
			model.EventKillMonster:  h.KillMonster,
			model.EventNextFloor:    h.NextFloor,
			model.EventPrevFloor:    h.PrevFloor,
			model.EventEnterBoss:    h.EnterBoss,
			model.EventKillBoss:     h.KillBoss,
			model.EventLeaveDungeon: h.LeaveDungeon,
			model.EventCannotCont:   h.CannotCont,
			model.EventHeal:         h.Heal,
			model.EventDamage:       h.Damage,
		},
	}
}

func (p *Processor) Run(events []model.Event, out io.Writer) []*model.Player {
	for _, e := range events {
		p.closeIfExpired(e.Time)

		h, ok := p.handlers[e.ID]
		if !ok {
			continue
		}
		h(e, out)
	}

	p.closeIfExpired(p.closeAt)

	return p.h.Players()
}

func (p *Processor) closeIfExpired(t time.Time) {
	if p.closed || t.Before(p.closeAt) {
		return
	}
	p.h.CloseDungeon(p.closeAt)
	p.closed = true
}
