package commands

import (
	"context"
	"sync"

	appErr "github.com/rojanmagar2001/gotodo/internal/application/errors"
)

type UndoAction func(ctx context.Context) error

type UndoManager struct {
	mu    sync.Mutex
	stack []UndoAction
}

func NewUndoManager() *UndoManager { return &UndoManager{} }

func (u *UndoManager) Push(a UndoAction) {
	u.mu.Lock()
	defer u.mu.Unlock()
	u.stack = append(u.stack, a)
}

func (u *UndoManager) Undo(ctx context.Context) error {
	u.mu.Lock()
	if len(u.stack) == 0 {
		u.mu.Unlock()
		return appErr.ErrUnExpected
	}
	last := u.stack[len(u.stack)-1]
	u.stack = u.stack[:len(u.stack)-1]
	u.mu.Unlock()

	if err := last(ctx); err != nil {
		return appErr.ErrUnExpected
	}

	return nil
}
