package ports

import (
	"context"

	"github.com/rojanmagar2001/gotodo/internal/domain/todo"
)

type TodoRepository interface {
	Create(ctx context.Context, t todo.Todo) error
	Update(ctx context.Context, t todo.Todo) error
	GetByID(ctx context.Context, id todo.TodoID) (todo.Todo, error)

	List(ctx context.Context, spec ListSpec) ([]todo.Todo, error)

	SoftDelete(ctx context.Context, id todo.TodoID) error
	HardDelete(ctx context.Context, id todo.TodoID) error
}
