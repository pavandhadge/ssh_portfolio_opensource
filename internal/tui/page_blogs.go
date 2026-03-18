package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/pavandhadge/ssh_portfolio_opensource/internal/platform"
)

func (b blogsPage) View() string {
	if len(b.posts) == 0 {
		return mutedStyle.Render("No blog posts added yet.")
	}

	leftW, _, _ := paneMetrics()
	list := renderPaneList(b.labels, b.cursor, leftW)
	return renderTwoPaneListDetail(list, b.viewport.View())
}

func (b blogsPage) Update(msg tea.Msg) (blogsPage, tea.Cmd) {
	if len(b.posts) == 0 {
		return b, nil
	}

	var cmd tea.Cmd
	needSync := false
	prevCursor := b.cursor
	prevVersion := b.posts[b.cursor].version
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if b.blogFocus {
			switch msg.String() {
			case "left", "h", "right", "l", "tab":
				b.posts[b.cursor], cmd = b.posts[b.cursor].Update(msg)
				if b.posts[b.cursor].version != prevVersion {
					needSync = true
				}
			case "enter":
				b.posts[b.cursor], cmd = b.posts[b.cursor].Update(msg)
				if b.posts[b.cursor].version != prevVersion {
					needSync = true
				}
			case "shift+tab", "esc", "backspace":
				b.blogFocus = false
			case "up", "k", "down", "j", "pgup", "pgdown", "ctrl+b", "ctrl+f":
				var vcmd tea.Cmd
				b.viewport, vcmd = b.viewport.Update(msg)
				cmd = tea.Batch(cmd, vcmd)
			}
		} else {
			switch msg.String() {
			case "tab", "enter", "right", "l":
				b.blogFocus = true
			case "up", "k":
				b.cursor = (b.cursor - 1 + len(b.posts)) % len(b.posts)
			case "down", "j":
				b.cursor = (b.cursor + 1) % len(b.posts)
			}
		}
	case tea.MouseMsg:
		if b.blogFocus {
			var vcmd tea.Cmd
			b.viewport, vcmd = b.viewport.Update(msg)
			cmd = tea.Batch(cmd, vcmd)
		}
	}

	if b.cursor != prevCursor {
		b.viewport.GotoTop()
		needSync = true
	}

	rightW, rightH := detailPaneSize()
	if b.viewport.Width != rightW || b.viewport.Height != rightH {
		needSync = true
	}
	if b.detailCursor != b.cursor || b.detailVersion != b.posts[b.cursor].version {
		needSync = true
	}
	if needSync {
		b.viewport = syncViewport(b.viewport, b.posts[b.cursor].render(rightW), rightW, rightH)
		b.detailCursor = b.cursor
		b.detailVersion = b.posts[b.cursor].version
	}
	return b, cmd
}

func (bp blogPost) render(contentW int) string {
	contentW -= 2
	if contentW < 24 {
		contentW = 24
	}
	meta := bp.date
	if bp.readingTime != "" {
		if meta != "" {
			meta += " | "
		}
		meta += bp.readingTime
	}
	linkLine := mutedStyle.Render("no links available")
	if bp.url != "" {
		linkLine = renderLabeledURL("url:", bp.url, bp.highlightIdx == 0)
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		titleStyle.Render(bp.title),
		metaStyle.Render(meta),
		"",
		linkLine,
		"",
		lipgloss.NewStyle().Width(contentW).Render(bp.summary),
	)
}

func (bp blogPost) View() string {
	rightW, _ := detailPaneSize()
	return bp.render(rightW)
}

func (bp blogPost) Update(msg tea.Msg) (blogPost, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k", "left", "h", "right", "l", "down", "j", "tab":
			if bp.highlightIdx != 0 {
				bp.highlightIdx = 0
				bp.version++
			}
		case "enter":
			return bp, platform.OpenTargetCmd(bp.url)
		}
	}
	return bp, nil
}
