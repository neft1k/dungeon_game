package config

import (
	"os"
	"testing"
	"time"
)

type configTestCase struct {
	name    string
	cfg     Config
	wantErr bool
}

type openTimeTestCase struct {
	name    string
	openAt  string
	want    time.Time
	wantErr bool
}

func mustParseTime(s string) time.Time {
	t, err := time.Parse("15:04:05", s)
	if err != nil {
		panic(err)
	}
	return t
}

func TestOpenTime(t *testing.T) {
	tests := []openTimeTestCase{
		{
			name:   "корректное время",
			openAt: "14:05:00",
			want:   mustParseTime("14:05:00"),
		},
		{
			name:    "неверный формат",
			openAt:  "25:00:00",
			wantErr: true,
		},
		{
			name:    "пустая строка",
			openAt:  "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := Config{OpenAt: tt.openAt}
			got, err := cfg.OpenTime()
			if tt.wantErr {
				if err == nil {
					t.Fatal("ожидалась ошибка, но её нет")
				}
				return
			}
			if err != nil {
				t.Fatalf("неожиданная ошибка: %v", err)
			}
			if got != tt.want {
				t.Errorf("получено %v, ожидалось %v", got, tt.want)
			}
		})
	}
}

func TestCloseTime(t *testing.T) {
	cfg := Config{OpenAt: "14:00:00", Duration: 2}
	got, err := cfg.CloseTime()
	if err != nil {
		t.Fatalf("неожиданная ошибка: %v", err)
	}
	want := mustParseTime("16:00:00")
	if got != want {
		t.Errorf("получено %v, ожидалось %v", got, want)
	}
}

func TestValidate(t *testing.T) {
	tests := []configTestCase{
		{
			name: "корректный конфиг",
			cfg:  Config{Floors: 2, Monsters: 3, OpenAt: "14:00:00", Duration: 2},
		},
		{
			name:    "Floors = 0",
			cfg:     Config{Floors: 0, Monsters: 3, OpenAt: "14:00:00", Duration: 2},
			wantErr: true,
		},
		{
			name:    "Monsters = 0",
			cfg:     Config{Floors: 2, Monsters: 0, OpenAt: "14:00:00", Duration: 2},
			wantErr: true,
		},
		{
			name:    "Duration = 0",
			cfg:     Config{Floors: 2, Monsters: 3, OpenAt: "14:00:00", Duration: 0},
			wantErr: true,
		},
		{
			name:    "неверный OpenAt",
			cfg:     Config{Floors: 2, Monsters: 3, OpenAt: "bad", Duration: 2},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cfg.validate()
			if tt.wantErr {
				if err == nil {
					t.Fatal("ожидалась ошибка, но её нет")
				}
				return
			}
			if err != nil {
				t.Fatalf("неожиданная ошибка: %v", err)
			}
		})
	}
}

func writeTempConfig(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp("", "config*.json")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Remove(f.Name()) })
	if _, err := f.WriteString(content); err != nil {
		t.Fatal(err)
	}
	if err := f.Close(); err != nil {
		t.Fatal(err)
	}
	return f.Name()
}

func TestLoad(t *testing.T) {
	t.Run("корректный файл", func(t *testing.T) {
		path := writeTempConfig(t, `{"Floors":2,"Monsters":3,"OpenAt":"14:00:00","Duration":2}`)
		cfg, err := Load(path)
		if err != nil {
			t.Fatalf("неожиданная ошибка: %v", err)
		}
		if cfg.Floors != 2 || cfg.Monsters != 3 || cfg.Duration != 2 {
			t.Errorf("неверные поля: %+v", cfg)
		}
	})

	t.Run("файл не существует", func(t *testing.T) {
		_, err := Load("nonexistent.json")
		if err == nil {
			t.Fatal("ожидалась ошибка, но её нет")
		}
	})

	t.Run("невалидный json", func(t *testing.T) {
		path := writeTempConfig(t, `not json`)
		_, err := Load(path)
		if err == nil {
			t.Fatal("ожидалась ошибка, но её нет")
		}
	})

	t.Run("невалидные поля", func(t *testing.T) {
		path := writeTempConfig(t, `{"Floors":0,"Monsters":3,"OpenAt":"14:00:00","Duration":2}`)
		_, err := Load(path)
		if err == nil {
			t.Fatal("ожидалась ошибка, но её нет")
		}
	})
}
