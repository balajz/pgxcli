package completer_test

import (
	"testing"

	"github.com/balajz/pgxcli/internal/completer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadPgKeywords(t *testing.T) {
	keywords := completer.LoadPgKeywords()
	require.NotEmpty(t, keywords)
	require.Len(t, keywords, 511)
	assert.Contains(t, keywords, "SELECT")
	assert.Contains(t, keywords, "INSERT")
	assert.Contains(t, keywords, "UPDATE")
	assert.Contains(t, keywords, "DELETE")
}
