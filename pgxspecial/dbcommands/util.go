package dbcommands

import (
	"strings"
)

//nolint:gocyclo
func sqlNamePattern(pattern string) (schema, table string) {
	inQuotes := false
	var buf strings.Builder
	var schemaBuf *string

	for i := 0; i < len(pattern); i++ {
		c := pattern[i]

		switch {
		case c == '"':
			if inQuotes && i+1 < len(pattern) && pattern[i+1] == '"' {
				buf.WriteByte('"')
				i++
			} else {
				inQuotes = !inQuotes
			}

		case !inQuotes && c >= 'A' && c <= 'Z':
			buf.WriteByte(c + 32)

		case !inQuotes && c == '*':
			buf.WriteString(".*")

		case !inQuotes && c == '?':
			buf.WriteByte('.')

		case !inQuotes && c == '.':
			s := buf.String()
			schemaBuf = &s
			buf.Reset()

		default:
			if c == '$' || (inQuotes && strings.ContainsRune("|*+?()[]{}.^\\", rune(c))) {
				buf.WriteByte('\\')
			}
			buf.WriteByte(c)
		}
	}

	if buf.Len() > 0 {
		table = "^(" + buf.String() + ")$"
	}
	if schemaBuf != nil {
		schema = "^(" + *schemaBuf + ")$"
	}

	return schema, table
}
