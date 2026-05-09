package adapter

import (
	"strings"
	"testing"
)

var benchEvents = func() string {
	lines := []string{
		"[14:00:00] 1 1",
		"[14:00:00] 2 1",
		"[14:10:00] 2 2",
		"[14:14:00] 2 3",
		"[14:27:00] 2 11 60",
		"[14:29:00] 2 11 50",
		"[14:40:00] 1 2",
		"[14:41:00] 1 3",
		"[14:44:00] 1 11 50",
		"[14:48:00] 1 4",
		"[14:48:00] 1 6",
		"[14:59:00] 1 7",
		"[15:04:00] 1 8",
	}
	var sb strings.Builder
	for i := 0; i < 1000; i++ {
		for _, l := range lines {
			sb.WriteString(l)
			sb.WriteByte('\n')
		}
	}
	return sb.String()
}()

func BenchmarkParseEvents(b *testing.B) {
	for b.Loop() {
		_, _ = ParseEvents(strings.NewReader(benchEvents))
	}
}
