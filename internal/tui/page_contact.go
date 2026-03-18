package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/pavandhadge/ssh_portfolio_opensource/internal/platform"
)

func (c contactPage) View() string {
	width := activeContentWidth - 4
	if width < 28 {
		width = 28
	}
	if width > 72 {
		width = 72
	}
	return renderPaneList(c.labels, c.focused, width)
}

func (c contactPage) Update(msg tea.Msg) (contactPage, tea.Cmd) {
	total := 5
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			c.focused = (c.focused - 1 + total) % total
		case "down", "j":
			c.focused = (c.focused + 1) % total
		case "enter":
			switch c.focused {
			case 0:
				return c, platform.OpenTargetCmd(c.email)
			case 1:
				return c, platform.OpenTargetCmd(c.github)
			case 2:
				return c, platform.OpenTargetCmd(c.linkedin)
			case 3:
				return c, platform.OpenTargetCmd(c.twitter)
			case 4:
				return c, platform.OpenTargetCmd(c.portfolio)
			}
		}
	}
	return c, nil
}
