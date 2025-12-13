package todo

import "strings"

type Priority string

const (
	PriorityLow    Priority = "low"
	PriorityMedium Priority = "medium"
	PriorityHigh   Priority = "high"
)

func NewPriority(raw string) (Priority, error) {
	v := strings.TrimSpace(strings.ToLower(raw))
	switch Priority(v) {
	case PriorityLow, PriorityMedium, PriorityHigh:
		return Priority(v), nil
	default:
		return "", ErrInvalidPriority
	}
}

func (p Priority) String() string { return string(p) }
