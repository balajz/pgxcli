package dbcommands_test

import (
	"context"
	"testing"

	"github.com/balajz/pgxcli/pgxspecial/dbcommands"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
)

func TestListSchemas(t *testing.T) {
	db := connectTestDB(t)
	defer db.(*pgxpool.Pool).Close()

	pattern := ""
	verbose := false

	schemaNames := []string{"test_schema1", "test_schema2"}
	for _, schema := range schemaNames {
		CreateSchema(t, context.Background(), db.(*pgxpool.Pool), schema)
		defer DropSchema(t, context.Background(), db.(*pgxpool.Pool), schema)
	}

	res, err := dbcommands.ListSchemas(context.Background(), db, pattern, verbose)
	if err != nil {
		t.Fatalf("ListSchemas failed: %v", err)
	}
	result := RequiresRowResult(t, res)

	defer result.Rows.Close()

	fds := result.Rows.FieldDescriptions()
	assert.NotNil(t, fds)

	columnsExpected := []string{
		"name",
		"owner",
	}
	assert.Equal(t, columnsExpected, getColumnNames(fds), "Column names do not match expected")
	assert.Len(t, fds, 2)

	var allRows []map[string]interface{}
	allRows, err = RowsToMaps(result.Rows)
	if err != nil {
		t.Fatalf("Failed to read rows: %v", err)
	}
	assert.GreaterOrEqual(t, len(allRows), 2, "Expected at least two schemas matching the pattern")
	for _, schema := range schemaNames {
		assert.True(t, containsByField(allRows, "name", schema), "Expected schema %s not found", schema)
	}
}

func TestListSchemasWithPattern(t *testing.T) {
	db := connectTestDB(t)
	defer db.(*pgxpool.Pool).Close()

	pattern := "test_schema*"
	verbose := false

	schemaNames := []string{"test_schema1", "test_schema2"}
	for _, schema := range schemaNames {
		CreateSchema(t, context.Background(), db.(*pgxpool.Pool), schema)
		defer DropSchema(t, context.Background(), db.(*pgxpool.Pool), schema)
	}

	res, err := dbcommands.ListSchemas(context.Background(), db, pattern, verbose)
	if err != nil {
		t.Fatalf("ListSchemas failed: %v", err)
	}
	result := RequiresRowResult(t, res)

	defer result.Rows.Close()

	fds := result.Rows.FieldDescriptions()
	assert.NotNil(t, fds)

	columnsExpected := []string{
		"name",
		"owner",
	}
	assert.Equal(t, columnsExpected, getColumnNames(fds), "Column names do not match expected")
	assert.Len(t, fds, 2)

	var allRows []map[string]interface{}
	allRows, err = RowsToMaps(result.Rows)
	if err != nil {
		t.Fatalf("Failed to read rows: %v", err)
	}
	assert.GreaterOrEqual(t, len(allRows), 2, "Expected at least two schemas matching the pattern")
	for _, schema := range schemaNames {
		assert.True(t, containsByField(allRows, "name", schema), "Expected schema %s not found", schema)
	}
}

func TestListSchemasWithNoMatchingPattern(t *testing.T) {
	db := connectTestDB(t)
	defer db.(*pgxpool.Pool).Close()

	pattern := "non_existing_schema"
	verbose := false

	res, err := dbcommands.ListSchemas(context.Background(), db, pattern, verbose)
	if err != nil {
		t.Fatalf("ListSchemas failed: %v", err)
	}
	result := RequiresRowResult(t, res)

	defer result.Rows.Close()

	fds := result.Rows.FieldDescriptions()
	assert.NotNil(t, fds)

	columnsExpected := []string{
		"name",
		"owner",
	}
	assert.Equal(t, columnsExpected, getColumnNames(fds), "Column names do not match expected")
	assert.Len(t, fds, 2)

	var allRows []map[string]interface{}
	allRows, err = RowsToMaps(result.Rows)
	if err != nil {
		t.Fatalf("Failed to read rows: %v", err)
	}
	assert.Len(t, allRows, 0, "Expected no schemas matching the pattern")
}

