package report

import (
	"time"

	"dungeon_game/internal/model"
)

func Build(players []*model.Player) []model.Report {
	reports := make([]model.Report, 0, len(players))
	for _, p := range players {
		reports = append(reports, build(p))
	}
	return reports
}

func build(p *model.Player) model.Report {
	return model.Report{
		Result:    resolveResult(p),
		PlayerID:  p.ID,
		TotalTime: totalTime(p),
		AvgFloor:  avgFloorTime(p),
		BossTime:  bossTime(p),
		HP:        p.HP,
	}
}

func resolveResult(p *model.Player) model.Result {
	switch p.Status {
	case model.StatusDisqual:
		return model.ResultDisqual
	case model.StatusDead:
		return model.ResultFail
	default:
		if isDungeonComplete(p) {
			return model.ResultSuccess
		}
		return model.ResultFail
	}
}

func isDungeonComplete(p *model.Player) bool {
	if !p.BossDefeated {
		return false
	}
	for _, count := range p.MonstersLeft {
		if count > 0 {
			return false
		}
	}
	return true
}

func totalTime(p *model.Player) time.Duration {
	if p.EnteredAt.IsZero() || p.LeftAt.IsZero() {
		return 0
	}
	return p.LeftAt.Sub(p.EnteredAt)
}

func avgFloorTime(p *model.Player) time.Duration {
	if len(p.FloorTimes) == 0 {
		return 0
	}
	var total time.Duration
	for _, t := range p.FloorTimes {
		total += t
	}
	return total / time.Duration(len(p.FloorTimes))
}

func bossTime(p *model.Player) time.Duration {
	if p.BossKillAt.IsZero() {
		return 0
	}
	return p.BossKillAt.Sub(p.BossEnter)
}
