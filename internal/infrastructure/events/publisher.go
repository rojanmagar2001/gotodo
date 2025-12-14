package events

import (
	"context"
	"log"

	"github.com/rojanmagar2001/gotodo/internal/domain/todo"
)

type LogPublisher struct {
	L *log.Logger
}

func (p LogPublisher) Publish(ctx context.Context, evs []todo.Event) error {
	if len(evs) == 0 {
		return nil
	}

	for _, e := range evs {
		// keep it minimal; richer later
		p.L.Printf("event: %T", e)
	}
	return nil
}
