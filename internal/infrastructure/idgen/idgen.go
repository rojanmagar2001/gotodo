package idgen

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/rojanmagar2001/gotodo/internal/domain/todo"
)

type RandomIDGen struct{}

func (RandomIDGen) NewTodoID() todo.TodoID {
	var b [8]byte // 16 hex characters
	_, _ = rand.Read(b[:])
	return todo.TodoID(hex.EncodeToString(b[:]))
}
