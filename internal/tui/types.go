package tui

import "github.com/charmbracelet/bubbles/viewport"

const (
	tabOverview = iota
	tabProjects
	tabExperience
	tabContact
	tabBlogs
	tabHelp
)

type job struct {
	title           string
	tags            []string
	body            string
	companyName     string
	startDate       string
	endDate         string
	companyLocation string
	jobMode         string
	companyURL      string
	highlightIdx    int
	version         int
}

type projectComponent struct {
	name         string
	description  string
	stack        []string
	link         string
	github       string
	stars        int
	highlightIdx int
	version      int
}

type blogPost struct {
	title        string
	date         string
	url          string
	readingTime  string
	summary      string
	highlightIdx int
	version      int
}

type overviewPage struct {
	bio      string
	status   string
	location string
	skills   []string
	quote    string
}

type helpPage struct{}

type experiencePage struct {
	jobs            []job
	labels          []string
	cursor          int
	experienceFocus bool
	viewport        viewport.Model
	detailCursor    int
	detailVersion   int
}

type projectsPage struct {
	projects     []projectComponent
	labels       []string
	cursor       int
	projectFocus bool
	viewport     viewport.Model
	detailCursor int
	detailVer    int
}

type blogsPage struct {
	posts         []blogPost
	labels        []string
	cursor        int
	blogFocus     bool
	viewport      viewport.Model
	detailCursor  int
	detailVersion int
}

type contactPage struct {
	email     string
	github    string
	linkedin  string
	twitter   string
	portfolio string
	labels    []string

	focused int
}

type mainModel struct {
	pages     []string
	activeTab int

	width        int
	height       int
	ready        bool
	loadingBlink bool
	loadingDone  bool

	overview   overviewPage
	help       helpPage
	experience experiencePage
	projects   projectsPage
	blogs      blogsPage
	contact    contactPage
}
