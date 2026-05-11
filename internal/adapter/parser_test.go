package adapter

import (
	"io"
	"strings"
	"testing"
	"time"

	"dungeon_game/internal/model"
)

func mustParseTime(s string) time.Time {
	t, err := time.Parse("15:04:05", s)
	if err != nil {
		panic(err)
	}
	return t
}

type parseTestCase struct {
	name    string
	input   string
	want    []model.Event
	wantErr bool
}

type eventCollector struct {
	events []model.Event
}

func (c *eventCollector) Handle(e model.Event) {
	c.events = append(c.events, e)
}

func collectEvents(r io.Reader) ([]model.Event, error) {
	c := &eventCollector{}
	err := ParseEvents(r, c)
	return c.events, err
}

func TestParseEvents(t *testing.T) {
	tests := []parseTestCase{
		{
			name:  "регистрация игрока",
			input: "[14:00:00] 1 1\n",
			want: []model.Event{
				{Time: mustParseTime("14:00:00"), ID: model.EventRegister, PlayerID: 1},
			},
		},
		{
			name:  "вход в данж",
			input: "[14:10:00] 2 2\n",
			want: []model.Event{
				{Time: mustParseTime("14:10:00"), ID: model.EventEnterDungeon, PlayerID: 2},
			},
		},
		{
			name:  "урон с параметром",
			input: "[14:27:00] 2 11 60\n",
			want: []model.Event{
				{Time: mustParseTime("14:27:00"), ID: model.EventDamage, PlayerID: 2, Extra: "60"},
			},
		},
		{
			name:  "cannot continue с несколькими словами",
			input: "[14:00:00] 1 9 out of potions\n",
			want: []model.Event{
				{Time: mustParseTime("14:00:00"), ID: model.EventCannotCont, PlayerID: 1, Extra: "out of potions"},
			},
		},
		{
			name:  "пустые строки пропускаются",
			input: "[14:00:00] 1 1\n\n[14:10:00] 2 2\n",
			want: []model.Event{
				{Time: mustParseTime("14:00:00"), ID: model.EventRegister, PlayerID: 1},
				{Time: mustParseTime("14:10:00"), ID: model.EventEnterDungeon, PlayerID: 2},
			},
		},
		{
			name:  "несколько событий подряд",
			input: "[14:00:00] 1 1\n[14:00:00] 2 1\n",
			want: []model.Event{
				{Time: mustParseTime("14:00:00"), ID: model.EventRegister, PlayerID: 1},
				{Time: mustParseTime("14:00:00"), ID: model.EventRegister, PlayerID: 2},
			},
		},
		{
			name:    "нет открывающей скобки",
			input:   "14:00:00] 1 1\n",
			wantErr: true,
		},
		{
			name:    "нет закрывающей скобки",
			input:   "[14:00:00 1 1\n",
			wantErr: true,
		},
		{
			name:    "неверный формат времени",
			input:   "[99:99:99] 1 1\n",
			wantErr: true,
		},
		{
			name:    "playerID не число",
			input:   "[14:00:00] abc 1\n",
			wantErr: true,
		},
		{
			name:    "eventID не число",
			input:   "[14:00:00] 1 abc\n",
			wantErr: true,
		},
		{
			name:    "слишком мало полей",
			input:   "[14:00:00] 1\n",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := collectEvents(strings.NewReader(tt.input))

			if tt.wantErr {
				if err == nil {
					t.Fatal("ожидалась ошибка, но её нет")
				}
				return
			}

			if err != nil {
				t.Fatalf("неожиданная ошибка: %v", err)
			}
			if len(got) != len(tt.want) {
				t.Fatalf("ожидалось %d событий, получено %d", len(tt.want), len(got))
			}
			for i, e := range got {
				if e.Time != tt.want[i].Time {
					t.Errorf("событие %d: Time = %v, хотели %v", i, e.Time, tt.want[i].Time)
				}
				if e.ID != tt.want[i].ID {
					t.Errorf("событие %d: ID = %v, хотели %v", i, e.ID, tt.want[i].ID)
				}
				if e.PlayerID != tt.want[i].PlayerID {
					t.Errorf("событие %d: PlayerID = %v, хотели %v", i, e.PlayerID, tt.want[i].PlayerID)
				}
				if e.Extra != tt.want[i].Extra {
					t.Errorf("событие %d: Extra = %q, хотели %q", i, e.Extra, tt.want[i].Extra)
				}
			}
		})
	}
}
