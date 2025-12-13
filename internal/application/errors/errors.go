package errors

import "errors"

var (
	ErrNotFound   = errors.New("not found")
	ErrConflict   = errors.New("conflict")
	ErrValidation = errors.New("validation failed")
	ErrUnExpected = errors.New("unxpected error")
)
