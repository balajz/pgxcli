// Package render contains helpers for formatting special-command results.
package render

import "github.com/jedib0t/go-pretty/v6/table"

// Tables renders a list of pretty tables using the provided style.
func Tables(tables []table.Writer, style table.Style) string {
	var str string
	for _, table := range tables {
		table.SetStyle(style)
		str += table.Render() + "\n"
	}
	return str
}
