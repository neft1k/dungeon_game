package model

import "time"

type EventID int

const (
	EventRegister     EventID = 1
	EventEnterDungeon EventID = 2
	EventKillMonster  EventID = 3
	EventNextFloor    EventID = 4
	EventPrevFloor    EventID = 5
	EventEnterBoss    EventID = 6
	EventKillBoss     EventID = 7
	EventLeaveDungeon EventID = 8
	EventCannotCont   EventID = 9
	EventHeal         EventID = 10
	EventDamage       EventID = 11
)

type Event struct {
	Time     time.Time
	ID       EventID
	PlayerID int
	Extra    string
}

type OutEvent struct {
	Time     time.Time
	ID       EventID
	PlayerID int
	Extra    string
}
