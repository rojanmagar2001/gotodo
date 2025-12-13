package commands

import (
	"context"

	appErr "github.com/rojanmagar2001/gotodo/internal/application/errors"
	"github.com/rojanmagar2001/gotodo/internal/application/ports"
	"github.com/rojanmagar2001/gotodo/internal/application/result"
	"github.com/rojanmagar2001/gotodo/internal/domain/todo"
)

type AddTodo struct {
	Repo      ports.TodoRepository
	Clock     ports.Clock
	IDGen     ports.IDGenerator
	Publisher ports.EventPublisher
}

type AddTodoInput struct {
	Title    string
	Priority string
	Tags     []string
	DueDate  *string // YYYY-MM-DD
}

func (uc AddTodo) Execute(ctx context.Context, in AddTodoInput) result.Result[todo.Todo] {
	title, err := todo.NewTitle(in.Title)
	if err != nil {
		return result.Fail[todo.Todo](appErr.ErrValidation)
	}

	priority, err := todo.NewPriority(in.Priority)
	if err != nil {
		return result.Fail[todo.Todo](appErr.ErrValidation)
	}

	var due *todo.DueDate
	if in.DueDate != nil {
		d, err := todo.ParseDueDate(*in.DueDate)
		if err != nil {
			return result.Fail[todo.Todo](appErr.ErrValidation)
		}
		due = &d
	}

	td, events, err := todo.NewTodo(todo.NewTodoParams{
		ID:       uc.IDGen.NewTodoID(),
		Title:    title,
		Priority: priority,
		Tags:     todo.NewTags(in.Tags),
		DueDate:  due,
		Now:      uc.Clock.Now(),
	})
	if err != nil {
		return result.Fail[todo.Todo](appErr.ErrUnExpected)
	}

	if err := uc.Repo.Create(ctx, td); err != nil {
		return result.Fail[todo.Todo](appErr.ErrUnExpected)
	}

	_ = uc.Publisher.Publish(ctx, events)

	return result.Ok(td)
}
