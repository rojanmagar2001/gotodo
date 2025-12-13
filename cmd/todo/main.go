package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rojanmagar2001/gotodo/internal/interfaces/tui"
)

func main() {
	p := tea.NewProgram(tui.NewModel(), tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
