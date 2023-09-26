package main

import tea "github.com/charmbracelet/bubbletea"

// ScreenModel is a tea.Model plus behaviors to validate and accept state changes.
type ScreenModel[S any] interface {
	tea.Model
	// Validates that the screen will accept the parameter S.
	Validate(args S) (bool, string)
	// Accept merges the parameter's data into the model. Cannot fail, because we already promised to accept it during
	// Validate.
	Accept(args S) (ScreenModel[S], tea.Cmd)
}
