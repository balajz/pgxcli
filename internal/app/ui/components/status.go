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
	IssueLink      string
	SeparatorStyle lipgloss.Style
	StatusBarStyle lipgloss.Style
}

func NewStatusModel(version, issueLink string) StatusModel {
	return StatusModel{
		Version:   version,
		IssueLink: issueLink,
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

	name := "pgxcli " + m.Version

	link := m.StatusBarStyle.
		Underline(true).
		Hyperlink(m.IssueLink).
		Render("Report Issue")

	separator := m.SeparatorStyle.Render(strings.Repeat("─", m.Width))

	innerWidth := m.Width - m.StatusBarStyle.GetHorizontalPadding()
	if innerWidth < 0 {
		innerWidth = 0
	}

	usedWidth := lipgloss.Width(name) + lipgloss.Width(link)
	paddingWidth := innerWidth - usedWidth
	if paddingWidth < 0 {
		paddingWidth = 0
	}
	padding := strings.Repeat(" ", paddingWidth)

	statusBar := m.StatusBarStyle.Width(m.Width).Render(name + padding + link)

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
