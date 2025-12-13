package ports

import "github.com/rojanmagar2001/gotodo/internal/domain/todo"

type IDGenerator interface {
	NewTodoID() todo.TodoID
}
