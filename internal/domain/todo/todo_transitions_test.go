package todo

import (
	"testing"
	"time"
)

func TestTodo_LifecycleAndEvents(t *testing.T) {
	now := time.Date(2025, 12, 13, 10, 0, 0, 0, time.UTC)

	title, _ := NewTitle("Learn Go")
	pri, _ := NewPriority("high")
	tags := NewTags([]string{"study"})

	td, ev, err := NewTodo(NewTodoParams{
		ID:       TodoID("t1"),
		Title:    title,
		Priority: pri,
		Tags:     tags,
		DueDate:  nil,
		Now:      now,
	})
	if err != nil {
		t.Fatalf("NewTodo err: %v", err)
	}
	if td.Status != StatusActive {
		t.Fatalf("status=%s want active", td.Status)
	}
	if len(ev) != 1 {
		t.Fatalf("events=%d want 1", len(ev))
	}

	// Complete
	td2, ev, err := td.Complete(now.Add(time.Minute))
	if err != nil {
		t.Fatalf("Complete err: %v", err)
	}
	if td2.Status != StatusDone || td2.CompletedAt == nil {
		t.Fatalf("expected done + completedAt, got %v", td2)
	}
	if len(ev) != 1 {
		t.Fatalf("Complete events=%d want 1", len(ev))
	}

	// Complete again (idempotent)
	td3, ev, err := td2.Complete(now.Add(2 * time.Minute))
	if err != nil {
		t.Fatalf("Complete idempotent err: %v", err)
	}
	if td3.Status != StatusDone || len(ev) != 0 {
		t.Fatalf("expected no-op; status=%s events=%d", td3.Status, len(ev))
	}

	// Archive from done
	td4, ev, err := td3.Archive(now.Add(3 * time.Minute))
	if err != nil {
		t.Fatalf("Archive err: %v", err)
	}
	if td4.Status != StatusArchived || td4.ArchivedAt == nil {
		t.Fatalf("expected archived + archivedAt")
	}
	if len(ev) != 1 {
		t.Fatalf("Archive events=%d want 1", len(ev))
	}

	// Restore to active
	td5, ev, err := td4.Restore(now.Add(4 * time.Minute))
	if err != nil {
		t.Fatalf("Restore err: %v", err)
	}
	if td5.Status != StatusActive {
		t.Fatalf("expected active got %s", td5.Status)
	}
	if len(ev) != 1 {
		t.Fatalf("Restore events=%d want 1", len(ev))
	}
}

func TestTodo_InvalidTransition_ActiveToArchive(t *testing.T) {
	now := time.Date(2025, 12, 13, 10, 0, 0, 0, time.UTC)

	title, _ := NewTitle("X")
	pri, _ := NewPriority("low")

	td, _, _ := NewTodo(NewTodoParams{
		ID:       TodoID("t1"),
		Title:    title,
		Priority: pri,
		Tags:     NewTags(nil),
		Now:      now,
	})

	_, _, err := td.Archive(now.Add(time.Minute))
	if err == nil {
		t.Fatalf("expected error")
	}
	if err != ErrInvalidTransition {
		t.Fatalf("err=%v want=%v", err, ErrInvalidTransition)
	}
}

func TestTodo_SoftDeleteBlocksEdits(t *testing.T) {
	now := time.Date(2025, 12, 13, 10, 0, 0, 0, time.UTC)
	title, _ := NewTitle("Keep it clean")
	pri, _ := NewPriority("medium")

	td, _, _ := NewTodo(NewTodoParams{
		ID:       TodoID("t1"),
		Title:    title,
		Priority: pri,
		Tags:     NewTags(nil),
		Now:      now,
	})

	td2, ev, err := td.SoftDelete(now.Add(time.Minute))
	if err != nil || len(ev) != 1 {
		t.Fatalf("SoftDelete err=%v events=%d", err, len(ev))
	}

	newTitle, _ := NewTitle("Nope")
	_, _, err = td2.ChangeTitle(newTitle, now.Add(2*time.Minute))
	if err != ErrDeletedTodo {
		t.Fatalf("err=%v want=%v", err, ErrDeletedTodo)
	}
}
