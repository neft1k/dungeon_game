package processor

import (
	"io"
	"time"

	"dungeon_game/internal/model"
)

type EventHandler interface {
	Register(e model.Event, out io.Writer)
	EnterDungeon(e model.Event, out io.Writer)
	KillMonster(e model.Event, out io.Writer)
	NextFloor(e model.Event, out io.Writer)
	PrevFloor(e model.Event, out io.Writer)
	EnterBoss(e model.Event, out io.Writer)
	KillBoss(e model.Event, out io.Writer)
	LeaveDungeon(e model.Event, out io.Writer)
	CannotCont(e model.Event, out io.Writer)
	Heal(e model.Event, out io.Writer)
	Damage(e model.Event, out io.Writer)
	CloseDungeon(closeAt time.Time)
	Players() []*model.Player
}
