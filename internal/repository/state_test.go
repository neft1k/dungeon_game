package repository

import (
	"testing"

	"dungeon_game/internal/model"
)

func newState() *State {
	return New(&model.Dungeon{
		TotalFloors:      2,
		MonstersPerFloor: 2,
		BossFloor:        2,
	})
}

func TestRegisterPlayer(t *testing.T) {
	s := newState()
	p := s.RegisterPlayer(1)
	if p == nil || p.ID != 1 || p.HP != 100 || p.Status != model.StatusActive {
		t.Errorf("неверный начальный стейт игрока: %+v", p)
	}
}

func TestGetPlayer(t *testing.T) {
	s := newState()
	if s.GetPlayer(1) != nil {
		t.Error("несуществующий игрок должен возвращать nil")
	}
	s.RegisterPlayer(1)
	if s.GetPlayer(1) == nil {
		t.Error("зарегистрированный игрок должен возвращаться")
	}
}

func TestActiveInDungeon(t *testing.T) {
	s := newState()
	s.RegisterPlayer(1)
	p := s.GetPlayer(1)

	if s.ActiveInDungeon(1) != nil {
		t.Error("игрок не в данже — должен возвращать nil")
	}

	p.InDungeon = true
	if s.ActiveInDungeon(1) == nil {
		t.Error("активный игрок в данже должен возвращаться")
	}

	p.Status = model.StatusDead
	if s.ActiveInDungeon(1) != nil {
		t.Error("мёртвый игрок не должен возвращаться")
	}
}

func TestAllPlayers_SortedByID(t *testing.T) {
	s := newState()
	s.RegisterPlayer(3)
	s.RegisterPlayer(1)
	s.RegisterPlayer(2)

	players := s.AllPlayers()
	if len(players) != 3 {
		t.Fatalf("ожидалось 3 игрока, получено %d", len(players))
	}
	for i, p := range players {
		if p.ID != i+1 {
			t.Errorf("позиция %d: ожидался ID %d, получено %d", i, i+1, p.ID)
		}
	}
}

func TestDungeon(t *testing.T) {
	s := newState()
	d := s.Dungeon()
	if d.TotalFloors != 2 || d.BossFloor != 2 {
		t.Errorf("неверный данж: %+v", d)
	}
}
