package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/pavandhadge/ssh_portfolio_opensource/internal/platform"
)

func (e experiencePage) View() string {
	if len(e.jobs) == 0 {
		return mutedStyle.Render("No experience added yet.")
	}

	leftW, _, _ := paneMetrics()
	list := renderPaneList(e.labels, e.cursor, leftW)
	return renderTwoPaneListDetail(list, e.viewport.View())
}

func (e experiencePage) Update(msg tea.Msg) (experiencePage, tea.Cmd) {
	if len(e.jobs) == 0 {
		return e, nil
	}

	var cmd tea.Cmd
	needSync := false
	prevCursor := e.cursor
	prevVersion := e.jobs[e.cursor].version
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if e.experienceFocus {
			switch msg.String() {
			case "left", "h", "right", "l", "tab":
				e.jobs[e.cursor], cmd = e.jobs[e.cursor].Update(msg)
				if e.jobs[e.cursor].version != prevVersion {
					needSync = true
				}
			case "enter":
				e.jobs[e.cursor], cmd = e.jobs[e.cursor].Update(msg)
				if e.jobs[e.cursor].version != prevVersion {
					needSync = true
				}
			case "shift+tab", "esc", "backspace":
				e.experienceFocus = false
			case "up", "k", "down", "j", "pgup", "pgdown", "ctrl+b", "ctrl+f":
				var vcmd tea.Cmd
				e.viewport, vcmd = e.viewport.Update(msg)
				cmd = tea.Batch(cmd, vcmd)
			}
		} else {
			switch msg.String() {
			case "tab", "enter", "right", "l":
				e.experienceFocus = true
			case "up", "k":
				e.cursor = (e.cursor - 1 + len(e.jobs)) % len(e.jobs)
			case "down", "j":
				e.cursor = (e.cursor + 1) % len(e.jobs)
			}
		}
	case tea.MouseMsg:
		if e.experienceFocus {
			var vcmd tea.Cmd
			e.viewport, vcmd = e.viewport.Update(msg)
			cmd = tea.Batch(cmd, vcmd)
		}
	}

	if e.cursor != prevCursor {
		e.viewport.GotoTop()
		needSync = true
	}

	rightW, rightH := detailPaneSize()
	if e.viewport.Width != rightW || e.viewport.Height != rightH {
		needSync = true
	}
	if e.detailCursor != e.cursor || e.detailVersion != e.jobs[e.cursor].version {
		needSync = true
	}
	if needSync {
		e.viewport = syncViewport(e.viewport, e.jobs[e.cursor].render(rightW), rightW, rightH)
		e.detailCursor = e.cursor
		e.detailVersion = e.jobs[e.cursor].version
	}
	return e, cmd
}

func (j job) render(contentW int) string {
	contentW -= 2
	if contentW < 24 {
		contentW = 24
	}
	meta := j.startDate + " - " + j.endDate
	if j.jobMode != "" {
		meta += " | " + j.jobMode
	}
	if j.companyLocation != "" {
		meta += " | " + j.companyLocation
	}
	linkLine := mutedStyle.Render("no links available")
	if j.companyURL != "" {
		linkLine = renderLabeledURL("company:", j.companyURL, j.highlightIdx == 0)
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		titleStyle.Render(j.title),
		subtitleStyle.Bold(true).Render(j.companyName),
		metaStyle.Render(meta),
		renderWrappedChips(j.tags, contentW),
		"",
		linkLine,
		"",
		lipgloss.NewStyle().Foreground(colorText).Width(contentW).Render(j.body),
	)
}

func (j job) View() string {
	rightW, _ := detailPaneSize()
	return j.render(rightW)
}

func (j job) Update(msg tea.Msg) (job, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k", "left", "h", "right", "l", "down", "j", "tab":
			if j.highlightIdx != 0 {
				j.highlightIdx = 0
				j.version++
			}
		case "enter":
			return j, platform.OpenTargetCmd(j.companyURL)
		}
	}
	return j, nil
}
