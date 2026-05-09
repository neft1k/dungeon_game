package handler

import (
	"strings"
	"testing"
	"time"

	"dungeon_game/internal/model"
)

type fakeRepo struct {
	players map[int]*model.Player
	dungeon *model.Dungeon
}

func newFakeRepo() *fakeRepo {
	return &fakeRepo{
		players: make(map[int]*model.Player),
		dungeon: &model.Dungeon{
			TotalFloors:      2,
			MonstersPerFloor: 2,
			BossFloor:        2,
		},
	}
}

func (r *fakeRepo) RegisterPlayer(id int) *model.Player {
	p := &model.Player{ID: id, HP: 100, Status: model.StatusActive}
	r.players[id] = p
	return p
}

func (r *fakeRepo) GetPlayer(id int) *model.Player { return r.players[id] }
func (r *fakeRepo) Dungeon() *model.Dungeon        { return r.dungeon }
func (r *fakeRepo) AllPlayers() []*model.Player {
	players := make([]*model.Player, 0, len(r.players))
	for _, p := range r.players {
		players = append(players, p)
	}
	return players
}
func (r *fakeRepo) ActiveInDungeon(id int) *model.Player {
	p := r.players[id]
	if p == nil || p.Status != model.StatusActive || !p.InDungeon {
		return nil
	}
	return p
}

func TestRegister(t *testing.T) {
	h, _ := newHandler()
	b := out()
	h.Register(model.Event{Time: mustTime("14:00:00"), ID: model.EventRegister, PlayerID: 1}, b)
	if !strings.Contains(b.String(), "registered") {
		t.Errorf("получено: %q", b.String())
	}
}

func TestEnterDungeon_Unregistered(t *testing.T) {
	h, repo := newHandler()
	b := out()
	h.EnterDungeon(model.Event{Time: mustTime("14:30:00"), ID: model.EventEnterDungeon, PlayerID: 1}, b)
	if repo.GetPlayer(1).Status != model.StatusDisqual {
		t.Error("незарегистрированный игрок должен быть дисквалифицирован")
	}
}

func TestEnterDungeon_Valid(t *testing.T) {
	h, repo := newHandler()
	repo.RegisterPlayer(1)
	b := out()
	h.EnterDungeon(model.Event{Time: mustTime("14:30:00"), ID: model.EventEnterDungeon, PlayerID: 1}, b)
	p := repo.GetPlayer(1)
	if !p.InDungeon || p.Floor != 1 {
		t.Errorf("ожидался вход в данж на 1 этаже, получено InDungeon=%v Floor=%d", p.InDungeon, p.Floor)
	}
}

func TestKillMonster(t *testing.T) {
	h, repo := newHandler()
	p := setupInDungeon(h, repo)
	before := p.MonstersLeft[1]
	h.KillMonster(model.Event{Time: mustTime("14:30:00"), ID: model.EventKillMonster, PlayerID: 1}, out())
	if p.MonstersLeft[1] != before-1 {
		t.Errorf("монстров должно стать %d, получено %d", before-1, p.MonstersLeft[1])
	}
}

func TestNextFloor(t *testing.T) {
	h, repo := newHandler()
	p := setupInDungeon(h, repo)
	h.NextFloor(model.Event{Time: mustTime("14:30:00"), ID: model.EventNextFloor, PlayerID: 1}, out())
	if p.Floor != 2 {
		t.Errorf("этаж должен стать 2, получено %d", p.Floor)
	}
}

func TestPrevFloor_OnFirstFloor(t *testing.T) {
	h, repo := newHandler()
	setupInDungeon(h, repo)
	b := out()
	h.PrevFloor(model.Event{Time: mustTime("14:30:00"), ID: model.EventPrevFloor, PlayerID: 1}, b)
	if !strings.Contains(b.String(), "makes imposible move") {
		t.Errorf("с 1 этажа нельзя идти назад, получено: %q", b.String())
	}
}

func TestEnterBoss(t *testing.T) {
	h, repo := newHandler()
	p := setupInDungeon(h, repo)
	p.Floor = 2
	h.EnterBoss(model.Event{Time: mustTime("14:30:00"), ID: model.EventEnterBoss, PlayerID: 1}, out())
	if p.BossEnter.IsZero() {
		t.Error("BossEnter должен быть установлен")
	}
}

func TestKillBoss(t *testing.T) {
	h, repo := newHandler()
	p := setupInDungeon(h, repo)
	p.Floor = 2
	p.BossEnter = mustTime("14:28:00")
	h.KillBoss(model.Event{Time: mustTime("14:30:00"), ID: model.EventKillBoss, PlayerID: 1}, out())
	if !p.BossDefeated {
		t.Error("BossDefeated должен быть true")
	}
}

func TestLeaveDungeon(t *testing.T) {
	h, repo := newHandler()
	p := setupInDungeon(h, repo)
	h.LeaveDungeon(model.Event{Time: mustTime("14:30:00"), ID: model.EventLeaveDungeon, PlayerID: 1}, out())
	if p.InDungeon || p.LeftAt.IsZero() {
		t.Error("игрок должен покинуть данж с установленным LeftAt")
	}
}

func TestCannotCont(t *testing.T) {
	h, repo := newHandler()
	repo.RegisterPlayer(1)
	b := out()
	h.CannotCont(model.Event{Time: mustTime("14:30:00"), ID: model.EventCannotCont, PlayerID: 1, Extra: "out of potions"}, b)
	if repo.GetPlayer(1).Status != model.StatusDisqual {
		t.Error("игрок должен быть дисквалифицирован")
	}
}

func TestHeal_CapAt100(t *testing.T) {
	h, repo := newHandler()
	p := setupInDungeon(h, repo)
	p.HP = 80
	h.Heal(model.Event{Time: mustTime("14:30:00"), ID: model.EventHeal, PlayerID: 1, Extra: "80"}, out())
	if p.HP != 100 {
		t.Errorf("HP не должно превышать 100, получено %d", p.HP)
	}
}

func TestDamage_PlayerDies(t *testing.T) {
	h, repo := newHandler()
	p := setupInDungeon(h, repo)
	p.HP = 30
	h.Damage(model.Event{Time: mustTime("14:30:00"), ID: model.EventDamage, PlayerID: 1, Extra: "50"}, out())
	if p.Status != model.StatusDead || p.InDungeon || p.HP != 0 {
		t.Errorf("игрок должен умереть: status=%v InDungeon=%v HP=%d", p.Status, p.InDungeon, p.HP)
	}
}

func TestCloseDungeon(t *testing.T) {
	h, repo := newHandler()
	p := setupInDungeon(h, repo)
	closeAt := mustTime("16:00:00")
	h.CloseDungeon(closeAt)
	if p.InDungeon || p.LeftAt != closeAt {
		t.Errorf("игрок должен быть выведен при закрытии: InDungeon=%v LeftAt=%v", p.InDungeon, p.LeftAt)
	}
}

func mustTime(s string) time.Time {
	t, err := time.Parse("15:04:05", s)
	if err != nil {
		panic(err)
	}
	return t
}

func newHandler() (*Handler, *fakeRepo) {
	repo := newFakeRepo()
	return New(repo, mustTime("14:00:00"), mustTime("16:00:00")), repo
}

func setupInDungeon(h *Handler, repo *fakeRepo) *model.Player {
	repo.RegisterPlayer(1)
	var out strings.Builder
	h.EnterDungeon(model.Event{Time: mustTime("14:30:00"), ID: model.EventEnterDungeon, PlayerID: 1}, &out)
	return repo.GetPlayer(1)
}

func out() *strings.Builder { return &strings.Builder{} }
