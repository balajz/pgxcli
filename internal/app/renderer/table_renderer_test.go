package renderer

import (
	"strings"
	"testing"

	"github.com/balajz/pgxcli/internal/config"
	"github.com/balaji01-4d/pgxspecial"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testRowsResult struct {
	columns []string
	data    [][]any
}

func (r testRowsResult) ResultKind() pgxspecial.SpecialResultKind {
	return pgxspecial.ResultKindRows
}

func (r testRowsResult) Columns() []string {
	return r.columns
}

func (r testRowsResult) Data() [][]any {
	return r.data
}

func TestRowsResult(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name         string
		result       pgxspecial.SpecialCommandResult
		wantContains []string
		wantErr      string
	}{
		{
			name: "renders special rows",
			result: testRowsResult{
				columns: []string{"id", "name"},
				data:    [][]any{{1, "alice"}},
			},
			wantContains: []string{"id", "name", "1", "alice"},
		},
		{
			name:    "rejects unsupported type",
			result:  pgxspecial.ExtensionVerboseListResult{},
			wantErr: "invalid row result type",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got, err := RowsResult(tc.result, &config.Config{})
			if tc.wantErr != "" {
				require.Error(t, err)
				assert.ErrorContains(t, err, tc.wantErr)
				return
			}

			require.NoError(t, err)
			assertContainsFold(t, got, tc.wantContains...)
		})
	}
}

func TestDescribeTableResult(t *testing.T) {
	t.Parallel()

	view := "select * from users"
	testCases := []struct {
		name         string
		result       pgxspecial.SpecialCommandResult
		wantContains []string
		wantErr      string
	}{
		{
			name: "renders describe table sections",
			result: pgxspecial.DescribeTableListResult{
				Results: []pgxspecial.DescribeTableResult{{
					Columns: []string{"Column", "Type"},
					Data:    [][]string{{"id", "integer"}},
					TableMetaData: pgxspecial.TableFooterMeta{
						Indexes:        []string{"users_pkey"},
						ViewDefinition: &view,
					},
				}},
			},
			wantContains: []string{"Column", "Type", "id", "integer", "Indexes:", "users_pkey", "View definition:", view},
		},
		{
			name:    "rejects unsupported type",
			result:  testRowsResult{},
			wantErr: "invalid describe table result type",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got, err := DescribeTableResult(tc.result, &config.Config{})
			if tc.wantErr != "" {
				require.Error(t, err)
				assert.ErrorContains(t, err, tc.wantErr)
				return
			}

			require.NoError(t, err)
			assertContainsFold(t, got, tc.wantContains...)
		})
	}
}

func TestExtensionVerboseResult(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name         string
		result       pgxspecial.SpecialCommandResult
		wantContains []string
		wantErr      string
	}{
		{
			name: "renders extension verbose sections",
			result: pgxspecial.ExtensionVerboseListResult{
				Results: []pgxspecial.ExtensionVerboseResult{{
					Name:        "pgcrypto",
					Description: []string{"cryptographic functions"},
				}},
			},
			wantContains: []string{"Object Description", "cryptographic functions", "pgcrypto"},
		},
		{
			name:    "rejects unsupported type",
			result:  testRowsResult{},
			wantErr: "invalid extension verbose result type",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got, err := ExtensionVerboseResult(tc.result, &config.Config{})
			if tc.wantErr != "" {
				require.Error(t, err)
				assert.ErrorContains(t, err, tc.wantErr)
				return
			}

			require.NoError(t, err)
			assertContainsFold(t, got, tc.wantContains...)
		})
	}
}

func TestRenderTableFooter(t *testing.T) {
	t.Parallel()

	view := "select * from users"
	partitionKey := "LIST (tenant_id)"
	ownedBy := "postgres.user_id_seq"
	options := "fillfactor=80"
	server := "remote"
	fdwOptions := "(schema_name 'public')"
	tableType := "public.audit_log"
	hasOIDs := true

	testCases := []struct {
		name string
		meta pgxspecial.TableFooterMeta
		want string
	}{
		{
			name: "empty metadata",
			meta: pgxspecial.TableFooterMeta{},
			want: "",
		},
		{
			name: "renders populated metadata in order",
			meta: pgxspecial.TableFooterMeta{
				Indexes:          []string{"users_pkey"},
				CheckConstraints: []string{"users_name_check"},
				ViewDefinition:   &view,
				PartitionKey:     &partitionKey,
				HasOIDs:          &hasOIDs,
				Options:          &options,
				Server:           &server,
				FDWOptions:       &fdwOptions,
				OwnedBy:          &ownedBy,
				TypedTableOf:     &tableType,
			},
			want: strings.Join([]string{
				"Indexes:",
				"    users_pkey",
				"Check constraints:",
				"    users_name_check",
				"View definition:",
				"select * from users",
				"Partition key: LIST (tenant_id)",
				"Typed table of type: public.audit_log",
				"Has OIDs: yes",
				"Options: fillfactor=80",
				"Server: remote",
				"FDW Options: (schema_name 'public')",
				"Owned by: postgres.user_id_seq",
			}, "\n") + "\n",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := renderTableFooter(tc.meta)
			assert.Equal(t, tc.want, got)
		})
	}
}
