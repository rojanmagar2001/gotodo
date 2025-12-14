package queries

import (
	"context"
	"testing"
	"time"

	appErr "github.com/rojanmagar2001/gotodo/internal/application/errors"
	"github.com/rojanmagar2001/gotodo/internal/domain/todo"
)

func TestGetTodo_Found(t *testing.T) {
	ctx := context.Background()
	base := time.Date(2025, 12, 10, 10, 0, 0, 0, time.UTC)

	td := mkTodo(t, "1", "Buy milk", todo.StatusActive, todo.PriorityLow, []string{"home"}, nil, base)

	repo := newInMemoryRepo(td)
	q := GetTodo{Repo: repo}

	res := q.Execute(ctx, todo.TodoID("1"))
	if res.Err != nil {
		t.Fatalf("err=%v", res.Err)
	}
	if res.Value.ID != "1" || res.Value.Title != "Buy milk" {
		t.Fatalf("dto=%+v", res.Value)
	}
}

func TestGetTodo_NotFound(t *testing.T) {
	ctx := context.Background()
	repo := newInMemoryRepo()
	q := GetTodo{Repo: repo}

	res := q.Execute(ctx, todo.TodoID("missing"))
	if res.Err != appErr.ErrNotFound {
		t.Fatalf("err=%v want=%v", res.Err, appErr.ErrNotFound)
	}
}
