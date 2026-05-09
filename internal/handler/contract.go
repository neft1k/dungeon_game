package handler

import "dungeon_game/internal/model"

type Repository interface {
	RegisterPlayer(id int) *model.Player
	GetPlayer(id int) *model.Player
	ActiveInDungeon(id int) *model.Player
	AllPlayers() []*model.Player
	Dungeon() *model.Dungeon
}
