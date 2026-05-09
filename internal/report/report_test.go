package report

import (
	"testing"
	"time"

	"dungeon_game/internal/model"
)

type reportTestCase struct {
	name   string
	player *model.Player
	want   model.Result
}

type durationTestCase struct {
	name   string
	player *model.Player
	want   time.Duration
}

func mustTime(s string) time.Time {
	t, err := time.Parse("15:04:05", s)
	if err != nil {
		panic(err)
	}
	return t
}

func TestResolveResult(t *testing.T) {
	tests := []reportTestCase{
		{
			name:   "дисквалифицирован",
			player: &model.Player{Status: model.StatusDisqual},
			want:   model.ResultDisqual,
		},
		{
			name:   "мёртв",
			player: &model.Player{Status: model.StatusDead},
			want:   model.ResultFail,
		},
		{
			name: "данж пройден",
			player: &model.Player{
				Status:       model.StatusActive,
				BossDefeated: true,
				MonstersLeft: map[int]int{1: 0, 2: 0},
			},
			want: model.ResultSuccess,
		},
		{
			name: "босс не убит",
			player: &model.Player{
				Status:       model.StatusActive,
				BossDefeated: false,
				MonstersLeft: map[int]int{1: 0},
			},
			want: model.ResultFail,
		},
		{
			name: "остались монстры",
			player: &model.Player{
				Status:       model.StatusActive,
				BossDefeated: true,
				MonstersLeft: map[int]int{1: 2},
			},
			want: model.ResultFail,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := resolveResult(tt.player)
			if got != tt.want {
				t.Errorf("получено %v, ожидалось %v", got, tt.want)
			}
		})
	}
}

func TestTotalTime(t *testing.T) {
	tests := []durationTestCase{
		{
			name: "нормальный случай",
			player: &model.Player{
				EnteredAt: mustTime("14:00:00"),
				LeftAt:    mustTime("14:30:00"),
			},
			want: 30 * time.Minute,
		},
		{
			name:   "не входил в данж",
			player: &model.Player{},
			want:   0,
		},
		{
			name: "не вышел из данжа",
			player: &model.Player{
				EnteredAt: mustTime("14:00:00"),
			},
			want: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := totalTime(tt.player)
			if got != tt.want {
				t.Errorf("получено %v, ожидалось %v", got, tt.want)
			}
		})
	}
}

func TestAvgFloorTime(t *testing.T) {
	tests := []durationTestCase{
		{
			name: "два этажа",
			player: &model.Player{
				FloorTimes: []time.Duration{10 * time.Minute, 20 * time.Minute},
			},
			want: 15 * time.Minute,
		},
		{
			name: "один этаж",
			player: &model.Player{
				FloorTimes: []time.Duration{12 * time.Minute},
			},
			want: 12 * time.Minute,
		},
		{
			name:   "нет данных по этажам",
			player: &model.Player{},
			want:   0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := avgFloorTime(tt.player)
			if got != tt.want {
				t.Errorf("получено %v, ожидалось %v", got, tt.want)
			}
		})
	}
}

func TestBossTime(t *testing.T) {
	tests := []durationTestCase{
		{
			name: "босс убит",
			player: &model.Player{
				BossEnter:  mustTime("14:48:00"),
				BossKillAt: mustTime("14:59:00"),
			},
			want: 11 * time.Minute,
		},
		{
			name:   "босс не убит",
			player: &model.Player{},
			want:   0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := bossTime(tt.player)
			if got != tt.want {
				t.Errorf("получено %v, ожидалось %v", got, tt.want)
			}
		})
	}
}

func TestBuild(t *testing.T) {
	p := &model.Player{
		ID:           1,
		HP:           35,
		Status:       model.StatusActive,
		EnteredAt:    mustTime("14:40:00"),
		LeftAt:       mustTime("15:04:00"),
		FloorTimes:   []time.Duration{5 * time.Minute, 19 * time.Minute},
		BossEnter:    mustTime("14:48:00"),
		BossKillAt:   mustTime("14:59:00"),
		BossDefeated: true,
		MonstersLeft: map[int]int{1: 0, 2: 0},
	}

	reports := Build([]*model.Player{p})

	if len(reports) != 1 {
		t.Fatalf("ожидался 1 отчёт, получено %d", len(reports))
	}
	r := reports[0]

	if r.Result != model.ResultSuccess {
		t.Errorf("Result: получено %v, ожидалось SUCCESS", r.Result)
	}
	if r.PlayerID != 1 {
		t.Errorf("PlayerID: получено %d, ожидалось 1", r.PlayerID)
	}
	if r.TotalTime != 24*time.Minute {
		t.Errorf("TotalTime: получено %v, ожидалось 24m", r.TotalTime)
	}
	if r.AvgFloor != 12*time.Minute {
		t.Errorf("AvgFloor: получено %v, ожидалось 12m", r.AvgFloor)
	}
	if r.BossTime != 11*time.Minute {
		t.Errorf("BossTime: получено %v, ожидалось 11m", r.BossTime)
	}
	if r.HP != 35 {
		t.Errorf("HP: получено %d, ожидалось 35", r.HP)
	}
}
