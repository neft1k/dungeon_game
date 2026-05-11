package adapter

import "dungeon_game/internal/model"

type EventConsumer interface {
	Handle(e model.Event)
}
