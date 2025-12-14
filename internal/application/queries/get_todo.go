package queries

import (
	"context"

	appErr "github.com/rojanmagar2001/gotodo/internal/application/errors"
	"github.com/rojanmagar2001/gotodo/internal/application/ports"
	"github.com/rojanmagar2001/gotodo/internal/application/result"
	"github.com/rojanmagar2001/gotodo/internal/domain/todo"
)

type GetTodo struct {
	Repo ports.TodoRepository
}

func (q GetTodo) Execute(ctx context.Context, id todo.TodoID) result.Result[TodoDTO] {
	td, err := q.Repo.GetByID(ctx, id)
	if err != nil {
		return result.Fail[TodoDTO](appErr.ErrNotFound)
	}
	return result.Ok(ToDTO(td))
}
