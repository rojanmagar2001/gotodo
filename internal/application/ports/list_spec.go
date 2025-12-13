package ports

import "github.com/rojanmagar2001/gotodo/internal/domain/todo"

type ListSpec struct {
	Status *todo.Status
	Tag    *string
	Search *string
	SortBy SortField
}

type SortField string

const (
	SortByCreated  SortField = "created"
	SortByDueDate  SortField = "due"
	SortByPriority SortField = "priority"
)
