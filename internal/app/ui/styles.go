package ui

import "charm.land/lipgloss/v2"

type Styles struct {
	// UserInput styles the user's input query after enter key is pressed.
	UserInput lipgloss.Style
	// AppOutput styles the output from the application after executing the query.
	AppOutput lipgloss.Style
	// ErrorOutput styles the error output from the application after executing the query.
	ErrorOutput lipgloss.Style

	// UserInputSepartor styles the separator between user input from previous result.
	UserInputSepartor lipgloss.Style

	// InputSeparator is a top and bottom border for editline input area.
	InputSeparator lipgloss.Style

	// Spinner styles the spinner animation.
	Spinner        lipgloss.Style
	SpinnerCaption lipgloss.Style

	// StatusBar styles the status bar at the bottom.
	StatusBar lipgloss.Style

	// ClampNotice styles the notice when user input is clamped.
	ClampNotice lipgloss.Style
}

func DefaultStyles() Styles {
	return Styles{
		UserInput: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#A78BFA")),
		AppOutput: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#C4B5FD")),
		ErrorOutput: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF6B6B")),
		UserInputSepartor: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#B8A2FF")),
		InputSeparator: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#8B5CF6")),
		Spinner: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#C4B5FD")).
			PaddingLeft(2).
			Bold(true),
		SpinnerCaption: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#A78BFA")).
			Italic(true).
			Faint(true),
		StatusBar: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#C4B5FD")).
			Background(lipgloss.Color("#2A273F")).
			Padding(0, 1),
		ClampNotice: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#A78BFA")).
			Italic(true).
			Faint(true),
	}
}
