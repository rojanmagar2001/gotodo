package ports

import "github.com/rojanmagar2001/gotodo/internal/domain/todo"

type ListSpec struct {
	// filters
	Status *todo.Status
	Tag    *string

	// search
	Search *string // full-text-isj: title contains (case-insensitive)

	// sort
	SortBy    SortField
	SortOrder SortOrder

	// paging (optional now, used later by ui)
	Limit  int
	Offset int

	// include soft-deleted?
	IncludeDeleted bool
}

type (
	SortField string
	SortOrder string
)

const (
	SortByCreated  SortField = "created"
	SortByDueDate  SortField = "due"
	SortByPriority SortField = "priority"
	SortByTitle    SortField = "title"
	SortByUpdated  SortField = "updated"

	OrderAsc  SortOrder = "asc"
	OrderDesc SortOrder = "desc"
)
