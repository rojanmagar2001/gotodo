package queries

import (
	"time"

	"github.com/rojanmagar2001/gotodo/internal/domain/todo"
)

type TodoDTO struct {
	ID       string
	Title    string
	Status   string
	Priority string
	Tags     []string
	DueDate  *string

	CreatedAt   time.Time
	UpdatedAt   time.Time
	CompletedAt *time.Time
	ArchivedAt  *time.Time
	DeletedAt   *time.Time
}

func ToDTO(t todo.Todo) TodoDTO {
	var due *string
	if t.DueDate != nil {
		s := t.DueDate.String()
		due = &s
	}

	// copy tags to avoid sharing underlying slice
	tags := make([]string, len(t.Tags))
	copy(tags, t.Tags)

	return TodoDTO{
		ID:       t.ID.String(),
		Title:    t.Title.String(),
		Status:   string(t.Status),
		Priority: t.Priority.String(),
		Tags:     tags,
		DueDate:  due,

		CreatedAt:   t.CreatedAt,
		UpdatedAt:   t.UpdatedAt,
		CompletedAt: t.CompletedAt,
		ArchivedAt:  t.ArchivedAt,
		DeletedAt:   t.DeletedAt,
	}
}
