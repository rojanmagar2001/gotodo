package jsonstore

import "errors"

var (
	ErrCorruptData = errors.New("jsonstore: corrupt data")
	ErrLocked      = errors.New("jsonstore: store is locked")
)
