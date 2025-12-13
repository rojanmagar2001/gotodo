package ports

import (
	"context"

	"github.com/rojanmagar2001/gotodo/internal/domain/todo"
)

type EventPublisher interface {
	Publish(ctx context.Context, events []todo.Event) error
}
