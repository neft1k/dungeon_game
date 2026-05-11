package adapter

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"dungeon_game/internal/model"
)

func ParseEvents(r io.Reader, consumer EventConsumer) error {
	scanner := bufio.NewScanner(r)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		event, err := parseLine(line)
		if err != nil {
			return fmt.Errorf("line %d: %w", lineNum, err)
		}

		consumer.Handle(event)
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("reading input: %w", err)
	}

	return nil
}

func parseLine(line string) (model.Event, error) {
	if len(line) < 10 || line[0] != '[' {
		return model.Event{}, fmt.Errorf("invalid format: %q", line)
	}

	closeBracket := strings.Index(line, "]")
	if closeBracket == -1 {
		return model.Event{}, fmt.Errorf("missing closing bracket: %q", line)
	}

	timeStr := line[1:closeBracket]
	t, err := time.Parse("15:04:05", timeStr)
	if err != nil {
		return model.Event{}, fmt.Errorf("invalid time %q: %w", timeStr, err)
	}

	rest := strings.TrimSpace(line[closeBracket+1:])
	parts := strings.SplitN(rest, " ", 3)
	if len(parts) < 2 {
		return model.Event{}, fmt.Errorf("missing event id or player id: %q", line)
	}

	playerID, err := strconv.Atoi(parts[0])
	if err != nil {
		return model.Event{}, fmt.Errorf("invalid player id %q: %w", parts[0], err)
	}

	eventID, err := strconv.Atoi(parts[1])
	if err != nil {
		return model.Event{}, fmt.Errorf("invalid event id %q: %w", parts[1], err)
	}

	extra := ""
	if len(parts) == 3 {
		extra = parts[2]
	}

	return model.Event{
		Time:     t,
		ID:       model.EventID(eventID),
		PlayerID: playerID,
		Extra:    extra,
	}, nil
}
