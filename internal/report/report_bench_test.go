package report

import (
	"testing"
	"time"

	"dungeon_game/internal/model"
)

func makePlayers(n int) []*model.Player {
	players := make([]*model.Player, n)
	for i := range players {
		players[i] = &model.Player{
			ID:           i + 1,
			HP:           35,
			Status:       model.StatusActive,
			EnteredAt:    time.Time{}.Add(time.Duration(i) * time.Hour),
			LeftAt:       time.Time{}.Add(time.Duration(i)*time.Hour + 24*time.Minute),
			FloorTimes:   []time.Duration{5 * time.Minute, 19 * time.Minute},
			BossEnter:    time.Time{}.Add(time.Duration(i)*time.Hour + 8*time.Minute),
			BossKillAt:   time.Time{}.Add(time.Duration(i)*time.Hour + 19*time.Minute),
			BossDefeated: true,
			MonstersLeft: map[int]int{1: 0, 2: 0},
		}
	}
	return players
}

func BenchmarkBuild(b *testing.B) {
	players := makePlayers(1000)
	for b.Loop() {
		Build(players)
	}
}
