package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"strings"
)

type onScreen2BeforeRender struct {
	filenames []string
}
type screen2model struct {
	filenames []string
	error     string
}

func screen2InitialModel() tea.Model {
	return screen2model{
		filenames: []string{},
		error:     "",
	}
}

func (m screen2model) Init() tea.Cmd {
	return nil
}

func (m screen2model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case onScreen2BeforeRender:
		if strings.Contains(strings.Join(msg.filenames, ","), ".sum") {
			m.error = "no, bad, no select .sum files"
		} else {
			m.filenames = msg.filenames
		}
	}
	return m, nil
}

func (m screen2model) View() string {
	if len(m.error) > 0 {
		return m.error + "\nHit [left] and fix your input."
	}

	return fmt.Sprintf("Filename selected: %s", strings.Join(m.filenames, ", "))
}
