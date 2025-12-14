package tui

import (
	"github.com/rojanmagar2001/gotodo/internal/application/commands"
	"github.com/rojanmagar2001/gotodo/internal/application/queries"
)

type App struct {
	// Commands
	Add      commands.AddTodo
	Complete commands.CompleteTodo

	// Queries
	List  queries.ListTodos
	Get   queries.GetTodo
	Stats queries.Stats
}
