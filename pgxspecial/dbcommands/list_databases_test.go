package dbcommands_test

import (
	"context"
	"testing"

	"github.com/balajz/pgxcli/pgxspecial/dbcommands"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
)

func TestListDatabases(t *testing.T) {
	db := connectTestDB(t)
	defer db.(*pgxpool.Pool).Close()

	pattern := ""
	verbose := false

	res, err := dbcommands.ListDatabases(context.Background(), db, pattern, verbose)
	if err != nil {
		t.Fatalf("ListDatabases failed: %v", err)
	}

	result := RequiresRowResult(t, res)

	defer result.Rows.Close()

	fds := result.Rows.FieldDescriptions()
	assert.NotNil(t, fds)

	columnsExpected := []string{
		"name",
		"owner",
		"encoding",
		"collate",
		"ctype",
		"access_privileges",
	}
	assert.Equal(t, columnsExpected, getColumnNames(fds), "Column names do not match expected")
	// expecting 6 columns: Name Owner Encoding Collate Ctype Access privileges
	assert.Len(t, fds, 6)

	var allRows []map[string]interface{}
	allRows, err = RowsToMaps(result.Rows)
	if err != nil {
		t.Fatalf("Failed to read rows: %v", err)
	}
	assert.True(t, containsDB(allRows, "template0"))
	assert.True(t, containsDB(allRows, "template1"))
	assert.True(t, containsDB(allRows, "postgres"))
}

func TestListDatabasesVerbose(t *testing.T) {
	db := connectTestDB(t)
	defer db.(*pgxpool.Pool).Close()

	pattern := ""
	verbose := true

	res, err := dbcommands.ListDatabases(context.Background(), db, pattern, verbose)
	if err != nil {
		t.Fatalf("ListDatabases failed: %v", err)
	}
	result := RequiresRowResult(t, res)

	defer result.Rows.Close()

	fds := result.Rows.FieldDescriptions()
	assert.NotNil(t, fds)

	columnsExpected := []string{
		"name",
		"owner",
		"encoding",
		"collate",
		"ctype",
		"access_privileges",
		"size",
		"Tablespace",
		"description",
	}
	assert.Equal(t, columnsExpected, getColumnNames(fds), "Column names do not match expected")
	// expecting 9 columns
	assert.Len(t, fds, 9)

	var allRows []map[string]interface{}
	allRows, err = RowsToMaps(result.Rows)
	if err != nil {
		t.Fatalf("Failed to read rows: %v", err)
	}
	assert.True(t, containsDB(allRows, "template0"))
	assert.True(t, containsDB(allRows, "template1"))
	assert.True(t, containsDB(allRows, "postgres"))
}

func TestListDatabaseWithExactPattern(t *testing.T) {
	db := connectTestDB(t)
	defer db.(*pgxpool.Pool).Close()

	pattern := "postgres"
	verbose := false

	res, err := dbcommands.ListDatabases(context.Background(), db, pattern, verbose)
	if err != nil {
		t.Fatalf("ListDatabases failed: %v", err)
	}
	result := RequiresRowResult(t, res)
	defer result.Rows.Close()

	fds := result.Rows.FieldDescriptions()
	assert.NotNil(t, fds)

	columnsExpected := []string{
		"name",
		"owner",
		"encoding",
		"collate",
		"ctype",
		"access_privileges",
	}
	assert.Equal(t, columnsExpected, getColumnNames(fds), "Column names do not match expected")
	// expecting 6 columns: Name Owner Encoding Collate Ctype Access privileges
	assert.Len(t, fds, 6)

	var allRows []map[string]interface{}
	allRows, err = RowsToMaps(result.Rows)
	if err != nil {
		t.Fatalf("Failed to read rows: %v", err)
	}
	assert.Len(t, allRows, 1, "Expected only one database matching the pattern")
	assert.False(t, containsDB(allRows, "template0"))
	assert.False(t, containsDB(allRows, "template1"))
	assert.True(t, containsDB(allRows, "postgres"))
}

func TestListDatabaseWithPattern(t *testing.T) {
	db := connectTestDB(t)
	defer db.(*pgxpool.Pool).Close()

	pattern := "templ*"
	verbose := false

	res, err := dbcommands.ListDatabases(context.Background(), db, pattern, verbose)
	if err != nil {
		t.Fatalf("ListDatabases failed: %v", err)
	}

	result := RequiresRowResult(t, res)
	defer result.Rows.Close()

	fds := result.Rows.FieldDescriptions()
	assert.NotNil(t, fds)

	columnsExpected := []string{
		"name",
		"owner",
		"encoding",
		"collate",
		"ctype",
		"access_privileges",
	}
	assert.Equal(t, columnsExpected, getColumnNames(fds), "Column names do not match expected")
	// expecting 6 columns: Name Owner Encoding Collate Ctype Access privileges
	assert.Len(t, fds, 6)

	var allRows []map[string]interface{}
	allRows, err = RowsToMaps(result.Rows)
	if err != nil {
		t.Fatalf("Failed to read rows: %v", err)
	}
	assert.Len(t, allRows, 2, "Expected only one database matching the pattern")
	assert.True(t, containsDB(allRows, "template0"))
	assert.True(t, containsDB(allRows, "template1"))
	assert.False(t, containsDB(allRows, "postgres"))
}

func TestListDatabaseWithNoMatchingPattern(t *testing.T) {
	db := connectTestDB(t)
	defer db.(*pgxpool.Pool).Close()

	pattern := "pastgres" // typo intentional
	verbose := false

	res, err := dbcommands.ListDatabases(context.Background(), db, pattern, verbose)
	if err != nil {
		t.Fatalf("ListDatabases failed: %v", err)
	}

	result := RequiresRowResult(t, res)

	defer result.Rows.Close()

	fds := result.Rows.FieldDescriptions()
	assert.NotNil(t, fds)

	columnsExpected := []string{
		"name",
		"owner",
		"encoding",
		"collate",
		"ctype",
		"access_privileges",
	}
	assert.Equal(t, columnsExpected, getColumnNames(fds), "Column names do not match expected")
	// expecting 6 columns: Name Owner Encoding Collate Ctype Access privileges
	assert.Len(t, fds, 6)

	var allRows []map[string]interface{}
	allRows, err = RowsToMaps(result.Rows)
	if err != nil {
		t.Fatalf("Failed to read rows: %v", err)
	}
	assert.Len(t, allRows, 0, "Expected no database matching the pattern")
}
