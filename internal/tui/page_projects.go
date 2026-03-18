package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/pavandhadge/ssh_portfolio_opensource/internal/platform"
)

func (p projectsPage) View() string {
	if len(p.projects) == 0 {
		return mutedStyle.Render("No projects added yet.")
	}
	leftW, _, _ := paneMetrics()
	list := renderPaneList(p.labels, p.cursor, leftW)
	return renderTwoPaneListDetail(list, p.viewport.View())
}

func (p projectsPage) Update(msg tea.Msg) (projectsPage, tea.Cmd) {
	if len(p.projects) == 0 {
		return p, nil
	}
	var cmd tea.Cmd
	needSync := false
	prevCursor := p.cursor
	prevVersion := p.projects[p.cursor].version
	if p.projectFocus {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "left", "h", "right", "l", "tab", "enter":
				p.projects[p.cursor], cmd = p.projects[p.cursor].Update(msg)
				if p.projects[p.cursor].version != prevVersion {
					needSync = true
				}
			case "shift+tab", "esc", "backspace":
				p.projectFocus = false
			case "up", "k", "down", "j", "pgup", "pgdown", "ctrl+b", "ctrl+f":
				var vcmd tea.Cmd
				p.viewport, vcmd = p.viewport.Update(msg)
				cmd = tea.Batch(cmd, vcmd)
			}
		case tea.MouseMsg:
			var vcmd tea.Cmd
			p.viewport, vcmd = p.viewport.Update(msg)
			cmd = tea.Batch(cmd, vcmd)
		}
	} else if msg, ok := msg.(tea.KeyMsg); ok {
		switch msg.String() {
		case "tab", "enter", "right", "l":
			p.projectFocus = true
		case "up", "k":
			p.cursor = (p.cursor - 1 + len(p.projects)) % len(p.projects)
		case "down", "j":
			p.cursor = (p.cursor + 1) % len(p.projects)
		}
	}

	if p.cursor != prevCursor {
		p.viewport.GotoTop()
		needSync = true
	}

	rightW, rightH := detailPaneSize()
	if p.viewport.Width != rightW || p.viewport.Height != rightH {
		needSync = true
	}
	if p.detailCursor != p.cursor || p.detailVer != p.projects[p.cursor].version {
		needSync = true
	}
	if needSync {
		p.viewport = syncViewport(p.viewport, p.projects[p.cursor].render(rightW), rightW, rightH)
		p.detailCursor = p.cursor
		p.detailVer = p.projects[p.cursor].version
	}
	return p, cmd
}

func (pc projectComponent) render(contentW int) string {
	contentW -= 2
	if contentW < 24 {
		contentW = 24
	}
	return lipgloss.JoinVertical(
		lipgloss.Left,
		titleStyle.Render(pc.name),
		renderWrappedChips(pc.stack, contentW),
		"",
		renderLabeledURL("github:", pc.github, pc.highlightIdx == 0),
		renderLabeledURL("link:", pc.link, pc.highlightIdx == 1),
		"",
		lipgloss.NewStyle().Width(contentW).Render(pc.description),
	)
}

func (pc projectComponent) View() string {
	rightW, _ := detailPaneSize()
	return pc.render(rightW)
}

func (pc projectComponent) Update(msg tea.Msg) (projectComponent, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k", "left", "h":
			if pc.highlightIdx != 0 {
				pc.highlightIdx = 0
				pc.version++
			}
			if pc.highlightIdx == 0 && pc.github == "" && pc.link != "" {
				pc.highlightIdx = 1
				pc.version++
			}
		case "down", "j", "right", "l":
			if pc.highlightIdx != 1 {
				pc.highlightIdx = 1
				pc.version++
			}
			if pc.highlightIdx == 1 && pc.link == "" && pc.github != "" {
				pc.highlightIdx = 0
				pc.version++
			}
		case "tab":
			if pc.highlightIdx == 0 {
				pc.highlightIdx = 1
			} else {
				pc.highlightIdx = 0
			}
			if pc.highlightIdx == 0 && pc.github == "" && pc.link != "" {
				pc.highlightIdx = 1
			}
			if pc.highlightIdx == 1 && pc.link == "" && pc.github != "" {
				pc.highlightIdx = 0
			}
			pc.version++
		case "enter":
			target := ""
			if pc.highlightIdx == 0 {
				target = pc.github
				if target == "" {
					target = pc.link
				}
			} else {
				target = pc.link
				if target == "" {
					target = pc.github
				}
			}
			return pc, platform.OpenTargetCmd(target)
		}
	}
	return pc, nil
}
