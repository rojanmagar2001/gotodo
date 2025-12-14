package main

import (
	"fmt"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/rojanmagar2001/gotodo/internal/application/commands"
	"github.com/rojanmagar2001/gotodo/internal/application/queries"
	"github.com/rojanmagar2001/gotodo/internal/infrastructure/clock"
	"github.com/rojanmagar2001/gotodo/internal/infrastructure/events"
	"github.com/rojanmagar2001/gotodo/internal/infrastructure/idgen"
	"github.com/rojanmagar2001/gotodo/internal/infrastructure/jsonstore"
	"github.com/rojanmagar2001/gotodo/internal/infrastructure/logging"
	"github.com/rojanmagar2001/gotodo/internal/interfaces/tui"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "seed" {
		if err := runSeedCommand(os.Args[2:]); err != nil {
			fmt.Fprintln(os.Stderr, "seed error:", err)
			os.Exit(1)
		}
		return
	}

	logger := logging.New()

	// storage path (for now: ~/.gotodo/todos.json)
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	dir := filepath.Join(home, ".gotodo")
	_ = os.MkdirAll(dir, 0o700)
	dbPath := filepath.Join(dir, "todos.json")

	repo := jsonstore.NewRepository(dbPath)

	clk := clock.RealClock{}
	ids := idgen.RandomIDGen{}
	pub := events.LogPublisher{L: logger}

	// undo := commands.NewUndoManager()

	// Commands
	add := commands.AddTodo{Repo: repo, Clock: clk, IDGen: ids, Publisher: pub}
	complete := commands.CompleteTodo{Repo: repo, Clock: clk, Publisher: pub}

	// Queries
	list := queries.ListTodos{Repo: repo}
	get := queries.GetTodo{Repo: repo}
	stats := queries.Stats{Repo: repo, Clock: clk}

	app := tui.App{
		Add:      add,
		Complete: complete,
		List:     list,
		Get:      get,
		Stats:    stats,
	}

	p := tea.NewProgram(tui.NewModel(app), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
