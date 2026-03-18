package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func (o overviewPage) Update(msg tea.Msg) (overviewPage, tea.Cmd) { return o, nil }
func (h helpPage) Update(msg tea.Msg) (helpPage, tea.Cmd)         { return h, nil }

func (o overviewPage) View() string {
	contentW := activeContentWidth - 2
	if contentW < 24 {
		contentW = 24
	}
	skills := renderWrappedChips(o.skills, contentW)
	top := lipgloss.JoinVertical(
		lipgloss.Left,
		lipgloss.NewStyle().Foreground(colorOrange).Bold(true).Render("Open Source Developer"),
		mutedStyle.Render(o.location+"  |  "+o.status),
		"",
		lipgloss.NewStyle().Width(contentW).Render(o.bio),
		"",
		mutedStyle.Render("Skills"),
		skills,
	)
	quote := quoteStyle.Render("\"with love by pavan dhadge\"")
	bodyH := activeBodyHeight - 2
	if bodyH < 8 {
		bodyH = 8
	}
	space := bodyH - lipgloss.Height(top)
	if space < 2 {
		space = 2
	}
	return lipgloss.JoinVertical(
		lipgloss.Left,
		top,
		lipgloss.Place(contentW, space, lipgloss.Center, lipgloss.Bottom, quote),
	)
}

func (h helpPage) View() string {
	return lipgloss.JoinVertical(
		lipgloss.Left,
		titleStyle.Render("Help"),
		"",
		mutedStyle.Render("Navigation"),
		"o overview  |  p projects  |  e experience  |  c contact  |  b blogs  |  ? help",
		"",
		mutedStyle.Render("List & Focus"),
		"up/down select entry  |  tab/enter/right focus details  |  shift+tab or esc back to list",
		"",
		mutedStyle.Render("Detail Pane"),
		"up/down scroll  |  left/right/tab switch links (when available)  |  enter show link for local open",
		"",
		mutedStyle.Render("Session"),
		"q quit",
	)
}
