package jsonstore

import "time"

const schemaVersion = 1

type fileSchema struct {
	Version int       `json:"version"`
	SavedAt time.Time `json:"savedAt"`
	Todos   []todoRow `json:"todos"`
}

type todoRow struct {
	ID       string   `json:"id"`
	Title    string   `json:"title"`
	Status   string   `json:"status"`
	Priority string   `json:"priority"`
	Tags     []string `json:"tags"`
	DueDate  *string  `json:"dueDate"`

	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
	CompletedAt *time.Time `json:"completedAt"`
	ArchivedAt  *time.Time `json:"archivedAt"`
	DeletedAt   *time.Time `json:"deletedAt"`
}
