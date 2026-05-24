package dbcommands_test

import (
	"context"
	"testing"

	"github.com/balajz/pgxcli/pgxspecial/dbcommands"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
)

func TestListPrivileges(t *testing.T) {
	db := connectTestDB(t).(*pgxpool.Pool)
	defer db.Close()

	pattern := ""
	pool := db
	ctx := context.Background()

	db.Exec(ctx, "CREATE TABLE test_tbl (id int)")
	db.Exec(ctx, "CREATE ROLE test_user")

	GrantPrivilege(t, ctx, pool, "SELECT", "test_tbl", "test_user")
	defer RevokePrivilege(t, ctx, pool, "SELECT", "test_tbl", "test_user")

	res, err := dbcommands.ListPrivileges(context.Background(), db, pattern, false)
	if err != nil {
		t.Fatalf("ListPrivileges failed: %v", err)
	}
	result := RequiresRowResult(t, res)

	defer result.Rows.Close()

	fds := result.Rows.FieldDescriptions()
	assert.NotNil(t, fds)

	columnsExpected := []string{
		"schema",
		"name",
		"type",
		"access_privileges",
		"column_privileges",
		"policies",
	}
	assert.Equal(t, columnsExpected, getColumnNames(fds), "Column names do not match expected")
	// expecting 6 columns
	assert.Len(t, fds, 6)

	var allRows []map[string]interface{}
	allRows, err = RowsToMaps(result.Rows)
	if err != nil {
		t.Fatalf("Failed to read rows: %v", err)
	}
	assert.Greater(t, len(allRows), 0, "Expected at least one privilege entry")
}
