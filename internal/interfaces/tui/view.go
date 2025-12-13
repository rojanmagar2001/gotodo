package tui

import "strings"

func (m Model) View() string {
	var b strings.Builder
	b.WriteString("Todo TUI (Milestone 0)\n")
	b.WriteString("----------------------\n\n")
	b.WriteString("Keys:\n")
	b.WriteString("  q / ctrl+c  quit\n\n")

	if !m.ready {
		b.WriteString("(waiting for terminal size...)\n")
	}
	return b.String()
}
