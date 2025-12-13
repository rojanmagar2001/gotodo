package commands

import (
	"context"

	appErr "github.com/rojanmagar2001/gotodo/internal/application/errors"
	"github.com/rojanmagar2001/gotodo/internal/application/ports"
	"github.com/rojanmagar2001/gotodo/internal/application/result"
	"github.com/rojanmagar2001/gotodo/internal/domain/todo"
)

type EditTodo struct {
	Repo      ports.TodoRepository
	Clock     ports.Clock
	Publisher ports.EventPublisher
	Undo      *UndoManager
}

type EditTodoInput struct {
	ID       todo.TodoID
	Title    *string
	Priority *string
	Tags     *[]string
	DueDate  **string
}

func (uc EditTodo) Execute(ctx context.Context, in EditTodoInput) result.Result[todo.Todo] {
	current, err := uc.Repo.GetByID(ctx, in.ID)
	if err != nil {
		return result.Fail[todo.Todo](appErr.ErrNotFound)
	}
	before := current
	now := uc.Clock.Now()

	var events []todo.Event

	if in.Title != nil {
		tt, err := todo.NewTitle(*in.Title)
		if err != nil {
			return result.Fail[todo.Todo](appErr.ErrValidation)
		}
		updated, ev, err := current.ChangeTitle(tt, now)
		if err != nil {
			return result.Fail[todo.Todo](appErr.MapDomainError(err))
		}
		current = updated
		events = append(events, ev...)
	}

	if in.Priority != nil {
		pp, err := todo.NewPriority(*in.Priority)
		if err != nil {
			return result.Fail[todo.Todo](appErr.ErrValidation)
		}
		// add a domain method later if you want invariants/events; for new set directly
		current.Priority = pp
		current.UpdatedAt = now
		// optional: add TodoPriorityChanged event in domain if you want parity
	}

	if in.Tags != nil {
		current.Tags = todo.NewTags(*in.Tags)
		current.UpdatedAt = now
	}

	if in.DueDate != nil {
		// clear
		if *in.DueDate == nil {
			current.DueDate = nil
			current.UpdatedAt = now
		} else {
			d, err := todo.ParseDueDate(**in.DueDate)
			if err != nil {
				return result.Fail[todo.Todo](appErr.ErrValidation)
			}
			current.DueDate = &d
			current.UpdatedAt = now
		}
	}

	if err := uc.Repo.Update(ctx, current); err != nil {
		return result.Fail[todo.Todo](appErr.ErrUnExpected)
	}

	_ = uc.Publisher.Publish(ctx, events)

	changed := len(events) > 0

	// Undo: restore full snapshot
	if uc.Undo != nil && changed {
		uc.Undo.Push(func(ctx context.Context) error {
			return uc.Repo.Update(ctx, before)
		})
	}

	return result.Ok(current)
}