func TestListSchemasVerbose(t *testing.T) {
	db := connectTestDB(t)
	defer db.(*pgxpool.Pool).Close()

	pattern := ""
	verbose := true

	schemaNames := []string{"test_schema1", "test_schema2"}
	for _, schema := range schemaNames {
		CreateSchema(t, context.Background(), db.(*pgxpool.Pool), schema)
		defer DropSchema(t, context.Background(), db.(*pgxpool.Pool), schema)
	}

	res, err := dbcommands.ListSchemas(context.Background(), db, pattern, verbose)
	if err != nil {
		t.Fatalf("ListSchemas failed: %v", err)
	}
	result := RequiresRowResult(t, res)

	defer result.Rows.Close()

	fds := result.Rows.FieldDescriptions()
	assert.NotNil(t, fds)

	columnsExpected := []string{
		"name",
		"owner",
		"access_privileges",
		"description",
	}
	assert.Equal(t, columnsExpected, getColumnNames(fds), "Column names do not match expected")
	assert.Len(t, fds, 4)

	var allRows []map[string]interface{}
	allRows, err = RowsToMaps(result.Rows)
	if err != nil {
		t.Fatalf("Failed to read rows: %v", err)
	}
	assert.GreaterOrEqual(t, len(allRows), 2, "Expected at least two schemas matching the pattern")
}

func TestListSchemasWithPatternVerbose(t *testing.T) {
	db := connectTestDB(t)
	defer db.(*pgxpool.Pool).Close()

	pattern := "test_schema*"
	verbose := true

	schemaNames := []string{"test_schema1", "test_schema2"}
	for _, schema := range schemaNames {
		CreateSchema(t, context.Background(), db.(*pgxpool.Pool), schema)
		defer DropSchema(t, context.Background(), db.(*pgxpool.Pool), schema)
	}

	res, err := dbcommands.ListSchemas(context.Background(), db, pattern, verbose)
	if err != nil {
		t.Fatalf("ListSchemas failed: %v", err)
	}
	result := RequiresRowResult(t, res)

	defer result.Rows.Close()

	fds := result.Rows.FieldDescriptions()
	assert.NotNil(t, fds)

	columnsExpected := []string{
		"name",
		"owner",
		"access_privileges",
		"description",
	}
	assert.Equal(t, columnsExpected, getColumnNames(fds), "Column names do not match expected")
	assert.Len(t, fds, 4)

	var allRows []map[string]interface{}
	allRows, err = RowsToMaps(result.Rows)
	if err != nil {
		t.Fatalf("Failed to read rows: %v", err)
	}
	assert.GreaterOrEqual(t, len(allRows), 2, "Expected at least two schemas matching the pattern")

}
func TestListSchemasWithNoMatchingPatternVerbose(t *testing.T) {
	db := connectTestDB(t)
	defer db.(*pgxpool.Pool).Close()

	pattern := "non_existing_schema"
	verbose := true
	res, err := dbcommands.ListSchemas(context.Background(), db, pattern, verbose)
	if err != nil {
		t.Fatalf("ListSchemas failed: %v", err)
	}
	result := RequiresRowResult(t, res)

	defer result.Rows.Close()

	fds := result.Rows.FieldDescriptions()
	assert.NotNil(t, fds)

	columnsExpected := []string{
		"name",
		"owner",
		"access_privileges",
		"description",
	}
	assert.Equal(t, columnsExpected, getColumnNames(fds), "Column names do not match expected")
	assert.Len(t, fds, 4)

	var allRows []map[string]interface{}
	allRows, err = RowsToMaps(result.Rows)
	if err != nil {
		t.Fatalf("Failed to read rows: %v", err)
	}
	assert.Len(t, allRows, 0, "Expected no schemas matching the pattern")
}
