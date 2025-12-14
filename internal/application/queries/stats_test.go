package queries

import (
	"context"
	"testing"
	"time"

	"github.com/rojanmagar2001/gotodo/internal/domain/todo"
)

func TestStats_CountsAndDueBuckets(t *testing.T) {
	ctx := context.Background()

	// Clock "now" is Dec 14 2025 10:00 UTC
	now := time.Date(2025, 12, 14, 10, 0, 0, 0, time.UTC)
	clock := fakeClock{t: now}

	base := time.Date(2025, 12, 10, 10, 0, 0, 0, time.UTC)

	overdue := "2025-12-13"
	today := "2025-12-14"
	soon := "2025-12-18"
	later := "2026-01-01"

	activeOverdue := mkTodo(t, "1", "A overdue", todo.StatusActive, todo.PriorityLow, nil, &overdue, base)
	activeToday := mkTodo(t, "2", "B today", todo.StatusActive, todo.PriorityLow, nil, &today, base.Add(time.Minute))
	activeSoon := mkTodo(t, "3", "C soon", todo.StatusActive, todo.PriorityLow, nil, &soon, base.Add(2*time.Minute))
	activeLater := mkTodo(t, "4", "D later", todo.StatusActive, todo.PriorityLow, nil, &later, base.Add(3*time.Minute))

	doneToday := mkTodo(t, "5", "E done", todo.StatusDone, todo.PriorityLow, nil, &today, base.Add(4*time.Minute))
	arch := mkTodo(t, "6", "F arch", todo.StatusArchived, todo.PriorityLow, nil, nil, base.Add(5*time.Minute))

	deleted := mkTodo(t, "7", "G deleted", todo.StatusActive, todo.PriorityLow, nil, nil, base.Add(6*time.Minute))
	delAt := base.Add(7 * time.Minute)
	deleted.DeletedAt = &delAt

	repo := newInMemoryRepo(activeOverdue, activeToday, activeSoon, activeLater, doneToday, arch, deleted)

	q := Stats{Repo: repo, Clock: clock}
	res := q.Execute(ctx)
	if res.Err != nil {
		t.Fatalf("err=%v", res.Err)
	}

	s := res.Value
	if s.Total != 7 {
		t.Fatalf("Total=%d want=7", s.Total)
	}
	if s.Active != 4 {
		t.Fatalf("Active=%d want=4", s.Active)
	}
	if s.Done != 1 {
		t.Fatalf("Done=%d want=1", s.Done)
	}
	if s.Archived != 1 {
		t.Fatalf("Archived=%d want=1", s.Archived)
	}
	if s.Deleted != 1 {
		t.Fatalf("Deleted=%d want=1", s.Deleted)
	}

	// Due buckets should only count ACTIVE todos with a due date
	if s.Overdue != 1 {
		t.Fatalf("Overdue=%d want=1", s.Overdue)
	}
	if s.DueToday != 1 {
		t.Fatalf("DueToday=%d want=1", s.DueToday)
	}
	if s.DueSoon != 1 {
		t.Fatalf("DueSoon=%d want=1", s.DueSoon)
	}
}
