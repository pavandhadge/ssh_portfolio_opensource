package tui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/pavandhadge/ssh_portfolio_opensource/internal/platform"
)

type loadingTickMsg struct{}
type loadingDoneMsg struct{}

func loadingTickCmd() tea.Cmd {
	return tea.Tick(450*time.Millisecond, func(time.Time) tea.Msg { return loadingTickMsg{} })
}

func loadingDoneCmd() tea.Cmd {
	return tea.Tick(2*time.Second, func(time.Time) tea.Msg { return loadingDoneMsg{} })
}

func (m mainModel) Init() tea.Cmd {
	return tea.Batch(loadingTickCmd(), loadingDoneCmd())
}

func (m mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	routeToPage := true

	switch msg := msg.(type) {
	case loadingTickMsg:
		if m.loadingDone {
			return m, nil
		}
		m.loadingBlink = !m.loadingBlink
		return m, loadingTickCmd()
	case loadingDoneMsg:
		m.loadingDone = true
		return m, nil
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.ready = true
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "alt+g":
			return m, platform.OpenTargetCmd("https://github.com/ur_profile")
		case "alt+l":
			return m, platform.OpenTargetCmd("https://www.in.linkedin.com/in/ur_profile")
		case "alt+p":
			return m, platform.OpenTargetCmd("https://ur_portfolio.site")
		case "o":
			m.activeTab = tabOverview
			routeToPage = false
		case "p":
			m.activeTab = tabProjects
			routeToPage = false
		case "e":
			m.activeTab = tabExperience
			routeToPage = false
		case "c":
			m.activeTab = tabContact
			routeToPage = false
		case "b":
			m.activeTab = tabBlogs
			routeToPage = false
		case "?":
			m.activeTab = tabHelp
			routeToPage = false
		case "esc":
			m.activeTab = tabOverview
			routeToPage = false
		}
	}

	if !routeToPage {
		m.refreshActiveTabViewport()
		return m, nil
	}

	switch m.activeTab {
	case tabOverview:
		m.overview, cmd = m.overview.Update(msg)
	case tabProjects:
		m.projects, cmd = m.projects.Update(msg)
	case tabExperience:
		m.experience, cmd = m.experience.Update(msg)
	case tabContact:
		m.contact, cmd = m.contact.Update(msg)
	case tabBlogs:
		m.blogs, cmd = m.blogs.Update(msg)
	case tabHelp:
		m.help, cmd = m.help.Update(msg)
	}

	return m, cmd
}

func (m *mainModel) refreshActiveTabViewport() {
	rightW, rightH := detailPaneSize()
	switch m.activeTab {
	case tabProjects:
		if len(m.projects.projects) == 0 {
			return
		}
		cur := m.projects.cursor
		ver := m.projects.projects[cur].version
		if m.projects.viewport.Width != rightW || m.projects.viewport.Height != rightH || m.projects.detailCursor != cur || m.projects.detailVer != ver {
			m.projects.viewport = syncViewport(m.projects.viewport, m.projects.projects[cur].render(rightW), rightW, rightH)
			m.projects.detailCursor = cur
			m.projects.detailVer = ver
		}
	case tabExperience:
		if len(m.experience.jobs) == 0 {
			return
		}
		cur := m.experience.cursor
		ver := m.experience.jobs[cur].version
		if m.experience.viewport.Width != rightW || m.experience.viewport.Height != rightH || m.experience.detailCursor != cur || m.experience.detailVersion != ver {
			m.experience.viewport = syncViewport(m.experience.viewport, m.experience.jobs[cur].render(rightW), rightW, rightH)
			m.experience.detailCursor = cur
			m.experience.detailVersion = ver
		}
	case tabBlogs:
		if len(m.blogs.posts) == 0 {
			return
		}
		cur := m.blogs.cursor
		ver := m.blogs.posts[cur].version
		if m.blogs.viewport.Width != rightW || m.blogs.viewport.Height != rightH || m.blogs.detailCursor != cur || m.blogs.detailVersion != ver {
			m.blogs.viewport = syncViewport(m.blogs.viewport, m.blogs.posts[cur].render(rightW), rightW, rightH)
			m.blogs.detailCursor = cur
			m.blogs.detailVersion = ver
		}
	}
}

func (m mainModel) View() string {
	if !m.ready || !m.loadingDone {
		return renderLoading(m)
	}

	var body string
	switch m.activeTab {
	case tabOverview:
		body = m.overview.View()
	case tabProjects:
		body = m.projects.View()
	case tabExperience:
		body = m.experience.View()
	case tabContact:
		body = m.contact.View()
	case tabBlogs:
		body = m.blogs.View()
	case tabHelp:
		body = m.help.View()
	}

	return renderMainUI(m, body)
}
