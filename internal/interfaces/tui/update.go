package tui

import tea "github.com/charmbracelet/bubbletea"

func (m Model) Init() tea.Cmd { return nil }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case tea.WindowSizeMsg:
		m.ready = true
	case tea.KeyMsg:
		km := msg.(tea.KeyMsg)
		if km.String() == "q" || km.String() == "ctrl+c" {
			return m, tea.Quit
		}
	}
	return m, nil
}
