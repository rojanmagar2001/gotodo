package todo

import "time"

type Todo struct {
	ID       TodoID
	Title    Title
	Status   Status
	Priority Priority
	Tags     Tags
	DueDate  *DueDate

	CreatedAt   time.Time
	UpdatedAt   time.Time
	CompletedAt *time.Time
	ArchivedAt  *time.Time
	DeletedAt   *time.Time
}

type NewTodoParams struct {
	ID       TodoID
	Title    Title
	Priority Priority
	Tags     Tags
	DueDate  *DueDate
	Now      time.Time
}

func NewTodo(p NewTodoParams) (Todo, []Event, error) {
	if !p.ID.Valid() {
		return Todo{}, nil, ErrInvalidTransition // or create ErrInvalidId if you prefer
	}
	if !StatusActive.Valid() {
		return Todo{}, nil, ErrInvalidTransition
	}

	t := Todo{
		ID:        p.ID,
		Title:     p.Title,
		Status:    StatusActive,
		Priority:  p.Priority,
		Tags:      p.Tags,
		DueDate:   p.DueDate,
		CreatedAt: p.Now,
		UpdatedAt: p.Now,
	}

	events := []Event{TodoCreated{ID: t.ID, OccurredAt: p.Now}}

	return t, events, nil
}

func (t Todo) ensureNotDeleted() error {
	if t.DeletedAt != nil {
		return ErrDeletedTodo
	}
	return nil
}

func (t Todo) ChangeTitle(newTitle Title, now time.Time) (Todo, []Event, error) {
	if err := t.ensureNotDeleted(); err != nil {
		return t, nil, err
	}
	if t.Title == newTitle {
		return t, nil, nil // idempotent
	}
	t.Title = newTitle
	t.UpdatedAt = now
	return t, []Event{TodoTitleChanged{ID: t.ID, Title: newTitle, OccurredAt: now}}, nil
}

func (t Todo) Complete(now time.Time) (Todo, []Event, error) {
	if err := t.ensureNotDeleted(); err != nil {
		return t, nil, err
	}
	switch t.Status {
	case StatusDone:
		return t, nil, nil // idempotent
	case StatusActive:
		t.Status = StatusDone
		t.CompletedAt = ptrTime(now)
		t.UpdatedAt = now
		return t, []Event{TodoCompleted{ID: t.ID, OccurredAt: now}}, nil
	case StatusArchived:
		return t, nil, ErrInvalidTransition
	default:
		return t, nil, ErrInvalidTransition
	}
}

func (t Todo) Archive(now time.Time) (Todo, []Event, error) {
	if err := t.ensureNotDeleted(); err != nil {
		return t, nil, err
	}
	switch t.Status {
	case StatusArchived:
		return t, nil, nil
	case StatusDone:
		t.Status = StatusArchived
		t.ArchivedAt = ptrTime(now)
		t.UpdatedAt = now
		return t, []Event{TodoArchived{ID: t.ID, OccurredAt: now}}, nil
	default:
		return t, nil, ErrInvalidTransition
	}
}

func (t Todo) Restore(now time.Time) (Todo, []Event, error) {
	if err := t.ensureNotDeleted(); err != nil {
		return t, nil, err
	}
	switch t.Status {
	case StatusArchived:
		t.Status = StatusActive
		t.ArchivedAt = nil
		t.UpdatedAt = now
		return t, []Event{TodoRestored{ID: t.ID, OccurredAt: now}}, nil
	case StatusActive:
		return t, nil, nil
	case StatusDone:
		return t, nil, ErrInvalidTransition
	default:
		return t, nil, ErrInvalidTransition
	}
}

func (t Todo) SoftDelete(now time.Time) (Todo, []Event, error) {
	if t.DeletedAt != nil {
		return t, nil, nil
	}
	t.DeletedAt = ptrTime(now)
	return t, []Event{TodoDeleted{ID: t.ID, OccurredAt: now}}, nil
}

func (t Todo) Reopen(now time.Time) (Todo, []Event, error) {
	if err := t.ensureNotDeleted(); err != nil {
		return t, nil, err
	}

	switch t.Status {
	case StatusActive:
		return t, nil, nil // idempotent
	case StatusDone:
		t.Status = StatusActive
		t.CompletedAt = nil
		t.UpdatedAt = now
		return t, []Event{TodoReopened{ID: t.ID, OccurredAt: now}}, nil
	case StatusArchived:
		return t, nil, ErrInvalidTransition
	default:
		return t, nil, ErrInvalidTransition
	}
}

func ptrTime(t time.Time) *time.Time { return &t }
