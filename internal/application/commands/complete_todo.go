package commands

import (
	"context"

	appErr "github.com/rojanmagar2001/gotodo/internal/application/errors"
	"github.com/rojanmagar2001/gotodo/internal/application/ports"
	"github.com/rojanmagar2001/gotodo/internal/application/result"
	"github.com/rojanmagar2001/gotodo/internal/domain/todo"
)

type CompleteTodo struct {
	Repo      ports.TodoRepository
	Clock     ports.Clock
	Publisher ports.EventPublisher
}

func (uc CompleteTodo) Execute(ctx context.Context, id todo.TodoID) result.Result[todo.Todo] {
	td, err := uc.Repo.GetByID(ctx, id)
	if err != nil {
		return result.Fail[todo.Todo](appErr.ErrNotFound)
	}

	updated, events, err := td.Complete(uc.Clock.Now())
	if err != nil {
		return result.Fail[todo.Todo](appErr.ErrValidation)
	}

	if err := uc.Repo.Update(ctx, updated); err != nil {
		return result.Fail[todo.Todo](appErr.ErrUnExpected)
	}

	_ = uc.Publisher.Publish(ctx, events)

	return result.Ok(updated)
}
