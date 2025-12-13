package commands

import (
	"context"

	appErr "github.com/rojanmagar2001/gotodo/internal/application/errors"
	"github.com/rojanmagar2001/gotodo/internal/application/ports"
	"github.com/rojanmagar2001/gotodo/internal/application/result"
	"github.com/rojanmagar2001/gotodo/internal/domain/todo"
)

type SoftDeleteTodo struct {
	Repo      ports.TodoRepository
	Clock     ports.Clock
	Publisher ports.EventPublisher
	Undo      *UndoManager
}

func (uc SoftDeleteTodo) Execute(ctx context.Context, id todo.TodoID) result.Result[todo.Todo] {
	current, err := uc.Repo.GetByID(ctx, id)
	if err != nil {
		return result.Fail[todo.Todo](appErr.ErrNotFound)
	}
	before := current

	updated, events, err := current.SoftDelete(uc.Clock.Now())
	if err != nil {
		return result.Fail[todo.Todo](appErr.MapDomainError(err))
	}

	if err := uc.Repo.Update(ctx, updated); err != nil {
		return result.Fail[todo.Todo](appErr.ErrUnExpected)
	}
	_ = uc.Publisher.Publish(ctx, events)

	changed := len(events) > 0

	if uc.Undo != nil && changed {
		uc.Undo.Push(func(ctx context.Context) error {
			return uc.Repo.Update(ctx, before) // “undelete” by restoring snapshot
		})
	}
	return result.Ok(updated)
}

type HardDeleteTodo struct {
	Repo ports.TodoRepository
	Undo *UndoManager
}

func (uc HardDeleteTodo) Execute(ctx context.Context, id todo.TodoID) result.Result[struct{}] {
	// For undo, capture snapshot first (optional)
	var before *todo.Todo
	if uc.Undo != nil {
		td, err := uc.Repo.GetByID(ctx, id)
		if err == nil {
			before = &td
		}
	}

	if err := uc.Repo.HardDelete(ctx, id); err != nil {
		return result.Fail[struct{}](appErr.ErrUnExpected)
	}

	if uc.Undo != nil && before != nil {
		uc.Undo.Push(func(ctx context.Context) error {
			return uc.Repo.Create(ctx, *before)
		})
	}

	return result.Ok(struct{}{})
}
