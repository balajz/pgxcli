package components

import (
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

// MessageUpdateMsg is used to update the message in the status bar
type MessageUpdateMsg struct {
	Message string
}

// MessageResetMsg is used to reset the message in the status bar
type MessageResetMsg struct{}

// StatusModel manages the status bar and the executing spinner
type StatusModel struct {
	Message        string
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
	case MessageUpdateMsg:
		m.Message = msg.Message
	case MessageResetMsg:
		m.Message = ""
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

	firstLine := m.StatusBarStyle.Faint(true).Render(m.Message)
	secondLine := m.StatusBarStyle.Render(lipgloss.Sprintf("%s%s%s", name, padding, link))

	statusBar := lipgloss.JoinVertical(
		lipgloss.Top,
		firstLine,
		secondLine,
	)

	return lipgloss.JoinVertical(
		lipgloss.Top,
		separator,
		statusBar,
	)
}

// StaticHeight returns the height of the footer.
func (m StatusModel) StaticHeight() int {
	separator := m.SeparatorStyle.Render(strings.Repeat("─", m.Width))

	// Top separator + Bottom separator + message line + version line
	return lipgloss.Height(separator)*2 + 2
}
