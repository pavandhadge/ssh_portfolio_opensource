package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func clamp(v, lo int) int {
	if v < lo {
		return lo
	}
	return v
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func rule(width int) string {
	if width <= 0 {
		return ""
	}
	return sepStyle.Width(width).
		BorderTop(true).
		BorderBottom(false).
		BorderLeft(false).
		BorderRight(false).
		Render("")
}

func renderTabs(pages []string, active int) string {
	var b strings.Builder
	for i, p := range pages {
		if i > 0 {
			b.WriteByte(' ')
		}
		if i == active {
			b.WriteString(tabActiveStyle.Render(p))
		} else {
			b.WriteString(tabStyle.Render(p))
		}
	}
	return b.String()
}

func renderLoading(m mainModel) string {
	cursor := " "
	if m.loadingBlink {
		cursor = "_"
	}
	line := headerBrandStyle.Render("ur_portfolio.site") + " " + loadingCursorStyle.Render(cursor)
	w, h := m.width, m.height
	if w <= 0 {
		w = 80
	}
	if h <= 0 {
		h = 24
	}
	return lipgloss.Place(w, h, lipgloss.Center, lipgloss.Center, line)
}

func pageInstructions(m mainModel) (string, string) {
	base := lipgloss.JoinHorizontal(
		lipgloss.Left,
		keyStyle.Render("o/p/e/c/b/?"),
		" switch tab  |  ",
		keyStyle.Render("q"),
		" quit",
	)

	switch m.activeTab {
	case tabOverview:
		return base, "overview summary section"
	case tabContact:
		return base, lipgloss.JoinHorizontal(lipgloss.Left, keyStyle.Render("up/down"), " select  |  ", keyStyle.Render("enter"), " show link for local open")
	case tabHelp:
		return base, "help and control reference"
	case tabProjects:
		if m.projects.projectFocus {
			return base, lipgloss.JoinHorizontal(lipgloss.Left, keyStyle.Render("up/down"), " scroll  |  ", keyStyle.Render("left/right/tab"), " change link  |  ", keyStyle.Render("enter"), " show link for local open  |  ", keyStyle.Render("shift+tab/esc"), " back")
		}
		return base, lipgloss.JoinHorizontal(lipgloss.Left, keyStyle.Render("up/down"), " select project  |  ", keyStyle.Render("tab/enter/right"), " focus details")
	case tabExperience:
		if m.experience.experienceFocus {
			return base, lipgloss.JoinHorizontal(lipgloss.Left, keyStyle.Render("up/down"), " scroll details  |  ", keyStyle.Render("enter"), " show company link for local open  |  ", keyStyle.Render("shift+tab/esc"), " back")
		}
		return base, lipgloss.JoinHorizontal(lipgloss.Left, keyStyle.Render("up/down"), " select role  |  ", keyStyle.Render("tab/enter/right"), " focus details")
	case tabBlogs:
		if m.blogs.blogFocus {
			return base, lipgloss.JoinHorizontal(lipgloss.Left, keyStyle.Render("up/down"), " scroll details  |  ", keyStyle.Render("enter"), " show blog link for local open  |  ", keyStyle.Render("shift+tab/esc"), " back")
		}
		return base, lipgloss.JoinHorizontal(lipgloss.Left, keyStyle.Render("up/down"), " select post  |  ", keyStyle.Render("tab/enter/right"), " focus details")
	default:
		return base, ""
	}
}

func renderMainUI(m mainModel, body string) string {
	h := clamp(m.height-1, 12)
	usableW := clamp(m.width, 48)
	contentW := min(usableW, 110)
	activeContentWidth = contentW

	top := lipgloss.JoinVertical(
		lipgloss.Left,
		lipgloss.JoinHorizontal(lipgloss.Left, headerBrandStyle.Render("terminal"), "  ", renderTabs(m.pages, m.activeTab)),
		headerMetaStyle.Render("# terminal portfolio"),
		rule(contentW),
	)

	ins1, ins2 := pageInstructions(m)
	footer := lipgloss.JoinVertical(lipgloss.Left, rule(contentW), footerStyle.Render(ins1), footerStyle.Render(ins2))

	used := lipgloss.Height(top) + lipgloss.Height(footer)
	bodyH := clamp(h-used-1, 4)
	activeBodyHeight = bodyH

	pageHead := pageTitleStyle.Render(strings.ToUpper(m.pages[m.activeTab]))
	main := lipgloss.JoinVertical(
		lipgloss.Left,
		pageHead,
		rule(contentW),
		"",
		cardStyle.Width(contentW).Height(bodyH-2).Render(body),
	)

	layout := lipgloss.JoinVertical(lipgloss.Left, top, "", main, "", footer)
	return lipgloss.NewStyle().
		Width(usableW).
		Height(h).
		Align(lipgloss.Center, lipgloss.Top).
		Render(appStyle.Width(contentW).Height(h).Render(layout))
}
