package tui

import (
	"context"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rojanmagar2001/gotodo/internal/application/ports"
	"github.com/rojanmagar2001/gotodo/internal/application/queries"
)

type todosLoadedMsg struct {
	todos []queries.TodoDTO
	err   error
}

func (m Model) Init() tea.Cmd { return m.loadTodosCmd() }

func (m Model) loadTodosCmd() tea.Cmd {
	return func() tea.Msg {
		res := m.app.List.Execute(context.Background(), ports.ListSpec{})
		return todosLoadedMsg{todos: res.Value, err: res.Err}
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch x := msg.(type) {
	case tea.WindowSizeMsg:
		m.ready = true
	case todosLoadedMsg:
		m.todos = x.todos
		m.err = x.err
	case tea.KeyMsg:
		switch x.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}
	return m, nil
}
