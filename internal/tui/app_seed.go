package tui

func initialModel() mainModel {
	m := mainModel{
		pages:     []string{"", "", "", "", "", ""},
		activeTab: tabOverview,

		overview: overviewPage{
			bio:      "",
			status:   "",
			location: "",
			skills:   []string{"", "", "", "", "", "", "", ""},
			quote:    "",
		},
		experience: experiencePage{
			jobs: []job{
				{
					title:           "",
					companyName:     "",
					body:            contentText("", ""),
					startDate:       "",
					endDate:         "",
					companyLocation: "",
					jobMode:         "",
					tags:            []string{"", "", "", "", "", ""},
					companyURL:      "",
				},
				{
					title:           "",
					companyName:     "",
					body:            contentText("", ""),
					startDate:       "",
					endDate:         "",
					companyLocation: "",
					jobMode:         "",
					tags:            []string{"", "", "", "", ""},
					companyURL:      "",
				},
				{
					title:           "",
					companyName:     "",
					body:            contentText("", ""),
					startDate:       "",
					endDate:         "",
					companyLocation: "",
					jobMode:         "",
					tags:            []string{"", "", "", ""},
					companyURL:      "",
				},
				{
					title:           "",
					companyName:     "",
					body:            contentText("", ""),
					startDate:       "",
					endDate:         "",
					companyLocation: "",
					jobMode:         "",
					tags:            []string{"", "", "", ""},
					companyURL:      "",
				},
			},
		},
		projects: projectsPage{
			projects: []projectComponent{
				{
					name:        "",
					description: contentText("", ""),
					stack:       []string{"", "", "", "", ""},
					github:      "",
					link:        "",
					stars:       0,
				},
				{
					name:        "",
					description: contentText("", ""),
					stack:       []string{"", "", "", ""},
					github:      "",
					link:        "",
					stars:       0,
				},
				{
					name:        "",
					description: contentText("", ""),
					stack:       []string{"", "", "", ""},
					github:      "",
					link:        "",
					stars:       0,
				},
				{
					name:        "",
					description: contentText("", ""),
					stack:       []string{"", "", "", ""},
					github:      "",
					link:        "",
					stars:       0,
				},
				{
					name:        "",
					description: contentText("", ""),
					stack:       []string{"", "", ""},
					github:      "",
					link:        "",
					stars:       0,
				},
				{
					name:        "",
					description: contentText("", ""),
					stack:       []string{"", "", ""},
					github:      "",
					link:        "",
					stars:       0,
				},
				{
					name:        "",
					description: contentText("", ""),
					stack:       []string{"", "", ""},
					github:      "",
					link:        "",
					stars:       0,
				},
				{
					name:        "",
					description: contentText("", ""),
					stack:       []string{"", "", ""},
					github:      "",
					link:        "",
					stars:       0,
				},
				{
					name:        "",
					description: contentText("", ""),
					stack:       []string{"", "", "", ""},
					github:      "",
					link:        "",
					stars:       0,
				},
				{
					name:        "",
					description: contentText("", ""),
					stack:       []string{"", ""},
					github:      "",
					link:        "",
					stars:       0,
				},
				{
					name:        "",
					description: contentText("", ""),
					stack:       []string{"", ""},
					github:      "",
					link:        "",
					stars:       0,
				},
				{
					name:        "",
					description: contentText("", ""),
					stack:       []string{"", "", ""},
					github:      "",
					link:        "",
					stars:       0,
				},
			},
		},
		contact: contactPage{
			email:     "",
			github:    "",
			linkedin:  "",
			twitter:   "",
			portfolio: "",
		},
		blogs: blogsPage{
			posts: []blogPost{
				{
					title:       "",
					date:        "",
					readingTime: "",
					url:         "",
					summary:     contentText("", ""),
				},
				{
					title:       "",
					date:        "",
					readingTime: "",
					url:         "",
					summary:     contentText("", ""),
				},
				{
					title:       "",
					date:        "",
					readingTime: "",
					url:         "",
					summary:     contentText("", ""),
				},
				{
					title:       "",
					date:        "",
					readingTime: "",
					url:         "",
					summary:     contentText("", ""),
				},
				{
					title:       "",
					date:        "",
					readingTime: "",
					url:         "",
					summary:     contentText("", ""),
				},
				{
					title:       "",
					date:        "",
					readingTime: "",
					url:         "",
					summary:     contentText("", ""),
				},
				{
					title:       "",
					date:        "",
					readingTime: "",
					url:         "",
					summary:     contentText("", ""),
				},
				{
					title:       "",
					date:        "",
					readingTime: "",
					url:         "",
					summary:     contentText("", ""),
				},
				{
					title:       "",
					date:        "",
					readingTime: "",
					url:         "",
					summary:     contentText("", ""),
				},
				{
					title:       "",
					date:        "",
					readingTime: "",
					url:         "",
					summary:     contentText("", ""),
				},
				{
					title:       "",
					date:        "",
					readingTime: "",
					url:         "",
					summary:     contentText("", ""),
				},
				{
					title:       "",
					date:        "",
					readingTime: "",
					url:         "",
					summary:     contentText("", ""),
				},
			},
		},
	}
	initStaticLabels(&m)
	primePageViewports(&m)
	return m
}

func initStaticLabels(m *mainModel) {
	m.projects.labels = make([]string, 0, len(m.projects.projects))
	for _, project := range m.projects.projects {
		m.projects.labels = append(m.projects.labels, project.name)
	}

	m.experience.labels = make([]string, 0, len(m.experience.jobs))
	for _, j := range m.experience.jobs {
		m.experience.labels = append(m.experience.labels, j.title+""+j.companyName)
	}

	m.blogs.labels = make([]string, 0, len(m.blogs.posts))
	for _, post := range m.blogs.posts {
		m.blogs.labels = append(m.blogs.labels, post.title+""+post.date)
	}

	m.contact.labels = []string{
		"" + m.contact.email,
		"" + m.contact.github,
		"" + m.contact.linkedin,
		"" + m.contact.twitter,
		"" + m.contact.portfolio,
	}
}

func primePageViewports(m *mainModel) {
	projW, projH := detailPaneSize()
	if len(m.projects.projects) > 0 {
		m.projects.viewport = syncViewport(m.projects.viewport, m.projects.projects[m.projects.cursor].render(projW), projW, projH)
		m.projects.detailCursor = m.projects.cursor
		m.projects.detailVer = m.projects.projects[m.projects.cursor].version
	}
	if len(m.experience.jobs) > 0 {
		m.experience.viewport = syncViewport(m.experience.viewport, m.experience.jobs[m.experience.cursor].render(projW), projW, projH)
		m.experience.detailCursor = m.experience.cursor
		m.experience.detailVersion = m.experience.jobs[m.experience.cursor].version
	}
	if len(m.blogs.posts) > 0 {
		m.blogs.viewport = syncViewport(m.blogs.viewport, m.blogs.posts[m.blogs.cursor].render(projW), projW, projH)
		m.blogs.detailCursor = m.blogs.cursor
		m.blogs.detailVersion = m.blogs.posts[m.blogs.cursor].version
	}
}
