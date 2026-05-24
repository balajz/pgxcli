package dbcommands

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSqlNamePattern(t *testing.T) {
	tests := []struct {
		name           string
		pattern        string
		expectedSchema string
		expectedTable  string
	}{
		{
			name:           "simple table",
			pattern:        "users",
			expectedSchema: "",
			expectedTable:  "^(users)$",
		},
		{
			name:           "schema and table",
			pattern:        "public.users",
			expectedSchema: "^(public)$",
			expectedTable:  "^(users)$",
		},
		{
			name:           "uppercase conversion",
			pattern:        "Users",
			expectedSchema: "",
			expectedTable:  "^(users)$",
		},
		{
			name:           "quoted uppercase preserved",
			pattern:        `"Users"`,
			expectedSchema: "",
			expectedTable:  "^(Users)$",
		},
		{
			name:           "wildcard *",
			pattern:        "user*",
			expectedSchema: "",
			expectedTable:  "^(user.*)$",
		},
		{
			name:           "wildcard ?",
			pattern:        "user?",
			expectedSchema: "",
			expectedTable:  "^(user.)$",
		},
		{
			name:           "dot in quotes",
			pattern:        `"schema.table"`,
			expectedSchema: "",
			expectedTable:  `^(schema\.table)$`,
		},
		{
			name:           "quotes inside quotes",
			pattern:        `"User""Name"`,
			expectedSchema: "",
			expectedTable:  `^(User"Name)$`,
		},
		{
			name:           "complex schema and table with wildcards",
			pattern:        "Schema*.Table?",
			expectedSchema: "^(schema.*)$",
			expectedTable:  "^(table.)$",
		},
		{
			name:           "regex special chars in quotes",
			pattern:        `"foo$bar"`,
			expectedSchema: "",
			expectedTable:  `^(foo\$bar)$`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			schema, table := sqlNamePattern(tt.pattern)
			assert.Equal(t, tt.expectedSchema, schema)
			assert.Equal(t, tt.expectedTable, table)
		})
	}
}
