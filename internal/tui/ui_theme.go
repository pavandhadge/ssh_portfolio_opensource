package tui

import "github.com/charmbracelet/lipgloss"

var (
	colorBorder  = lipgloss.Color("240")
	colorText    = lipgloss.Color("252")
	colorMuted   = lipgloss.Color("245")
	colorPrimary = lipgloss.Color("15")
	colorOrange  = lipgloss.Color("#FFA726")
)

var (
	activeContentWidth = 80
	activeBodyHeight   = 12
)

var (
	appStyle = lipgloss.NewStyle().Padding(0, 0)

	headerBrandStyle = lipgloss.NewStyle().Foreground(colorOrange).Bold(true)
	headerMetaStyle  = lipgloss.NewStyle().Foreground(colorMuted)

	tabStyle       = lipgloss.NewStyle().Foreground(colorMuted)
	tabActiveStyle = lipgloss.NewStyle().Foreground(colorPrimary).Bold(true).Underline(true)

	cardStyle      = lipgloss.NewStyle().Border(lipgloss.HiddenBorder()).Padding(0, 0)
	pageTitleStyle = lipgloss.NewStyle().Foreground(colorOrange).Bold(true)
	titleStyle     = lipgloss.NewStyle().Foreground(colorPrimary).Bold(true)
	subtitleStyle  = lipgloss.NewStyle().Foreground(colorText)
	metaStyle      = lipgloss.NewStyle().Foreground(colorMuted).Faint(true)
	mutedStyle     = lipgloss.NewStyle().Foreground(colorMuted)

	listItemStyle       = lipgloss.NewStyle().Foreground(colorMuted)
	listItemActiveStyle = lipgloss.NewStyle().Foreground(colorPrimary).Bold(true)

	chipStyle       = lipgloss.NewStyle().Foreground(colorMuted).Padding(0, 1)
	linkStyle       = lipgloss.NewStyle().Foreground(colorMuted).Faint(true).Underline(true)
	linkActiveStyle = lipgloss.NewStyle().Foreground(colorOrange).Bold(true).Underline(true)
	quoteStyle      = lipgloss.NewStyle().Foreground(colorMuted).Italic(true)

	footerStyle        = lipgloss.NewStyle().Foreground(colorMuted).PaddingLeft(1)
	keyStyle           = lipgloss.NewStyle().Foreground(colorPrimary).Bold(true)
	loadingCursorStyle = lipgloss.NewStyle().Foreground(colorOrange).Bold(true)

	sepStyle = lipgloss.NewStyle().BorderStyle(lipgloss.NormalBorder()).BorderForeground(colorBorder)
)
