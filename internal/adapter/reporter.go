package adapter

import (
	"fmt"
	"io"
	"time"

	"dungeon_game/internal/model"
)

func FormatReport(reports []model.Report, out io.Writer) {
	_, _ = fmt.Fprintln(out, "Final report:")
	for _, r := range reports {
		_, _ = fmt.Fprintf(out, "[%s] %d [%s, %s, %s] HP:%d\n",
			r.Result,
			r.PlayerID,
			formatDuration(r.TotalTime),
			formatDuration(r.AvgFloor),
			formatDuration(r.BossTime),
			r.HP,
		)
	}
}

func formatDuration(d time.Duration) string {
	h := int(d.Hours())
	m := int(d.Minutes()) % 60
	s := int(d.Seconds()) % 60
	return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
}
