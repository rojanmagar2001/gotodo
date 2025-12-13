package errors

import (
	"errors"

	domain "github.com/rojanmagar2001/gotodo/internal/domain/todo"
)

func MapDomainError(err error) error {
	if err == nil {
		return nil
	}

	switch {
	case errors.Is(err, domain.ErrInvalidTitle),
		errors.Is(err, domain.ErrInvalidPriority),
		errors.Is(err, domain.ErrInvalidDueDate),
		errors.Is(err, domain.ErrInvalidTransition),
		errors.Is(err, domain.ErrDeletedTodo):
		return ErrValidation
	default:
		return ErrUnExpected
	}
}
