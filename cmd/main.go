package main

import (
	"fmt"
	"os"

	"dungeon_game/internal/adapter"
	"dungeon_game/internal/config"
	"dungeon_game/internal/server"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: dungeon_game <config.json> [events]")
		os.Exit(1)
	}

	cfg, err := config.Load(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "config error: %v\n", err)
		os.Exit(1)
	}

	in, closeIn, err := adapter.OpenInput(os.Args[2:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer closeIn()

	if err := server.Start(cfg, in, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
