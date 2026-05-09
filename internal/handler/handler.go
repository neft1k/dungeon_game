package handler

import (
	"fmt"
	"io"
	"strconv"
	"time"

	"dungeon_game/internal/model"
)

type Handler struct {
	repo    Repository
	openAt  time.Time
	closeAt time.Time
}

func New(repo Repository, openAt, closeAt time.Time) *Handler {
	return &Handler{repo: repo, openAt: openAt, closeAt: closeAt}
}

func (h *Handler) Register(e model.Event, out io.Writer) {
	if h.repo.GetPlayer(e.PlayerID) != nil {
		emit(out, e.Time, e.PlayerID, fmt.Sprintf("makes imposible move [%d]", e.ID))
		return
	}
	h.repo.RegisterPlayer(e.PlayerID)
	emit(out, e.Time, e.PlayerID, "registered")
}

func (h *Handler) EnterDungeon(e model.Event, out io.Writer) {
	player := h.repo.GetPlayer(e.PlayerID)
	if player == nil {
		p := h.repo.RegisterPlayer(e.PlayerID)
		p.Status = model.StatusDisqual
		emit(out, e.Time, e.PlayerID, "is disqualified")
		return
	}
	if player.Status != model.StatusActive {
		return
	}
	if e.Time.Before(h.openAt) || !e.Time.Before(h.closeAt) {
		player.Status = model.StatusDisqual
		emit(out, e.Time, e.PlayerID, "is disqualified")
		return
	}
	if player.InDungeon || !player.LeftAt.IsZero() {
		emit(out, e.Time, e.PlayerID, fmt.Sprintf("makes imposible move [%d]", e.ID))
		return
	}

	player.InDungeon = true
	player.EnteredAt = e.Time
	player.Floor = 1
	player.FloorEnter = e.Time
	player.MonstersLeft = initMonstersLeft(h.repo.Dungeon())
	player.BossDefeated = false

	emit(out, e.Time, e.PlayerID, "entered the dungeon")
}

func (h *Handler) KillMonster(e model.Event, out io.Writer) {
	player := h.repo.ActiveInDungeon(e.PlayerID)
	if player == nil {
		return
	}

	dungeon := h.repo.Dungeon()
	if player.Floor == dungeon.BossFloor || player.MonstersLeft[player.Floor] <= 0 {
		emit(out, e.Time, e.PlayerID, fmt.Sprintf("makes imposible move [%d]", e.ID))
		return
	}

	player.MonstersLeft[player.Floor]--
	emit(out, e.Time, e.PlayerID, "killed the monster")

	if player.MonstersLeft[player.Floor] == 0 {
		player.FloorTimes = append(player.FloorTimes, e.Time.Sub(player.FloorEnter))
	}
}

func (h *Handler) NextFloor(e model.Event, out io.Writer) {
	player := h.repo.ActiveInDungeon(e.PlayerID)
	if player == nil {
		return
	}

	if player.Floor == h.repo.Dungeon().BossFloor {
		emit(out, e.Time, e.PlayerID, fmt.Sprintf("makes imposible move [%d]", e.ID))
		return
	}

	player.Floor++
	player.FloorEnter = e.Time
	emit(out, e.Time, e.PlayerID, "went to the next floor")
}

func (h *Handler) PrevFloor(e model.Event, out io.Writer) {
	player := h.repo.ActiveInDungeon(e.PlayerID)
	if player == nil {
		return
	}

	if player.Floor == 1 {
		emit(out, e.Time, e.PlayerID, fmt.Sprintf("makes imposible move [%d]", e.ID))
		return
	}

	player.Floor--
	player.FloorEnter = e.Time
	emit(out, e.Time, e.PlayerID, "went to the previous floor")
}

func (h *Handler) EnterBoss(e model.Event, out io.Writer) {
	player := h.repo.ActiveInDungeon(e.PlayerID)
	if player == nil {
		return
	}

	dungeon := h.repo.Dungeon()
	if player.Floor != dungeon.BossFloor || !player.BossEnter.IsZero() || player.BossDefeated {
		emit(out, e.Time, e.PlayerID, fmt.Sprintf("makes imposible move [%d]", e.ID))
		return
	}

	player.BossEnter = e.Time
	emit(out, e.Time, e.PlayerID, "entered the boss's floor")
}

func (h *Handler) KillBoss(e model.Event, out io.Writer) {
	player := h.repo.ActiveInDungeon(e.PlayerID)
	if player == nil {
		return
	}

	dungeon := h.repo.Dungeon()
	if player.Floor != dungeon.BossFloor || player.BossEnter.IsZero() || player.BossDefeated {
		emit(out, e.Time, e.PlayerID, fmt.Sprintf("makes imposible move [%d]", e.ID))
		return
	}

	player.BossDefeated = true
	player.BossKillAt = e.Time
	emit(out, e.Time, e.PlayerID, "killed the boss")
}

func (h *Handler) LeaveDungeon(e model.Event, out io.Writer) {
	player := h.repo.ActiveInDungeon(e.PlayerID)
	if player == nil {
		return
	}

	player.InDungeon = false
	player.LeftAt = e.Time
	emit(out, e.Time, e.PlayerID, "left the dungeon")
}

func (h *Handler) CannotCont(e model.Event, out io.Writer) {
	player := h.repo.GetPlayer(e.PlayerID)
	if player == nil || player.Status != model.StatusActive || !player.LeftAt.IsZero() {
		return
	}

	player.Status = model.StatusDisqual
	player.InDungeon = false
	player.LeftAt = e.Time
	emit(out, e.Time, e.PlayerID, fmt.Sprintf("cannot continue due to [%s]", e.Extra))
}

func (h *Handler) Heal(e model.Event, out io.Writer) {
	player := h.repo.ActiveInDungeon(e.PlayerID)
	if player == nil {
		return
	}

	amount, err := strconv.Atoi(e.Extra)
	if err != nil || amount <= 0 {
		return
	}

	player.HP += amount
	if player.HP > 100 {
		player.HP = 100
	}
	emit(out, e.Time, e.PlayerID, fmt.Sprintf("has restored [%d] of health", amount))
}

func (h *Handler) Damage(e model.Event, out io.Writer) {
	player := h.repo.ActiveInDungeon(e.PlayerID)
	if player == nil {
		return
	}

	amount, err := strconv.Atoi(e.Extra)
	if err != nil || amount <= 0 {
		return
	}

	player.HP -= amount
	emit(out, e.Time, e.PlayerID, fmt.Sprintf("recieved [%d] of damage", amount))

	if player.HP <= 0 {
		player.HP = 0
		player.Status = model.StatusDead
		player.InDungeon = false
		player.LeftAt = e.Time
		emit(out, e.Time, e.PlayerID, "is dead")
	}
}

func (h *Handler) Players() []*model.Player {
	return h.repo.AllPlayers()
}

func (h *Handler) CloseDungeon(closeAt time.Time) {
	for _, player := range h.repo.AllPlayers() {
		if player.InDungeon {
			player.InDungeon = false
			player.LeftAt = closeAt
		}
	}
}

func emit(out io.Writer, t time.Time, playerID int, msg string) {
	_, _ = fmt.Fprintf(out, "[%s] Player [%d] %s\n", t.Format("15:04:05"), playerID, msg)
}

func initMonstersLeft(d *model.Dungeon) map[int]int {
	m := make(map[int]int, d.TotalFloors-1)
	for i := 1; i < d.TotalFloors; i++ {
		m[i] = d.MonstersPerFloor
	}
	return m
}
