package queries

import (
	"context"

	appErr "github.com/rojanmagar2001/gotodo/internal/application/errors"
	"github.com/rojanmagar2001/gotodo/internal/application/ports"
	"github.com/rojanmagar2001/gotodo/internal/application/result"
)

type ListTodos struct {
	Repo ports.TodoRepository
}

func (q ListTodos) Execute(ctx context.Context, spec ports.ListSpec) result.Result[[]TodoDTO] {
	tds, err := q.Repo.List(ctx, spec)
	if err != nil {
		return result.Fail[[]TodoDTO](appErr.ErrUnExpected)
	}

	out := make([]TodoDTO, 0, len(tds))
	for _, t := range tds {
		out = append(out, ToDTO(t))
	}

	return result.Ok(out)
}
