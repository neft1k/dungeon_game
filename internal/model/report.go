package model

import "time"

type Report struct {
	Result    Result
	PlayerID  int
	TotalTime time.Duration
	AvgFloor  time.Duration
	BossTime  time.Duration
	HP        int
}
