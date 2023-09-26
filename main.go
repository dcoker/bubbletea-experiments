package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"log"
	"os"
	"reflect"
	"time"
)

// pretend to do work
func fakeSlow() {
	var duration, _ = time.ParseDuration("1s")
	time.Sleep(duration)
}

// tea.Msg
type onScreenChange struct {
	from    int
	to      int
	message string
}

// tea.Cmd: pretends to do some work before a screen is selected.
func runScreenChange(from int, to int) tea.Msg {
	fakeSlow()
	return onScreenChange{message: fmt.Sprintf("Welcome to screen %d", to+1), from: from, to: to}
}

type model struct {
	// track the current screen number
	screen     int
	numScreens int

	// screen 1
	screen1 ScreenModel[struct{}]
	screen2 ScreenModel[Screen2InputArgs]

	// informative messages for the header
	error       string
	information string

	// when commands are running (such as database queries), we display a spinner.
	wait    string
	spinner spinner.Model
}

func initialModel() model {
	outerModel := model{
		screen1:    screen1InitialModel(),
		screen2:    screen2InitialModel(),
		screen:     0,
		numScreens: 2,
		spinner:    spinner.New(spinner.WithSpinner(spinner.Dot)),
	}
	return outerModel

}

func (m model) Init() tea.Cmd {
	return tea.Batch(m.screen1.Init(), m.screen2.Init(), m.spinner.Tick)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	typ := reflect.TypeOf(msg).String()
	if typ != "spinner.TickMsg" {
		log.Printf("model.Update: %s; state = %+v", typ, m)
	}

	// Any msg not handled in this switch statement is passed to spinner or the screen-specific models.
	switch msg := msg.(type) {
	case onScreenChange:
		m.wait = ""
		m.information = msg.message
		m.screen = msg.to
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return m, tea.Quit
		case "left":
			if m.screen > 0 {
				m.wait = "reticulating splines ..."
				return m, func() tea.Msg {
					return runScreenChange(m.screen, m.screen-1)
				}
			}
			return m, nil
		case "right":
			if m.screen < m.numScreens-1 {
				to := m.screen + 1

				// First, validate that the next screen will accept the input data from the current state.
				var ok bool
				var msg string
				switch to {
				case 0:
					ok, msg = m.screen1.Validate(struct{}{})
				case 1:
					ok, msg = m.screen2.Validate(Screen2InputArgs{filenames: m.screen1.(screen1model).FilenameArray()})
				}
				// If the next screen doesn't accept the data, display a helpful error message.
				if !ok {
					m.error = msg
					return m, nil
				} else {
					// If the data passes validation, then tell the screen's models to accept it, and allow it to
					// start a Cmd.
					var cmd tea.Cmd
					switch to {
					case 0:
						m.screen1, cmd = m.screen1.Accept(struct{}{})
					case 1:
						m.screen2, cmd = m.screen2.Accept(Screen2InputArgs{filenames: m.screen1.(screen1model).FilenameArray()})
					}
					var cmds []tea.Cmd
					if cmd != nil {
						cmds = append(cmds, cmd)
					}
					cmds = append(cmds, func() tea.Msg { // fake work
						return runScreenChange(m.screen, to)
					})
					m.wait = "thinking real hard ..."
					return m, tea.Batch(cmds...)
				}
			}
			return m, nil
		}
	}

	// Any other events will be passed to the current screen and to the spinner.
	var cmds []tea.Cmd

	var cmd tea.Cmd
	m.spinner, cmd = m.spinner.Update(msg)
	if cmd != nil {
		cmds = append(cmds, cmd)
	}

	switch m.screen {
	case 0:
		update, cmd := m.screen1.(tea.Model).Update(msg)
		m.screen1 = update.(screen1model)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	case 1:
		update, cmd := m.screen2.(tea.Model).Update(msg)
		m.screen2 = update.(screen2model)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	if len(m.wait) > 0 {
		return fmt.Sprintf("%s Please wait: %s", m.spinner.View(), m.wait)
	}
	s := ""
	if len(m.error) > 0 {
		s += fmt.Sprintf("ERROR: %s\n", m.error)
	}
	if len(m.information) > 0 {
		s += fmt.Sprintf("INFORMATIVE MESSAGE: %s\n", m.information)
	}

	s += fmt.Sprintf("You are on screen %d of %d.\n", m.screen+1, m.numScreens)

	switch m.screen {
	case 0:
		s += m.screen1.View()
	case 1:
		s += m.screen2.View()
	}
	return s

}

func main() {
	file, err := tea.LogToFile("debug.log", "")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("err: %v", err)
		os.Exit(1)
	}
}
