package ui

import (
	_ "embed"
	"fmt"
	"os"
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/muesli/termenv"
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

func orcaStr(out *termenv.Output) string {
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
		lineCol := out.Color(fmt.Sprintf("#%02x%02x%02x", r, g, b))

		for _, ch := range line {
			sb.WriteString(out.String(string(ch)).Foreground(lineCol).String())
		}
	}

	return sb.String()
}

func orcaView() string {
	return orcaStr(termenv.NewOutput(os.Stdout))
}

func BannerCmd(version string) tea.Cmd {
	out := termenv.NewOutput(os.Stdout)
	green := out.Color("#02BF87")

	str := fmt.Sprintf("%s\n  %s  %s\n\n",
		orcaStr(out),
		out.String("pgxcli v"+version).Foreground(green).Bold().String(),
		out.String("\\q to quit").Foreground(out.Color("240")).String(),
	)
	return tea.Printf("%s", str)
}
