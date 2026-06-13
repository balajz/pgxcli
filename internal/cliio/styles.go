package cliio

import "charm.land/lipgloss/v2"

var (
	// Primary
	Heading = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#9333EA")).
		Bold(true)

	Accent = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#C026D3"))

	Error = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#F87171")).
		MarginRight(2).
		Bold(true)

	Warning = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#B084F5")).
		Bold(true)

	Info = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#A855F7"))

	Success = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#8BDAA0")).
		Bold(true)

	Debug = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#7C6FA6"))

	// Text
	Detail = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#C4B5FD"))

	Dim = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#6B6680"))
)
