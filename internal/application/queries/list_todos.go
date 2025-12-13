package queries

import (
	"context"

	"github.com/rojanmagar2001/gotodo/internal/application/ports"
	"github.com/rojanmagar2001/gotodo/internal/application/result"
	"github.com/rojanmagar2001/gotodo/internal/domain/todo"
)

type ListTodos struct {
	Repo ports.TodoRepository
}

func (q ListTodos) Execute(ctx context.Context, spec ports.ListSpec) result.Result[[]todo.Todo] {
	tds, err := q.Repo.List(ctx, spec)
	if err != nil {
		return result.Fail[[]todo.Todo](err)
	}

	return result.Ok(tds)
}
