package components

import (
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

// StatusModel manages the status bar and the executing spinner.
type StatusModel struct {
	Version        string
	Width          int
	SeparatorStyle lipgloss.Style
	StatusBarStyle lipgloss.Style
}

func NewStatusModel(version string) StatusModel {
	return StatusModel{
		Version: version,
	}
}

func (m StatusModel) Init() tea.Cmd {
	return nil
}

func (m StatusModel) Update(msg tea.Msg) (StatusModel, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
	}
	return m, cmd
}

// ViewFooter returns the static footer with separator and status bar.
func (m StatusModel) View() string {
	if m.Width == 0 {
		return ""
	}
	separator := m.SeparatorStyle.Render(strings.Repeat("─", m.Width))
	statusBar := m.StatusBarStyle.Width(m.Width).Render("pgxcli " + m.Version)

	return lipgloss.JoinVertical(
		lipgloss.Top,
		separator,
		statusBar,
	)
}

// StaticHeight returns the height of the footer.
func (m StatusModel) StaticHeight() int {
	separator := m.SeparatorStyle.Render(strings.Repeat("─", m.Width))
	statusBar := m.StatusBarStyle.Width(m.Width).Render("pgxcli " + m.Version)

	// Top separator + Bottom separator + Status bar
	return lipgloss.Height(separator)*2 + lipgloss.Height(statusBar)
}
