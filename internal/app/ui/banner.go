package ui

import (
	_ "embed"
	"fmt"
	"os"
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/term"
)

var issueLink = "https://github.com/balajz/pgxcli/issues"

var (
	primaryStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#8B5CF6")).
			Bold(true)

	secondaryStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#A78BFA"))

	mutedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#A1A1AA"))

	accentStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#C4B5FD")).
			Italic(true)

	linkStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7C3AED")).
			Underline(true)
)

//go:embed ascii.txt
var asciiArt string

func gradientColor(t float64) (r, g, b int) {
	type rgb = [3]float64
	pine := rgb{62, 143, 176}
	foam := rgb{156, 207, 216}
	iris := rgb{196, 167, 231}

	lerp := func(a, b, t float64) float64 { return a + t*(b-a) }

	var c1, c2 rgb
	var t2 float64
	if t <= 0.5 {
		c1, c2, t2 = pine, foam, t*2
	} else {
		c1, c2, t2 = foam, iris, (t-0.5)*2
	}

	return int(lerp(c1[0], c2[0], t2)),
		int(lerp(c1[1], c2[1], t2)),
		int(lerp(c1[2], c2[2], t2))
}

func orcaStr() string {
	lines := strings.Split(asciiArt, "\n")
	for len(lines) > 0 && strings.TrimSpace(lines[0]) == "" {
		lines = lines[1:]
	}
	for len(lines) > 0 && strings.TrimSpace(lines[len(lines)-1]) == "" {
		lines = lines[:len(lines)-1]
	}
	for i, line := range lines {
		lines[i] = strings.TrimRight(line, " ")
	}

	total := len(lines)
	var sb strings.Builder

	for i, line := range lines {
		if i > 0 {
			sb.WriteByte('\n')
		}

		t := 0.0
		if total > 1 {
			t = float64(i) / float64(total-1)
		}
		r, g, b := gradientColor(t)
		hexColor := lipgloss.Color(fmt.Sprintf("#%02x%02x%02x", r, g, b))
		style := lipgloss.NewStyle().Foreground(hexColor)

		sb.WriteString(style.Render(line))
	}

	return sb.String()
}

func Banner(version string) string {
	leftPane := orcaStr()

	rightPane := lipgloss.JoinVertical(
		lipgloss.Left,
		secondaryStyle.Render("welcome to ")+
			primaryStyle.Render("pgxcli ")+
			mutedStyle.Render("v"+version)+"\n\n",

		accentStyle.Render("Happy Postgresing!\n\n"),

		linkStyle.
			Hyperlink(issueLink).
			Render("Report Issues ↗"),
	)

	content := lipgloss.JoinHorizontal(
		lipgloss.Top,
		leftPane,
		"      ", // gap between art and text
		rightPane,
	)

	banner := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#8B5CF6")).
		Padding(1, 4).
		Render(content)

	w, _, err := term.GetSize(os.Stdout.Fd())
	if err == nil && w > 0 {
		banner = lipgloss.Place(w, lipgloss.Height(banner), lipgloss.Center, lipgloss.Top, banner)
	}

	return banner
}
