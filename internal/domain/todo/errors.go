package todo

import "errors"

var (
	ErrInvalidTitle      = errors.New("invalid title")
	ErrInvalidPriority   = errors.New("invalid priority")
	ErrInvalidDueDate    = errors.New("invalid due date")
	ErrInvalidTransition = errors.New("invalid status transition")
	ErrDeletedTodo       = errors.New("todo is deleted")
)
