package tui

import (
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"
	"github.com/pavandhadge/ssh_portfolio_opensource/internal/platform"
)

var paneListStyleCache struct {
	rowW             int
	labelWrapStyle   lipgloss.Style
	activeRowStyle   lipgloss.Style
	inactiveRowStyle lipgloss.Style
}

func paneListStyles(rowW int) (lipgloss.Style, lipgloss.Style, lipgloss.Style) {
	if paneListStyleCache.rowW != rowW {
		paneListStyleCache.rowW = rowW
		paneListStyleCache.labelWrapStyle = lipgloss.NewStyle().Width(rowW - 3)
		paneListStyleCache.activeRowStyle = lipgloss.NewStyle().
			Width(rowW - 2).
			PaddingLeft(1).
			MarginBottom(1).
			BorderLeft(true).
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(colorOrange)
		paneListStyleCache.inactiveRowStyle = lipgloss.NewStyle().
			Width(rowW - 2).
			PaddingLeft(2).
			MarginBottom(1)
	}
	return paneListStyleCache.labelWrapStyle, paneListStyleCache.activeRowStyle, paneListStyleCache.inactiveRowStyle
}

func paneMetrics() (int, int, int) {
	total := activeContentWidth
	if total < 60 {
		total = 60
	}
	leftW := (total * 28) / 100
	if leftW < 22 {
		leftW = 22
	}
	if leftW > 34 {
		leftW = 34
	}
	rightW := total - leftW - 3
	if rightW < 24 {
		rightW = 24
	}
	bodyH := activeBodyHeight
	if bodyH < 6 {
		bodyH = 6
	}
	return leftW, rightW, bodyH
}

func syncViewport(vp viewport.Model, content string, width, height int) viewport.Model {
	if width < 1 {
		width = 1
	}
	if height < 1 {
		height = 1
	}
	if vp.Width == 0 || vp.Height == 0 {
		vp = viewport.New(width, height)
	}
	vp.Width = width
	vp.Height = height
	vp.SetContent(lipgloss.NewStyle().Width(width).Render(content))
	return vp
}

func detailPaneSize() (int, int) {
	_, rightW, bodyH := paneMetrics()
	rightH := bodyH - 1
	if rightH < 3 {
		rightH = 3
	}
	return rightW - 1, rightH
}

func renderPaneList(labels []string, cursor, leftW int) string {
	rows := make([]string, 0, len(labels))
	rowW := leftW - 3
	if rowW < 8 {
		rowW = 8
	}
	_, _, bodyH := paneMetrics()
	maxRows := bodyH - 2
	if maxRows < 4 {
		maxRows = 4
	}
	start := 0
	if len(labels) > maxRows {
		start = cursor - maxRows/2
		if start < 0 {
			start = 0
		}
		if start+maxRows > len(labels) {
			start = len(labels) - maxRows
		}
	}
	end := start + maxRows
	if end > len(labels) {
		end = len(labels)
	}

	if start > 0 {
		rows = append(rows, mutedStyle.Render("  ^ more"))
	}
	labelWrapStyle, activeRowStyle, inactiveRowStyle := paneListStyles(rowW)
	for i := start; i < end; i++ {
		label := labels[i]
		labelText := labelWrapStyle.Render(label)
		if i == cursor {
			// Border-left naturally spans full wrapped height of the active entry.
			rows = append(rows, activeRowStyle.Render(listItemActiveStyle.Render(labelText)))
		} else {
			rows = append(rows, inactiveRowStyle.Render(listItemStyle.Render(labelText)))
		}
	}
	if end < len(labels) {
		rows = append(rows, mutedStyle.Render("  v more"))
	}
	return lipgloss.JoinVertical(lipgloss.Left, rows...)
}

func renderTwoPaneListDetail(list, detail string) string {
	leftW, rightW, bodyH := paneMetrics()
	left := lipgloss.NewStyle().
		Width(leftW).
		MaxWidth(leftW).
		PaddingRight(1).
		Render(lipgloss.NewStyle().MaxWidth(leftW - 2).Render(list))
	sepH := bodyH - 6
	if sepH < 4 {
		sepH = 4
	}
	sep := lipgloss.Place(1, bodyH-1, lipgloss.Center, lipgloss.Center, strings.Repeat("│\n", sepH-1)+"│")
	right := lipgloss.NewStyle().
		Width(rightW).
		MaxWidth(rightW).
		PaddingLeft(1).
		Render(lipgloss.NewStyle().MaxWidth(rightW - 1).Render(detail))
	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		left,
		lipgloss.NewStyle().Foreground(colorBorder).Render(sep),
		right,
	)
}

func renderWrappedChips(items []string, width int) string {
	if len(items) == 0 {
		return ""
	}
	if width < 12 {
		width = 12
	}
	var b strings.Builder
	for i, item := range items {
		if i > 0 {
			b.WriteByte(' ')
		}
		b.WriteString(chipStyle.Render(item))
	}
	return lipgloss.NewStyle().Width(width).Render(b.String())
}

func renderLabeledURL(label, url string, active bool) string {
	labelPart := mutedStyle.Render(label + " ")
	if url == "" {
		return labelPart + mutedStyle.Render("-")
	}
	if active {
		return labelPart + platform.TerminalLinkLabel(linkActiveStyle.Render(url), url)
	}
	return labelPart + platform.TerminalLinkLabel(linkStyle.Render(url), url)
}
