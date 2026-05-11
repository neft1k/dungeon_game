package processor

import (
	"io"
	"strings"
	"testing"
	"time"

	"dungeon_game/internal/model"
)

type fakeHandler struct {
	received   []model.EventID
	closedAt   time.Time
	closeCount int
	players    []*model.Player
}

func (f *fakeHandler) Register(e model.Event, _ io.Writer)     { f.received = append(f.received, e.ID) }
func (f *fakeHandler) EnterDungeon(e model.Event, _ io.Writer) { f.received = append(f.received, e.ID) }
func (f *fakeHandler) KillMonster(e model.Event, _ io.Writer)  { f.received = append(f.received, e.ID) }
func (f *fakeHandler) NextFloor(e model.Event, _ io.Writer)    { f.received = append(f.received, e.ID) }
func (f *fakeHandler) PrevFloor(e model.Event, _ io.Writer)    { f.received = append(f.received, e.ID) }
func (f *fakeHandler) EnterBoss(e model.Event, _ io.Writer)    { f.received = append(f.received, e.ID) }
func (f *fakeHandler) KillBoss(e model.Event, _ io.Writer)     { f.received = append(f.received, e.ID) }
func (f *fakeHandler) LeaveDungeon(e model.Event, _ io.Writer) { f.received = append(f.received, e.ID) }
func (f *fakeHandler) CannotCont(e model.Event, _ io.Writer)   { f.received = append(f.received, e.ID) }
func (f *fakeHandler) Heal(e model.Event, _ io.Writer)         { f.received = append(f.received, e.ID) }
func (f *fakeHandler) Damage(e model.Event, _ io.Writer)       { f.received = append(f.received, e.ID) }
func (f *fakeHandler) CloseDungeon(t time.Time)                { f.closedAt = t; f.closeCount++ }
func (f *fakeHandler) Players() []*model.Player                { return f.players }

func mustTime(s string) time.Time {
	t, err := time.Parse("15:04:05", s)
	if err != nil {
		panic(err)
	}
	return t
}

func newProc(closeAt string) (*Processor, *fakeHandler) {
	h := &fakeHandler{}
	return New(h, mustTime(closeAt), &strings.Builder{}), h
}

func TestHandle_DispatchesKnownEvents(t *testing.T) {
	ids := []model.EventID{
		model.EventRegister, model.EventEnterDungeon, model.EventKillMonster,
		model.EventNextFloor, model.EventPrevFloor, model.EventEnterBoss,
		model.EventKillBoss, model.EventLeaveDungeon, model.EventCannotCont,
		model.EventHeal, model.EventDamage,
	}

	proc, h := newProc("16:00:00")
	for _, id := range ids {
		proc.Handle(model.Event{Time: mustTime("14:00:00"), ID: id})
	}

	if len(h.received) != len(ids) {
		t.Fatalf("ожидалось %d вызовов, получено %d", len(ids), len(h.received))
	}
	for i, id := range ids {
		if h.received[i] != id {
			t.Errorf("вызов %d: ожидался EventID=%d, получен %d", i, id, h.received[i])
		}
	}
}

func TestHandle_IgnoresUnknownEvent(t *testing.T) {
	proc, h := newProc("16:00:00")
	proc.Handle(model.Event{Time: mustTime("14:00:00"), ID: model.EventID(99)})
	if len(h.received) != 0 {
		t.Errorf("неизвестное событие не должно вызывать хендлер")
	}
}

func TestHandle_CloseDungeonOnExpiry(t *testing.T) {
	proc, h := newProc("15:00:00")
	proc.Handle(model.Event{Time: mustTime("15:00:00"), ID: model.EventRegister})
	if h.closeCount != 1 {
		t.Errorf("CloseDungeon должен вызваться один раз, вызвался %d раз", h.closeCount)
	}
	if h.closedAt != mustTime("15:00:00") {
		t.Errorf("CloseDungeon вызван с %v, ожидалось 15:00:00", h.closedAt)
	}
}

func TestHandle_NoCloseBeforeExpiry(t *testing.T) {
	proc, h := newProc("16:00:00")
	proc.Handle(model.Event{Time: mustTime("14:00:00"), ID: model.EventRegister})
	if h.closeCount != 0 {
		t.Errorf("CloseDungeon не должен вызываться до истечения времени")
	}
}

func TestHandle_ClosesOnlyOnce(t *testing.T) {
	proc, h := newProc("15:00:00")
	proc.Handle(model.Event{Time: mustTime("15:30:00"), ID: model.EventRegister})
	proc.Handle(model.Event{Time: mustTime("15:45:00"), ID: model.EventRegister})
	if h.closeCount != 1 {
		t.Errorf("CloseDungeon должен вызваться ровно один раз, вызвался %d раз", h.closeCount)
	}
}

func TestFinish_TriggersCloseIfNotYetClosed(t *testing.T) {
	proc, h := newProc("15:00:00")
	proc.Finish()
	if h.closeCount != 1 {
		t.Errorf("Finish должен вызвать CloseDungeon, вызвался %d раз", h.closeCount)
	}
}

func TestFinish_DoesNotDoubleClose(t *testing.T) {
	proc, h := newProc("15:00:00")
	proc.Handle(model.Event{Time: mustTime("15:30:00"), ID: model.EventRegister})
	proc.Finish()
	if h.closeCount != 1 {
		t.Errorf("CloseDungeon должен вызваться ровно один раз, вызвался %d раз", h.closeCount)
	}
}

func TestFinish_ReturnsPlayers(t *testing.T) {
	proc, h := newProc("16:00:00")
	h.players = []*model.Player{{ID: 1}, {ID: 2}}
	got := proc.Finish()
	if len(got) != 2 {
		t.Errorf("ожидалось 2 игрока, получено %d", len(got))
	}
}
