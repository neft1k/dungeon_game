package server

import (
	"strings"
	"testing"

	"dungeon_game/internal/config"
)

const readmeEvents = `[14:00:00] 1 1
[14:00:00] 2 1
[14:10:00] 2 2
[14:10:00] 3 2
[14:11:00] 2 5
[14:12:00] 3 3
[14:14:00] 2 3
[14:27:00] 2 11 60
[14:29:00] 2 11 50
[14:40:00] 1 2
[14:41:00] 1 3
[14:44:00] 1 11 50
[14:45:00] 1 3
[14:48:00] 1 4
[14:48:00] 1 6
[14:49:00] 1 11 25
[14:49:02] 1 10 80
[14:50:00] 1 11 65
[14:59:00] 1 7
[15:04:00] 1 8
`

const readmeOutput = `[14:00:00] Player [1] registered
[14:00:00] Player [2] registered
[14:10:00] Player [2] entered the dungeon
[14:10:00] Player [3] is disqualified
[14:11:00] Player [2] makes imposible move [5]
[14:14:00] Player [2] killed the monster
[14:27:00] Player [2] recieved [60] of damage
[14:29:00] Player [2] recieved [50] of damage
[14:29:00] Player [2] is dead
[14:40:00] Player [1] entered the dungeon
[14:41:00] Player [1] killed the monster
[14:44:00] Player [1] recieved [50] of damage
[14:45:00] Player [1] killed the monster
[14:48:00] Player [1] went to the next floor
[14:48:00] Player [1] entered the boss's floor
[14:49:00] Player [1] recieved [25] of damage
[14:49:02] Player [1] has restored [80] of health
[14:50:00] Player [1] recieved [65] of damage
[14:59:00] Player [1] killed the boss
[15:04:00] Player [1] left the dungeon
Final report:
[SUCCESS] 1 [00:24:00, 00:05:00, 00:11:00] HP:35
[FAIL] 2 [00:19:00, 00:00:00, 00:00:00] HP:0
[DISQUAL] 3 [00:00:00, 00:00:00, 00:00:00] HP:100
`

const cannotContEvents = `[14:00:00] 1 1
[14:10:00] 1 2
[14:20:00] 1 11 50
[14:25:00] 1 9 out of potions
`

const cannotContOutput = `[14:00:00] Player [1] registered
[14:10:00] Player [1] entered the dungeon
[14:20:00] Player [1] recieved [50] of damage
[14:25:00] Player [1] cannot continue due to [out of potions]
Final report:
[DISQUAL] 1 [00:15:00, 00:00:00, 00:00:00] HP:50
`

func TestIntegration_CannotContinue(t *testing.T) {
	cfg := &config.Config{
		Floors:   2,
		Monsters: 2,
		OpenAt:   "14:05:00",
		Duration: 2,
	}

	var out strings.Builder
	if err := Start(cfg, strings.NewReader(cannotContEvents), &out); err != nil {
		t.Fatalf("неожиданная ошибка: %v", err)
	}

	got := out.String()
	if got != cannotContOutput {
		t.Errorf("вывод не совпадает\nполучено:\n%s\nожидалось:\n%s", got, cannotContOutput)
	}
}

func TestIntegration_ReadmeExample(t *testing.T) {
	cfg := &config.Config{
		Floors:   2,
		Monsters: 2,
		OpenAt:   "14:05:00",
		Duration: 2,
	}

	var out strings.Builder
	if err := Start(cfg, strings.NewReader(readmeEvents), &out); err != nil {
		t.Fatalf("неожиданная ошибка: %v", err)
	}

	got := out.String()
	if got != readmeOutput {
		t.Errorf("вывод не совпадает\nполучено:\n%s\nожидалось:\n%s", got, readmeOutput)
	}
}
