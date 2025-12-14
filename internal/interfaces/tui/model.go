package tui

import "github.com/rojanmagar2001/gotodo/internal/application/queries"

type Model struct {
	app App

	// UI state
	todos []queries.TodoDTO
	err   error
	ready bool
}

func NewModel(app App) Model {
	return Model{app: app}
}
