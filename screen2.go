package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"log"
	"strings"
)

type Screen2InputArgs struct{ filenames []string }

type screen2model struct {
	filenames []string
	error     string
}

func (m screen2model) Validate(args Screen2InputArgs) (bool, string) {
	log.Printf("screen2model.Validate %+v", args)
	if len(args.filenames) == 0 {
		return false, "select file plz"
	}
	m.filenames = append(m.filenames, args.filenames...)
	return true, ""
}

func (m screen2model) Accept(args Screen2InputArgs) ScreenModel[Screen2InputArgs] {
	m.filenames = args.filenames
	return m
}

func screen2InitialModel() ScreenModel[Screen2InputArgs] {
	return screen2model{
		filenames: []string{},
		error:     "",
	}
}

func (m screen2model) Init() tea.Cmd {
	return nil
}

func (m screen2model) Update(_ tea.Msg) (tea.Model, tea.Cmd) {
	log.Printf("screen2model Update: %+v", m)
	return m, nil
}

func (m screen2model) View() string {
	return fmt.Sprintf("Filename selected: %s", strings.Join(m.filenames, ", "))
}
