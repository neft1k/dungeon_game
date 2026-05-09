package server

import (
	"io"

	"dungeon_game/internal/adapter"
	"dungeon_game/internal/config"
	"dungeon_game/internal/handler"
	"dungeon_game/internal/model"
	"dungeon_game/internal/processor"
	"dungeon_game/internal/report"
	"dungeon_game/internal/repository"
)

func Start(cfg *config.Config, in io.Reader, out io.Writer) error {
	events, err := adapter.ParseEvents(in)
	if err != nil {
		return err
	}

	openAt, _ := cfg.OpenTime()
	closeAt, _ := cfg.CloseTime()

	repo := repository.New(&model.Dungeon{
		TotalFloors:      cfg.Floors,
		MonstersPerFloor: cfg.Monsters,
		BossFloor:        cfg.Floors,
	})

	h := handler.New(repo, openAt, closeAt)
	proc := processor.New(h, closeAt)

	players := proc.Run(events, out)
	adapter.FormatReport(report.Build(players), out)

	return nil
}
