package parser_test

import (
	"testing"

	"github.com/balaji01-4d/pgxcli/internal/parser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadPgKeywords(t *testing.T) {
	keywords := parser.LoadPgKeywords()
	require.NotEmpty(t, keywords)
	require.Len(t, keywords, 511)
	assert.Contains(t, keywords, "select")
	assert.Contains(t, keywords, "insert")
	assert.Contains(t, keywords, "update")
	assert.Contains(t, keywords, "delete")
}
