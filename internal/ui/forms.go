// Package ui contains terminal UI flows used by the CLI.
package ui

import (
	"errors"
	"fmt"
	"image/color"
	"strconv"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/huh/v2"
	"charm.land/lipgloss/v2"
)

const maxWidth = 80

type styles struct {
	Base,
	HeaderText,
	Status,
	StatusHeader,
	Highlight,
	ErrorHeaderText,
	Help lipgloss.Style

	Red, Indigo, Green color.Color
}

func newStyles(hasDarkBg bool) *styles {
	var (
		s         = styles{}
		lightDark = lipgloss.LightDark(hasDarkBg)
	)

	s.Red = lightDark(lipgloss.Color("#FE5F86"), lipgloss.Color("#FE5F86"))
	s.Indigo = lightDark(lipgloss.Color("#5A56E0"), lipgloss.Color("#7571F9"))
	s.Green = lightDark(lipgloss.Color("#02BA84"), lipgloss.Color("#02BF87"))
	s.Base = lipgloss.NewStyle().
		Padding(1, 4, 0, 1)
	s.HeaderText = lipgloss.NewStyle().
		Foreground(s.Indigo).
		Bold(true).
		Padding(0, 2)
	s.Status = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(s.Indigo).
		PaddingLeft(1).
		MarginTop(1)
	s.StatusHeader = lipgloss.NewStyle().
		Foreground(s.Green).
		Bold(true)
	s.Highlight = lipgloss.NewStyle().
		Foreground(lipgloss.Color("212"))
	s.ErrorHeaderText = s.HeaderText.
		Foreground(s.Red)
	s.Help = lipgloss.NewStyle().
		Foreground(lipgloss.Color("240"))
	return &s
}

var errFormAborted = errors.New("connection form aborted")

// ConnectionValues contains connection details collected from the interactive form.
type ConnectionValues struct {
	Database string
	Username string
	Port     string
	Host     string
	Password string
}

type model struct {
	styles    func(bool) *styles
	form      *huh.Form
	values    ConnectionValues
	hasDarkBg bool
	height    int
	termWidth int
}

func newModel(database, username, host, port string) *model {
	port = strings.TrimSpace(port)
	if port == "" || port == "0" {
		port = "5432"
	}

	m := &model{
		styles: newStyles,
		values: ConnectionValues{
			Database: database,
			Username: username,
			Host:     host,
			Port:     port,
		},
	}

	m.form = huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Key("database").
				Title("Database").
				Placeholder("postgres").
				Validate(huh.ValidateNotEmpty()).
				Value(&m.values.Database),
			huh.NewInput().
				Key("username").
				Title("Username").
				Validate(huh.ValidateNotEmpty()).
				Value(&m.values.Username),
			huh.NewInput().
				Key("port").
				Title("Port").
				Placeholder("5432").
				Validate(validatePort).
				Value(&m.values.Port),
			huh.NewInput().
				Key("host").
				Title("Host").
				Description("unix socket if empty").
				Placeholder("localhost").
				Value(&m.values.Host),
			huh.NewInput().
				Key("password").
				Title("Password").
				EchoMode(huh.EchoModePassword).
				Value(&m.values.Password),
		),
	).
		WithWidth(45).
		WithShowHelp(false).
		WithShowErrors(false)
	return m
}

func (m *model) Init() tea.Cmd {
	return m.form.Init()
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.BackgroundColorMsg:
		m.hasDarkBg = msg.IsDark()
	case tea.WindowSizeMsg:
		m.termWidth = msg.Width
		m.height = msg.Height
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Interrupt
		case "esc", "q":
			return m, tea.Quit
		}
	}

	var cmds []tea.Cmd

	// Process the form
	form, cmd := m.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.form = f
		cmds = append(cmds, cmd)
	}

	if m.form.State == huh.StateCompleted {
		// Quit when the form is done.
		cmds = append(cmds, tea.Quit)
	}

	return m, tea.Batch(cmds...)
}

