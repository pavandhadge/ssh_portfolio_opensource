package tui

import tea "github.com/charmbracelet/bubbletea"

func NewModel() tea.Model {
	return initialModel()
}
