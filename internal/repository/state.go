package repository

import (
	"sort"

	"dungeon_game/internal/model"
)

type State struct {
	players map[int]*model.Player
	dungeon *model.Dungeon
}

func New(dungeon *model.Dungeon) *State {
	return &State{
		players: make(map[int]*model.Player),
		dungeon: dungeon,
	}
}

func (s *State) RegisterPlayer(id int) *model.Player {
	p := &model.Player{ID: id, HP: 100, Status: model.StatusActive}
	s.players[id] = p
	return p
}

func (s *State) GetPlayer(id int) *model.Player {
	return s.players[id]
}

func (s *State) ActiveInDungeon(id int) *model.Player {
	p := s.players[id]
	if p == nil || p.Status != model.StatusActive || !p.InDungeon {
		return nil
	}
	return p
}

func (s *State) AllPlayers() []*model.Player {
	players := make([]*model.Player, 0, len(s.players))
	for _, p := range s.players {
		players = append(players, p)
	}
	sort.Slice(players, func(i, j int) bool {
		return players[i].ID < players[j].ID
	})
	return players
}

func (s *State) Dungeon() *model.Dungeon {
	return s.dungeon
}
