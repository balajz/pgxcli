package renderer

import (
	"errors"
	"strings"
	"testing"

	"github.com/balajz/pgxcli/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTable(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name         string
		data         Data
		wantContains []string
		wantErr      string
	}{
		{
			name: "renders rows and caption",
			data: newDummyTableData(
				[]string{"id", "name"},
				[][]any{{1, "alice"}, {2, "bob"}},
				"2 rows",
			),
			wantContains: []string{"id", "name", "alice", "bob", "2 rows"},
		},
		{
			name:         "renders empty row set with header",
			data:         newDummyTableData([]string{"id"}, [][]any{}, "0 rows"),
			wantContains: []string{"id", "0 rows"},
		},
		{
			name:    "returns row error",
			data:    newDummyTableDataWithError(errors.New("boom")),
			wantErr: "boom",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var out strings.Builder
			err := Table(tc.data, &out, &config.Config{})

			if tc.wantErr != "" {
				require.Error(t, err)
				assert.ErrorContains(t, err, tc.wantErr)
				return
			}

			require.NoError(t, err)
			rendered := out.String()
			require.NotEmpty(t, rendered)
			assertContainsFold(t, rendered, tc.wantContains...)
		})
	}
}
