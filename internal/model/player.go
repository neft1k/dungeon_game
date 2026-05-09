package model

import "time"

type PlayerStatus int

const (
	StatusActive PlayerStatus = iota
	StatusDead
	StatusDisqual
)

type Result int

const (
	ResultSuccess Result = iota
	ResultFail
	ResultDisqual
)

func (r Result) String() string {
	switch r {
	case ResultSuccess:
		return "SUCCESS"
	case ResultFail:
		return "FAIL"
	default:
		return "DISQUAL"
	}
}

type Player struct {
	ID           int
	HP           int
	Floor        int
	Status       PlayerStatus
	EnteredAt    time.Time
	LeftAt       time.Time
	FloorTimes   []time.Duration
	FloorEnter   time.Time
	BossEnter    time.Time
	BossKillAt   time.Time
	InDungeon    bool
	MonstersLeft map[int]int
	BossDefeated bool
}
