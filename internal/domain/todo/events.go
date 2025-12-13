package todo

import "time"

type Event interface {
	eventName() string
}

type TodoCreated struct {
	ID         TodoID
	OccurredAt time.Time
}

func (TodoCreated) eventName() string { return "todo.created" }

type TodoTitleChanged struct {
	ID         TodoID
	Title      Title
	OccurredAt time.Time
}

func (TodoTitleChanged) eventName() string { return "todo.title_changed" }

type TodoCompleted struct {
	ID         TodoID
	OccurredAt time.Time
}

func (TodoCompleted) eventName() string { return "todo.completed" }

type TodoReopened struct {
	ID         TodoID
	OccurredAt time.Time
}

func (TodoReopened) eventName() string { return "todo.reopened" }

type TodoArchived struct {
	ID         TodoID
	OccurredAt time.Time
}

func (TodoArchived) eventName() string { return "todo.archived" }

type TodoRestored struct {
	ID         TodoID
	OccurredAt time.Time
}

func (TodoRestored) eventName() string { return "todo.restored" }

type TodoDeleted struct {
	ID         TodoID
	OccurredAt time.Time
}

func (TodoDeleted) eventName() string { return "todo.deleted" }
