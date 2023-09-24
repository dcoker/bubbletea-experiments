package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"os"
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
	screen1 tea.Model
	screen2 tea.Model

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

				// If the target screen takes input from another screen, pass it as an event to the new screen's
				// Update() call.
				switch to {
				case 1:
					s1m, _ := m.screen1.(screen1model)
					m.screen2, _ = m.screen2.Update(onScreen2BeforeRender{filenames: s1m.FilenameArray()})
					// sometimes new screens might want to do real work before rendering; if so, set .wait and
					// return a Cmd (faked below).
				}

				m.wait = "pwning newbs ..."
				return m, func() tea.Msg {
					return runScreenChange(m.screen, to)
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
		m.screen1, cmd = m.screen1.Update(msg)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	case 1:
		m.screen2, cmd = m.screen2.Update(msg)
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
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("err: %v", err)
		os.Exit(1)
	}

}
