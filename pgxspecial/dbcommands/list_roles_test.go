package dbcommands_test

import (
	"context"
	"testing"

	"github.com/balajz/pgxcli/pgxspecial/dbcommands"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
)

func TestListRoles(t *testing.T) {
	db := connectTestDB(t)
	defer db.(*pgxpool.Pool).Close()

	pattern := ""
	verbose := false

	res, err := dbcommands.ListRoles(context.Background(), db, pattern, verbose)
	if err != nil {
		t.Fatalf("ListRoles failed: %v", err)
	}
	result := RequiresRowResult(t, res)

	defer result.Rows.Close()

	fds := result.Rows.FieldDescriptions()
	assert.NotNil(t, fds)

	columnsExpected := []string{
		"rolname",
		"rolsuper",
		"rolinherit",
		"rolcreaterole",
		"rolcreatedb",
		"rolcanlogin",
		"rolconnlimit",
		"rolvaliduntil",
		"memberof",
		"rolreplication",
	}
	assert.Equal(t, columnsExpected, getColumnNames(fds), "Column names do not match expected")
	// expecting 10 columns
	assert.Len(t, fds, 10)

	var allRows []map[string]interface{}
	allRows, err = RowsToMaps(result.Rows)
	if err != nil {
		t.Fatalf("Failed to read rows: %v", err)
	}

	var essentialDefaultRoles = []string{
		"postgres",
		"pg_monitor",
		"pg_read_all_data",
		"pg_write_all_data",
	}

	for _, role := range essentialDefaultRoles {
		assert.True(t, containsByField(allRows, "rolname", role), "Expected role %s not found", role)
	}
}

func TestListRolesWithPattern(t *testing.T) {
	db := connectTestDB(t)
	defer db.(*pgxpool.Pool).Close()

	pattern := "pg_w*"
	verbose := false

	res, err := dbcommands.ListRoles(context.Background(), db, pattern, verbose)
	if err != nil {
		t.Fatalf("ListRoles failed: %v", err)
	}
	result := RequiresRowResult(t, res)

	defer result.Rows.Close()

	fds := result.Rows.FieldDescriptions()
	assert.NotNil(t, fds)

	columnsExpected := []string{
		"rolname",
		"rolsuper",
		"rolinherit",
		"rolcreaterole",
		"rolcreatedb",
		"rolcanlogin",
		"rolconnlimit",
		"rolvaliduntil",
		"memberof",
		"rolreplication",
	}
	assert.Equal(t, columnsExpected, getColumnNames(fds), "Column names do not match expected")
	// expecting 10 columns
	assert.Len(t, fds, 10)

	var allRows []map[string]interface{}
	allRows, err = RowsToMaps(result.Rows)
	if err != nil {
		t.Fatalf("Failed to read rows: %v", err)
	}
	assert.Len(t, allRows, 2)
	var expectedRoles = []string{
		"pg_write_all_data",
		"pg_write_server_files",
	}
	for _, role := range expectedRoles {
		assert.True(t, containsByField(allRows, "rolname", role), "Expected role %s not found", role)
	}
}

func TestListRolesWithNoMatchingPattern(t *testing.T) {
	db := connectTestDB(t)
	defer db.(*pgxpool.Pool).Close()

	pattern := "pg_xwrite*" // intentional typo
	verbose := false

	res, err := dbcommands.ListRoles(context.Background(), db, pattern, verbose)
	if err != nil {
		t.Fatalf("ListRoles failed: %v", err)
	}
	result := RequiresRowResult(t, res)

	defer result.Rows.Close()

	fds := result.Rows.FieldDescriptions()
	assert.NotNil(t, fds)

	columnsExpected := []string{
		"rolname",
		"rolsuper",
		"rolinherit",
		"rolcreaterole",
		"rolcreatedb",
		"rolcanlogin",
		"rolconnlimit",
		"rolvaliduntil",
		"memberof",
		"rolreplication",
	}
	assert.Equal(t, columnsExpected, getColumnNames(fds), "Column names do not match expected")
	// expecting 10 columns
	assert.Len(t, fds, 10)

	var allRows []map[string]interface{}
	allRows, err = RowsToMaps(result.Rows)
	if err != nil {
		t.Fatalf("Failed to read rows: %v", err)
	}
	assert.Len(t, allRows, 0, "Expected no roles matching the pattern")
}

func TestListRolesWithPatternVerbose(t *testing.T) {
	db := connectTestDB(t)
	defer db.(*pgxpool.Pool).Close()

	pattern := "pg_w*"
	verbose := true

	res, err := dbcommands.ListRoles(context.Background(), db, pattern, verbose)
	if err != nil {
		t.Fatalf("ListRoles failed: %v", err)
	}
	result := RequiresRowResult(t, res)

	defer result.Rows.Close()

	fds := result.Rows.FieldDescriptions()
	assert.NotNil(t, fds)

	columnsExpected := []string{
		"rolname",
		"rolsuper",
		"rolinherit",
		"rolcreaterole",
		"rolcreatedb",
		"rolcanlogin",
		"rolconnlimit",
		"rolvaliduntil",
		"memberof",
		"description",
		"rolreplication",
	}
	assert.Equal(t, columnsExpected, getColumnNames(fds), "Column names do not match expected")
	// expecting 11 columns
	assert.Len(t, fds, 11)

	var allRows []map[string]interface{}
	allRows, err = RowsToMaps(result.Rows)
	if err != nil {
		t.Fatalf("Failed to read rows: %v", err)
	}
	assert.Len(t, allRows, 2)
	var expectedRoles = []string{
		"pg_write_all_data",
		"pg_write_server_files",
	}
	for _, role := range expectedRoles {
		assert.True(t, containsByField(allRows, "rolname", role), "Expected role %s not found", role)
	}
}

func TestListRolesWithNoMatchingPatternVerbose(t *testing.T) {
	db := connectTestDB(t)
	defer db.(*pgxpool.Pool).Close()

	pattern := "pg_xwrite*" // intentional typo
	verbose := true

	res, err := dbcommands.ListRoles(context.Background(), db, pattern, verbose)
	if err != nil {
		t.Fatalf("ListRoles failed: %v", err)
	}
	result := RequiresRowResult(t, res)

	defer result.Rows.Close()

	fds := result.Rows.FieldDescriptions()
	assert.NotNil(t, fds)

	columnsExpected := []string{
		"rolname",
		"rolsuper",
		"rolinherit",
		"rolcreaterole",
		"rolcreatedb",
		"rolcanlogin",
		"rolconnlimit",
		"rolvaliduntil",
		"memberof",
		"description",
		"rolreplication",
	}
	assert.Equal(t, columnsExpected, getColumnNames(fds), "Column names do not match expected")
	// expecting 11 columns
	assert.Len(t, fds, 11)

	var allRows []map[string]interface{}
	allRows, err = RowsToMaps(result.Rows)
	if err != nil {
		t.Fatalf("Failed to read rows: %v", err)
	}
	assert.Len(t, allRows, 0, "Expected no roles matching the pattern")
}
