package queries

import (
	"context"
	"testing"
	"time"

	"github.com/rojanmagar2001/gotodo/internal/application/ports"
	"github.com/rojanmagar2001/gotodo/internal/domain/todo"
)

func TestListTodos_FilterByStatus(t *testing.T) {
	ctx := context.Background()
	base := time.Date(2025, 12, 10, 10, 0, 0, 0, time.UTC)

	td1 := mkTodo(t, "1", "Buy milk", todo.StatusActive, todo.PriorityLow, []string{"home"}, nil, base)
	td2 := mkTodo(t, "2", "Write report", todo.StatusDone, todo.PriorityHigh, []string{"work"}, nil, base.Add(time.Minute))
	td3 := mkTodo(t, "3", "Pay bills", todo.StatusActive, todo.PriorityMedium, []string{"home"}, nil, base.Add(2*time.Minute))

	repo := newInMemoryRepo(td1, td2, td3)
	q := ListTodos{Repo: repo}

	s := todo.StatusActive
	spec := ports.ListSpec{Status: &s, SortBy: ports.SortByCreated, SortOrder: ports.OrderAsc}

	res := q.Execute(ctx, spec)
	if res.Err != nil {
		t.Fatalf("err=%v", res.Err)
	}
	if len(res.Value) != 2 {
		t.Fatalf("len=%d want=2", len(res.Value))
	}
	if res.Value[0].ID != "1" || res.Value[1].ID != "3" {
		t.Fatalf("ids=%v want=[1 3]", []string{res.Value[0].ID, res.Value[1].ID})
	}
}

func TestListTodos_FilterByTag(t *testing.T) {
	ctx := context.Background()
	base := time.Date(2025, 12, 10, 10, 0, 0, 0, time.UTC)

	td1 := mkTodo(t, "1", "Buy milk", todo.StatusActive, todo.PriorityLow, []string{"home"}, nil, base)
	td2 := mkTodo(t, "2", "Write report", todo.StatusActive, todo.PriorityHigh, []string{"work"}, nil, base.Add(time.Minute))
	td3 := mkTodo(t, "3", "Email boss", todo.StatusActive, todo.PriorityMedium, []string{"work"}, nil, base.Add(2*time.Minute))

	repo := newInMemoryRepo(td1, td2, td3)
	q := ListTodos{Repo: repo}

	tag := "work"
	spec := ports.ListSpec{Tag: &tag, SortBy: ports.SortByCreated, SortOrder: ports.OrderAsc}

	res := q.Execute(ctx, spec)
	if res.Err != nil {
		t.Fatalf("err=%v", res.Err)
	}
	if len(res.Value) != 2 {
		t.Fatalf("len=%d want=2", len(res.Value))
	}
	if res.Value[0].ID != "2" || res.Value[1].ID != "3" {
		t.Fatalf("ids=%v want=[2 3]", []string{res.Value[0].ID, res.Value[1].ID})
	}
}

func TestListTodos_SearchCaseInsensitive(t *testing.T) {
	ctx := context.Background()
	base := time.Date(2025, 12, 10, 10, 0, 0, 0, time.UTC)

	td1 := mkTodo(t, "1", "Buy milk", todo.StatusActive, todo.PriorityLow, []string{"home"}, nil, base)
	td2 := mkTodo(t, "2", "BUY coffee", todo.StatusActive, todo.PriorityHigh, []string{"home"}, nil, base.Add(time.Minute))
	td3 := mkTodo(t, "3", "Write report", todo.StatusActive, todo.PriorityMedium, []string{"work"}, nil, base.Add(2*time.Minute))

	repo := newInMemoryRepo(td1, td2, td3)
	q := ListTodos{Repo: repo}

	search := "buy"
	spec := ports.ListSpec{Search: &search, SortBy: ports.SortByCreated, SortOrder: ports.OrderAsc}

	res := q.Execute(ctx, spec)
	if res.Err != nil {
		t.Fatalf("err=%v", res.Err)
	}
	if len(res.Value) != 2 {
		t.Fatalf("len=%d want=2", len(res.Value))
	}
	if res.Value[0].ID != "1" || res.Value[1].ID != "2" {
		t.Fatalf("ids=%v want=[1 2]", []string{res.Value[0].ID, res.Value[1].ID})
	}
}

func TestListTodos_ExcludeDeletedByDefault(t *testing.T) {
	ctx := context.Background()
	base := time.Date(2025, 12, 10, 10, 0, 0, 0, time.UTC)

	td1 := mkTodo(t, "1", "Keep me", todo.StatusActive, todo.PriorityLow, []string{"x"}, nil, base)
	td2 := mkTodo(t, "2", "Delete me", todo.StatusActive, todo.PriorityLow, []string{"x"}, nil, base.Add(time.Minute))
	delAt := base.Add(2 * time.Minute)
	td2.DeletedAt = &delAt

	repo := newInMemoryRepo(td1, td2)
	q := ListTodos{Repo: repo}

	spec := ports.ListSpec{SortBy: ports.SortByCreated, SortOrder: ports.OrderAsc}
	res := q.Execute(ctx, spec)
	if res.Err != nil {
		t.Fatalf("err=%v", res.Err)
	}
	if len(res.Value) != 1 {
		t.Fatalf("len=%d want=1", len(res.Value))
	}
	if res.Value[0].ID != "1" {
		t.Fatalf("id=%s want=1", res.Value[0].ID)
	}

	// include deleted
	spec.IncludeDeleted = true
	res2 := q.Execute(ctx, spec)
	if res2.Err != nil {
		t.Fatalf("err=%v", res2.Err)
	}
	if len(res2.Value) != 2 {
		t.Fatalf("len=%d want=2", len(res2.Value))
	}
}
