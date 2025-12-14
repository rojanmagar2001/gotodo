package tui

import (
	"strings"
)

func (m Model) View() string {
	if m.err != nil {
		return "Error: " + m.err.Error() + "\n\nPress q to quit.\n"
	}

	var b strings.Builder
	b.WriteString("Todo (Milestone 6)\n")
	b.WriteString("------------------\n\n")

	if len(m.todos) == 0 {
		b.WriteString("(no todos yet)\n")
	} else {
		for _, td := range m.todos {
			b.WriteString("- ")
			b.WriteString(td.Title)
			b.WriteString(" [")
			b.WriteString(td.Status)
			b.WriteString("]\n")
		}
	}

	b.WriteString("\nq: quit\n")
	return b.String()
}
