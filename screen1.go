package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"os"
)

type onScreen1DatabaseQueryDone struct {
	results []string
}

func runScreen1DatabaseQuery() tea.Msg {
	fakeSlow()
	entries, err := os.ReadDir(".")
	if err != nil {
		panic(err)
	}
	var files []string
	for _, entry := range entries {
		files = append(files, entry.Name())
		if len(files) > 5 {
			break
		}
	}
	return onScreen1DatabaseQueryDone{results: files}
}

type screen1model struct {
	choices  []string
	cursor   int
	selected map[int]struct{}
	spinner  spinner.Model
}

func (m screen1model) FilenameArray() []string {
	out := []string{}
	for i, v := range m.choices {
		if _, ok := m.selected[i]; ok {
			out = append(out, v)
		}
	}
	return out
}

func screen1InitialModel() tea.Model {
	return screen1model{
		choices:  []string{},
		selected: make(map[int]struct{}),
		spinner:  spinner.New(spinner.WithSpinner(spinner.Meter)),
	}
}

func (m screen1model) Init() tea.Cmd {
	return tea.Batch(runScreen1DatabaseQuery, m.spinner.Tick)
}

func (m screen1model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case onScreen1DatabaseQueryDone:
		m.choices = msg.results
	case tea.KeyMsg:
		switch msg.String() {
		case "up":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down":
			if m.cursor < len(m.choices) {
				m.cursor++
			}
		case "c", "r":
			m.choices = []string{}
			return m, runScreen1DatabaseQuery
		case "enter":
			_, ok := m.selected[m.cursor]
			if ok {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = struct{}{}
			}
		}
	}
	var cmd tea.Cmd
	m.spinner, cmd = m.spinner.Update(msg)
	return m, cmd
}

func (m screen1model) View() string {
	s := ""
	if len(m.choices) == 0 {
		return fmt.Sprintf("%s waiting for table list\n", m.spinner.View())
	}

	s += "[enter]: toggle; up/down: move arrow; right/left: change screen, c: clear and re-query; q: quit.\n\n"
	for i, choice := range m.choices {
		cursor := " "
		if m.cursor == i {
			cursor = ">"

		}
		checked := " "
		if _, ok := m.selected[i]; ok {
			checked = "x"
		}
		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
	}
	return s
}