func (m *model) View() tea.View {
	s := m.styles(m.hasDarkBg)

	switch m.form.State {
	case huh.StateCompleted:
		var b strings.Builder
		fmt.Fprintf(&b, "Connection details captured.\n\n")
		fmt.Fprintf(&b, "Database: %s\n", fallback(m.values.Database))
		fmt.Fprintf(&b, "Username: %s\n", fallback(m.values.Username))
		fmt.Fprintf(&b, "Port: %s\n", fallback(m.values.Port))
		fmt.Fprintf(&b, "Host: %s\n", fallback(m.values.Host))
		fmt.Fprintf(&b, "Password: %s", maskedPassword(m.values.Password))
		rendered := s.Status.Margin(0, 1).Padding(1, 2).Width(48).Render(b.String())
		view := tea.NewView(lipgloss.Place(m.termWidth, m.height, lipgloss.Center, lipgloss.Center, rendered))
		view.AltScreen = true
		return view
	default:
		// Orca (left side)
		orca := lipgloss.NewStyle().Margin(1, 4, 0, 0).Render(orcaView())

		// Form card (right side)
		v := strings.TrimSuffix(m.form.View(), "\n\n")
		card := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(s.Indigo).
			Padding(0, 1).
			Render(v)

		// Live DSN preview
		preview := lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")).
			Italic(true).
			MarginTop(1).
			Render(buildDSNPreview(
				m.form.GetString("database"),
				m.form.GetString("username"),
				m.form.GetString("host"),
				m.form.GetString("port"),
			))

		rightPanel := lipgloss.NewStyle().Margin(1, 0).Render(
			lipgloss.JoinVertical(lipgloss.Left, card, preview),
		)

		errors := m.form.Errors()
		header := m.appBoundaryView("PostgreSQL Connection Form")
		if len(errors) > 0 {
			header = m.appErrorBoundaryView(m.errorView())
		}
		body := lipgloss.JoinHorizontal(lipgloss.Center, orca, rightPanel)

		footer := m.appBoundaryView(m.form.Help().ShortHelpView(m.form.KeyBinds()))
		if len(errors) > 0 {
			footer = m.appErrorBoundaryView("")
		}

		content := s.Base.Render(header + "\n" + body + "\n\n" + footer)
		view := tea.NewView(lipgloss.Place(m.termWidth, m.height, lipgloss.Center, lipgloss.Center, content))
		view.AltScreen = true
		return view
	}
}

func (m *model) errorView() string {
	var s strings.Builder
	for _, err := range m.form.Errors() {
		s.WriteString(err.Error())
	}
	return s.String()
}

func (m *model) appBoundaryView(text string) string {
	s := m.styles(m.hasDarkBg)
	w := m.termWidth - s.Base.GetHorizontalFrameSize()
	if w <= 0 {
		w = maxWidth
	}
	return lipgloss.PlaceHorizontal(
		w,
		lipgloss.Left,
		s.HeaderText.Render(text),
		lipgloss.WithWhitespaceChars("─"),
		lipgloss.WithWhitespaceStyle(lipgloss.NewStyle().Foreground(s.Indigo)),
	)
}

func (m *model) appErrorBoundaryView(text string) string {
	s := m.styles(m.hasDarkBg)
	w := m.termWidth - s.Base.GetHorizontalFrameSize()
	if w <= 0 {
		w = maxWidth
	}
	return lipgloss.PlaceHorizontal(
		w,
		lipgloss.Left,
		s.ErrorHeaderText.Render(text),
		lipgloss.WithWhitespaceChars("─"),
		lipgloss.WithWhitespaceStyle(lipgloss.NewStyle().Foreground(s.Red)),
	)
}

func (m *model) result() (ConnectionValues, error) {
	if m.form.State != huh.StateCompleted {
		return ConnectionValues{}, errFormAborted
	}

	port, err := strconv.Atoi(strings.TrimSpace(m.values.Port))
	if err != nil || port < 1 || port > 65535 {
		return ConnectionValues{}, fmt.Errorf("invalid port: %q", m.values.Port)
	}

	values := m.values
	return values, nil
}

// RunConnectionForm starts the interactive connection form and returns collected values.
func RunConnectionForm(database, username, host, port string) (ConnectionValues, error) {
	finalModel, err := tea.NewProgram(newModel(database, username, host, port)).Run()
	if err != nil {
		return ConnectionValues{}, err
	}

	m, ok := finalModel.(*model)
	if !ok {
		return ConnectionValues{}, fmt.Errorf("unexpected model type: %T", finalModel)
	}

	return m.result()
}

func fallback(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return "(None)"
	}
	return value
}

func maskedPassword(password string) string {
	if password == "" {
		return "(None)"
	}
	return strings.Repeat("•", len(password))
}

func buildDSNPreview(db, user, host, port string) string {
	if db == "" {
		db = "<database>"
	}
	if user == "" {
		user = "<username>"
	}
	if host == "" {
		host = "localhost"
	}
	if port == "" {
		port = "5432"
	}
	return fmt.Sprintf("postgresql://%s@%s:%s/%s", user, host, port, db)
}

func validatePort(v string) error {
	v = strings.TrimSpace(v)
	if v == "" {
		return errors.New("port is required")
	}
	p, err := strconv.Atoi(v)
	if err != nil {
		return errors.New("port must be a number")
	}
	if p < 1 || p > 65535 {
		return errors.New("port must be between 1 and 65535")
	}
	return nil
}
