//go:build integration

package dbcommands_test

import (
	"context"
	"testing"

	"github.com/balajz/pgxcli/pgxspecial"
	"github.com/balajz/pgxcli/pgxspecial/dbcommands"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
)

func TestListExtensions(t *testing.T) {
	db := connectTestDB(t)
	defer db.(*pgxpool.Pool).Close()
	pattern := ""
	verbose := false

	res, err := dbcommands.ListExtensions(context.Background(), db, pattern, verbose)
	if err != nil {
		t.Fatalf("ListExtensions failed: %v", err)
	}
	result := RequiresRowResult(t, res)

	defer result.Rows.Close()

	fds := result.Rows.FieldDescriptions()
	if fds == nil {
		t.Fatalf("FieldDescriptions is nil")
	}

	columnsExpected := []string{
		"name",
		"version",
		"schema",
		"description",
	}
	// expecting 4 columns
	assert.Len(t, fds, 4, "Expected 4 columns")
	assert.Equal(t, columnsExpected, getColumnNames(fds), "Column names do not match expected")

	var allRows []map[string]interface{}
	allRows, err = RowsToMaps(result.Rows)
	if err != nil {
		t.Fatalf("Failed to read rows: %v", err)
	}
	assert.Greater(t, len(allRows), 0, "Expected at least one extension")
	assert.True(t, containsByField(allRows, "name", "plpgsql"), "Expected to find 'plpgsql' extension")
}

func TestListExtensionsWithPattern(t *testing.T) {
	db := connectTestDB(t)
	defer db.(*pgxpool.Pool).Close()
	pattern := "plpg*"
	verbose := false

	res, err := dbcommands.ListExtensions(context.Background(), db, pattern, verbose)
	if err != nil {
		t.Fatalf("ListExtensions failed: %v", err)
	}
	result := RequiresRowResult(t, res)

	defer result.Rows.Close()

	fds := result.Rows.FieldDescriptions()
	if fds == nil {
		t.Fatalf("FieldDescriptions is nil")
	}

	columnsExpected := []string{
		"name",
		"version",
		"schema",
		"description",
	}
	// expecting 4 columns
	assert.Len(t, fds, 4, "Expected 4 columns")
	assert.Equal(t, columnsExpected, getColumnNames(fds), "Column names do not match expected")

	var allRows []map[string]interface{}
	allRows, err = RowsToMaps(result.Rows)
	if err != nil {
		t.Fatalf("Failed to read rows: %v", err)
	}
	assert.Greater(t, len(allRows), 0, "Expected at least one extension")
	assert.True(t, containsByField(allRows, "name", "plpgsql"), "Expected to find 'plpgsql' extension")
}

func TestListExtensionsVerbose(t *testing.T) {
	db := connectTestDB(t)
	defer db.(*pgxpool.Pool).Close()
	pattern := ""
	verbose := true

	res, err := dbcommands.ListExtensions(context.Background(), db, pattern, verbose)
	if err != nil {
		t.Fatalf("ListExtensions failed: %v", err)
	}
	_, isExt := res.(pgxspecial.ExtensionVerboseListResult)
	if !isExt {
		t.Fatalf("Expected ExtensionVerboseListResult, got %T", res)
	}

	result, ok := res.(pgxspecial.ExtensionVerboseListResult)
	if !ok {
		t.Fatalf("Expected ExtensionVerboseListResult, got %T", res)
	}

	assert.Greater(t, len(result.Results), 0, "Expected at least one extension")
	found := false
	for _, ext := range result.Results {
		if ext.Name == "plpgsql" {
			found = true
			assert.Greater(t, len(ext.Description), 0, "Expected description for plpgsql extension")
			assert.Contains(t, ext.Description, "function plpgsql_call_handler()", "Unexpected description content")
			break
		}
	}

	assert.True(t, found, "Expected to find 'plpgsql' extension")
}
