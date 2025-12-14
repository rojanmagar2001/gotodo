package jsonstore

import (
	"context"
	"strings"
	"time"

	appErr "github.com/rojanmagar2001/gotodo/internal/application/errors"
	"github.com/rojanmagar2001/gotodo/internal/application/ports"
	"github.com/rojanmagar2001/gotodo/internal/domain/todo"
)

type Repository struct {
	store Store
}

func NewRepository(path string) *Repository {
	return &Repository{store: New(path)}
}

func (r *Repository) Create(ctx context.Context, t todo.Todo) error {
	return r.withLock(func(fs *fileSchema) error {
		for _, row := range fs.Todos {
			if row.ID == t.ID.String() {
				return appErr.ErrConflict
			}
		}
		fs.Todos = append(fs.Todos, toRow(t))
		return nil
	})
}

func (r *Repository) Update(ctx context.Context, t todo.Todo) error {
	return r.withLock(func(fs *fileSchema) error {
		for i := range fs.Todos {
			if fs.Todos[i].ID == t.ID.String() {
				fs.Todos[i] = toRow(t)
				return nil
			}
		}
		return appErr.ErrNotFound
	})
}

func (r *Repository) GetByID(ctx context.Context, id todo.TodoID) (todo.Todo, error) {
	fs, err := r.store.Load()
	if err != nil {
		return todo.Todo{}, err
	}

	for _, row := range fs.Todos {
		if row.ID == id.String() {
			return fromRow(row)
		}
	}
	return todo.Todo{}, appErr.ErrNotFound
}

func (r *Repository) List(ctx context.Context, spec ports.ListSpec) ([]todo.Todo, error) {
	fs, err := r.store.Load()
	if err != nil {
		return nil, err
	}

	// convert + filter
	var out []todo.Todo
	for _, row := range fs.Todos {
		td, err := fromRow(row)
		if err != nil {
			return nil, err
		}

		if !spec.IncludeDeleted && td.DeletedAt != nil {
			continue
		}
		if spec.Status != nil && td.Status != *spec.Status {
			continue
		}
		if spec.Tag != nil && !td.Tags.Contains(*spec.Tag) {
			continue
		}
		if spec.Search != nil {
			q := strings.ToLower(strings.TrimSpace(*spec.Search))
			if q != "" && !strings.Contains(strings.ToLower(td.Title.String()), q) {
				continue
			}
		}

		out = append(out, td)
	}

	return out, nil
}

func (r *Repository) SoftDelete(ctx context.Context, id todo.TodoID) error {
	// prefer application to call domain SoftDelete + Update,
	// but keep this for port completeness:
	td, err := r.GetByID(ctx, id)
	if err != nil {
		return err
	}
	updated, _, err := td.SoftDelete(timeNowUTC())
	if err != nil {
		return err
	}
	return r.Update(ctx, updated)
}

func (r *Repository) HardDelete(ctx context.Context, id todo.TodoID) error {
	return r.withLock(func(fs *fileSchema) error {
		n := fs.Todos[:0]
		found := false
		for _, row := range fs.Todos {
			if row.ID == id.String() {
				found = true
				continue
			}
			n = append(n, row)
		}
		fs.Todos = n
		if !found {
			return appErr.ErrNotFound
		}
		return nil
	})
}

// internal helper: load -> mutate -> save with lock
func (r *Repository) withLock(mut func(fs *fileSchema) error) error {
	l, err := acquireLock(r.store.Path)
	if err != nil {
		return err
	}
	defer func() { _ = l.release() }()

	fs, err := r.store.Load()
	if err != nil {
		return err
	}
	if err := mut(&fs); err != nil {
		return err
	}
	if err := r.store.Save(fs); err != nil {
		return err
	}

	return nil
}

// Best-effort UTC clock for infra convenience
func timeNowUTC() time.Time { return time.Now().UTC() }

// fromRow can return error when stored data violates domain constraints
func fromRow(row todoRow) (todo.Todo, error) {
	title, err := todo.NewTitle(row.Title)
	if err != nil {
		return todo.Todo{}, ErrCorruptData
	}

	priority, err := todo.NewPriority(row.Priority)
	if err != nil {
		return todo.Todo{}, ErrCorruptData
	}

	st := todo.Status(row.Status)
	if !st.Valid() {
		return todo.Todo{}, ErrCorruptData
	}

	var dd *todo.DueDate
	if row.DueDate != nil {
		d, err := todo.ParseDueDate(*row.DueDate)
		if err != nil {
			return todo.Todo{}, ErrCorruptData
		}
		dd = &d
	}

	td := todo.Todo{
		ID:          todo.TodoID(row.ID),
		Title:       title,
		Status:      st,
		Priority:    priority,
		Tags:        todo.NewTags(row.Tags),
		DueDate:     dd,
		CreatedAt:   row.CreatedAt,
		UpdatedAt:   row.UpdatedAt,
		CompletedAt: row.CompletedAt,
		ArchivedAt:  row.ArchivedAt,
		DeletedAt:   row.DeletedAt,
	}

	return td, nil
}

func toRow(t todo.Todo) todoRow {
	var due *string
	if t.DueDate != nil {
		s := t.DueDate.String()
		due = &s
	}

	// copy tags so we never serialize shared slice
	tags := make([]string, len(t.Tags))
	copy(tags, t.Tags)

	return todoRow{
		ID:       t.ID.String(),
		Title:    t.Title.String(),
		Status:   string(t.Status),
		Priority: t.Priority.String(),
		Tags:     tags,
		DueDate:  due,

		CreatedAt:   t.CreatedAt,
		UpdatedAt:   t.UpdatedAt,
		CompletedAt: t.CompletedAt,
		ArchivedAt:  t.ArchivedAt,
		DeletedAt:   t.DeletedAt,
	}
}
