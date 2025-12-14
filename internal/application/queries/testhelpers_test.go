package queries

import (
	"context"
	"errors"
	"sort"
	"strings"
	"sync"
	"time"

	appErr "github.com/rojanmagar2001/gotodo/internal/application/errors"
	"github.com/rojanmagar2001/gotodo/internal/application/ports"
	"github.com/rojanmagar2001/gotodo/internal/domain/todo"
)

type inMemoryRepo struct {
	mu   sync.RWMutex
	data map[todo.TodoID]todo.Todo
}

func newInMemoryRepo(seed ...todo.Todo) *inMemoryRepo {
	r := &inMemoryRepo{data: map[todo.TodoID]todo.Todo{}}
	for _, t := range seed {
		r.data[t.ID] = t
	}
	return r
}

func (r *inMemoryRepo) Create(ctx context.Context, t todo.Todo) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.data[t.ID]; ok {
		return appErr.ErrConflict
	}
	r.data[t.ID] = t
	return nil
}

func (r *inMemoryRepo) Update(ctx context.Context, t todo.Todo) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.data[t.ID]; !ok {
		return appErr.ErrNotFound
	}
	r.data[t.ID] = t
	return nil
}

func (r *inMemoryRepo) GetByID(ctx context.Context, id todo.TodoID) (todo.Todo, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	t, ok := r.data[id]
	if !ok {
		return todo.Todo{}, appErr.ErrNotFound
	}
	return t, nil
}

func (r *inMemoryRepo) SoftDelete(ctx context.Context, id todo.TodoID) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	t, ok := r.data[id]
	if !ok {
		return appErr.ErrNotFound
	}
	now := time.Now().UTC()
	t.DeletedAt = &now
	t.UpdatedAt = now
	r.data[id] = t
	return nil
}

func (r *inMemoryRepo) HardDelete(ctx context.Context, id todo.TodoID) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.data, id)
	return nil
}

func (r *inMemoryRepo) List(ctx context.Context, spec ports.ListSpec) ([]todo.Todo, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// defaults
	sortBy := spec.SortBy
	if sortBy == "" {
		sortBy = ports.SortByCreated
	}
	order := spec.SortOrder
	if order == "" {
		order = ports.OrderDesc
	}

	var items []todo.Todo
	for _, t := range r.data {
		// includeDeleted filter
		if !spec.IncludeDeleted && t.DeletedAt != nil {
			continue
		}
		// status filter
		if spec.Status != nil && t.Status != *spec.Status {
			continue
		}
		// tag filter
		if spec.Tag != nil {
			if !t.Tags.Contains(*spec.Tag) {
				continue
			}
		}
		// search
		if spec.Search != nil {
			q := strings.ToLower(strings.TrimSpace(*spec.Search))
			if q != "" {
				if !strings.Contains(strings.ToLower(t.Title.String()), q) {
					continue
				}
			}
		}
		items = append(items, t)
	}

	// sorting helpers
	priorityRank := func(p todo.Priority) int {
		switch p {
		case todo.PriorityHigh:
			return 3
		case todo.PriorityMedium:
			return 2
		case todo.PriorityLow:
			return 1
		default:
			return 0
		}
	}

	dueKey := func(t todo.Todo) (time.Time, bool) {
		if t.DueDate == nil {
			return time.Time{}, false
		}
		return t.DueDate.AsTimeUTC(), true
	}

	sort.Slice(items, func(i, j int) bool {
		a, b := items[i], items[j]

		less := false
		switch sortBy {
		case ports.SortByCreated:
			less = a.CreatedAt.Before(b.CreatedAt)
		case ports.SortByUpdated:
			less = a.UpdatedAt.Before(b.UpdatedAt)
		case ports.SortByTitle:
			less = strings.ToLower(a.Title.String()) < strings.ToLower(b.Title.String())
		case ports.SortByPriority:
			less = priorityRank(a.Priority) < priorityRank(b.Priority)
		case ports.SortByDueDate:
			ad, aok := dueKey(a)
			bd, bok := dueKey(b)

			// nil due dates go LAST in asc, FIRST in desc (pick a rule)
			if !aok && !bok {
				less = false
			} else if !aok && bok {
				less = false // a after b in asc
			} else if aok && !bok {
				less = true
			} else {
				less = ad.Before(bd)
			}
		default:
			less = a.CreatedAt.Before(b.CreatedAt)
		}

		if order == ports.OrderAsc {
			return less
		}
		return !less
	})

	// paging
	if spec.Offset < 0 {
		spec.Offset = 0
	}
	if spec.Limit < 0 {
		spec.Limit = 0
	}
	if spec.Offset >= len(items) {
		return []todo.Todo{}, nil
	}
	if spec.Offset > 0 {
		items = items[spec.Offset:]
	}
	if spec.Limit > 0 && spec.Limit < len(items) {
		items = items[:spec.Limit]
	}

	return items, nil
}

type fakeClock struct{ t time.Time }

func (f fakeClock) Now() time.Time { return f.t }

// helper to build a Todo quickly (keeps tests readable)
func mkTodo(t testingT, id string, title string, status todo.Status, pr todo.Priority, tags []string, due *string, created time.Time) todo.Todo {
	t.Helper()

	tt, err := todo.NewTitle(title)
	if err != nil {
		t.Fatalf("bad title: %v", err)
	}
	// priority is a value object but you already have typed constants; ensure consistent:
	p := pr

	var dd *todo.DueDate
	if due != nil {
		d, err := todo.ParseDueDate(*due)
		if err != nil {
			t.Fatalf("bad due: %v", err)
		}
		dd = &d
	}

	td := todo.Todo{
		ID:        todo.TodoID(id),
		Title:     tt,
		Status:    status,
		Priority:  p,
		Tags:      todo.NewTags(tags),
		DueDate:   dd,
		CreatedAt: created,
		UpdatedAt: created,
	}
	return td
}

// tiny interface so we can call Helper/Fatalf without importing testing in every helper
type testingT interface {
	Helper()
	Fatalf(format string, args ...any)
}

// Ensure our repo satisfies interface at compile time
var _ ports.TodoRepository = (*inMemoryRepo)(nil)

// quick sanity guard (optional)
var _ = errors.Is
