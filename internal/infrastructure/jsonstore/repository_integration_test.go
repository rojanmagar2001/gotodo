package jsonstore

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/rojanmagar2001/gotodo/internal/application/ports"
	"github.com/rojanmagar2001/gotodo/internal/domain/todo"
)

func TestRepository_CreateGetUpdate_List(t *testing.T) {
	ctx := context.Background()
	dir := t.TempDir()
	path := filepath.Join(dir, "todos.json")

	repo := NewRepository(path)

	base := time.Date(2025, 12, 14, 10, 0, 0, 0, time.UTC)
	title, _ := todo.NewTitle("Buy milk")
	pri, _ := todo.NewPriority("low")
	td, _, err := todo.NewTodo(todo.NewTodoParams{
		ID:       todo.TodoID("t1"),
		Title:    title,
		Priority: pri,
		Tags:     todo.NewTags([]string{"home"}),
		Now:      base,
	})
	if err != nil {
		t.Fatalf("NewTodo err=%v", err)
	}

	if err := repo.Create(ctx, td); err != nil {
		t.Fatalf("Create err=%v", err)
	}

	got, err := repo.GetByID(ctx, todo.TodoID("t1"))
	if err != nil {
		t.Fatalf("GetByID err=%v", err)
	}
	if got.Title.String() != "Buy milk" {
		t.Fatalf("title=%q", got.Title.String())
	}

	// update title
	newTitle, _ := todo.NewTitle("Buy oat milk")
	updated, _, err := got.ChangeTitle(newTitle, base.Add(time.Minute))
	if err != nil {
		t.Fatalf("ChangeTitle err=%v", err)
	}
	if err := repo.Update(ctx, updated); err != nil {
		t.Fatalf("Update err=%v", err)
	}

	list, err := repo.List(ctx, ports.ListSpec{})
	if err != nil {
		t.Fatalf("List err=%v", err)
	}
	if len(list) != 1 {
		t.Fatalf("len=%d want=1", len(list))
	}
	if list[0].Title.String() != "Buy oat milk" {
		t.Fatalf("title=%q", list[0].Title.String())
	}
}

func TestRepository_AtomicWrite_FileExists(t *testing.T) {
	ctx := context.Background()
	dir := t.TempDir()
	path := filepath.Join(dir, "todos.json")
	repo := NewRepository(path)

	base := time.Date(2025, 12, 14, 10, 0, 0, 0, time.UTC)
	title, _ := todo.NewTitle("A")
	pri, _ := todo.NewPriority("low")
	td, _, _ := todo.NewTodo(todo.NewTodoParams{
		ID:       todo.TodoID("t1"),
		Title:    title,
		Priority: pri,
		Tags:     todo.NewTags(nil),
		Now:      base,
	})

	if err := repo.Create(ctx, td); err != nil {
		t.Fatalf("Create err=%v", err)
	}

	// Ensure main file exists (atomic rename succeeded)
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("expected file to exist: %v", err)
	}
}
