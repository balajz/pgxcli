package dbcommands_test

import (
	"context"
	"testing"

	"github.com/balajz/pgxcli/pgxspecial/dbcommands"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
)

func TestListDatatypes(t *testing.T) {
	db := connectTestDB(t)
	defer db.(*pgxpool.Pool).Close()

	ctx := context.Background()
	typeNames := []string{"mood_enum", "status_enum", "priority_enum"}

	for _, typeName := range typeNames {
		CreateDatatype(t, ctx, db.(*pgxpool.Pool), typeName)
		defer DropDatatype(t, ctx, db.(*pgxpool.Pool), typeName)
	}

	pattern := ""
	verbose := false

	res, err := dbcommands.ListDatatypes(context.Background(), db, pattern, verbose)
	if err != nil {
		t.Fatalf("ListDatatypes failed: %v", err)
	}
	result := RequiresRowResult(t, res)
	defer result.Rows.Close()

	fds := result.Rows.FieldDescriptions()
	assert.NotNil(t, fds)

	columnsExpected := []string{
		"schema",
		"name",
		"description",
	}
	assert.Equal(t, columnsExpected, getColumnNames(fds), "Column names do not match expected")
	// expecting 3 columns
	assert.Len(t, fds, 3)

	var allRows []map[string]interface{}
	allRows, err = RowsToMaps(result.Rows)
	if err != nil {
		t.Fatalf("Failed to read rows: %v", err)
	}
	for _, typeName := range typeNames {
		assert.True(t, containsByField(allRows, "name", typeName))
	}
}

func TestListDatatypesWithPattern(t *testing.T) {
	db := connectTestDB(t)
	defer db.(*pgxpool.Pool).Close()

	ctx := context.Background()
	typeNames := []string{"mood_enum", "status_enum", "priority_enum"}

	for _, typeName := range typeNames {
		CreateDatatype(t, ctx, db.(*pgxpool.Pool), typeName)
		defer DropDatatype(t, ctx, db.(*pgxpool.Pool), typeName)
	}

	pattern := "*_enum"
	verbose := false

	res, err := dbcommands.ListDatatypes(context.Background(), db, pattern, verbose)
	if err != nil {
		t.Fatalf("ListDatatypes failed: %v", err)
	}
	result := RequiresRowResult(t, res)
	defer result.Rows.Close()

	fds := result.Rows.FieldDescriptions()
	assert.NotNil(t, fds)

	columnsExpected := []string{
		"schema",
		"name",
		"description",
	}
	assert.Equal(t, columnsExpected, getColumnNames(fds), "Column names do not match expected")
	// expecting 3 columns
	assert.Len(t, fds, 3)

	var allRows []map[string]interface{}
	allRows, err = RowsToMaps(result.Rows)
	if err != nil {
		t.Fatalf("Failed to read rows: %v", err)
	}
	for _, typeName := range typeNames {
		assert.True(t, containsByField(allRows, "name", typeName))
	}
}

func TestListDatatypesWithNoMatchingPattern(t *testing.T) {
	db := connectTestDB(t)
	defer db.(*pgxpool.Pool).Close()

	ctx := context.Background()
	typeNames := []string{"mood_enum", "status_enum", "priority_enum"}

	for _, typeName := range typeNames {
		CreateDatatype(t, ctx, db.(*pgxpool.Pool), typeName)
		defer DropDatatype(t, ctx, db.(*pgxpool.Pool), typeName)
	}

	pattern := "type_xenum"
	verbose := false

	res, err := dbcommands.ListDatatypes(context.Background(), db, pattern, verbose)
	if err != nil {
		t.Fatalf("ListDatatypes failed: %v", err)
	}
	result := RequiresRowResult(t, res)
	defer result.Rows.Close()

	fds := result.Rows.FieldDescriptions()
	assert.NotNil(t, fds)

	columnsExpected := []string{
		"schema",
		"name",
		"description",
	}
	assert.Equal(t, columnsExpected, getColumnNames(fds), "Column names do not match expected")
	// expecting 3 columns
	assert.Len(t, fds, 3)

	var allRows []map[string]interface{}
	allRows, err = RowsToMaps(result.Rows)
	if err != nil {
		t.Fatalf("Failed to read rows: %v", err)
	}
	assert.Len(t, allRows, 0)
}

func TestListDatatypesVerbose(t *testing.T) {
	db := connectTestDB(t)
	defer db.(*pgxpool.Pool).Close()

	ctx := context.Background()
	typeNames := []string{"mood_enum", "status_enum", "priority_enum"}

	for _, typeName := range typeNames {
		CreateDatatype(t, ctx, db.(*pgxpool.Pool), typeName)
		defer DropDatatype(t, ctx, db.(*pgxpool.Pool), typeName)
	}

	pattern := ""
	verbose := true

	res, err := dbcommands.ListDatatypes(context.Background(), db, pattern, verbose)
	if err != nil {
		t.Fatalf("ListDatatypes failed: %v", err)
	}
	result := RequiresRowResult(t, res)
	defer result.Rows.Close()

	fds := result.Rows.FieldDescriptions()
	assert.NotNil(t, fds)

	expectedColumns := []string{
		"schema",
		"name",
		"internal_name",
		"size",
		"elements",
		"access_privileges",
		"description",
	}
	assert.Equal(t, expectedColumns, getColumnNames(fds), "Column names do not match expected")
	// expecting 7 columns
	assert.Len(t, fds, 7)
	var allRows []map[string]interface{}
	allRows, err = RowsToMaps(result.Rows)
	if err != nil {
		t.Fatalf("Failed to read rows: %v", err)
	}
	for _, typeName := range typeNames {
		assert.True(t, containsByField(allRows, "name", typeName))
	}

}

func TestListDatatypesVerboseWithPattern(t *testing.T) {
	db := connectTestDB(t)
	defer db.(*pgxpool.Pool).Close()

	ctx := context.Background()
	typeNames := []string{"mood_enum", "status_enum", "priority_enum"}

	for _, typeName := range typeNames {
		CreateDatatype(t, ctx, db.(*pgxpool.Pool), typeName)
		defer DropDatatype(t, ctx, db.(*pgxpool.Pool), typeName)
	}

	pattern := "*_enum"
	verbose := true

	res, err := dbcommands.ListDatatypes(context.Background(), db, pattern, verbose)
	if err != nil {
		t.Fatalf("ListDatatypes failed: %v", err)
	}
	result := RequiresRowResult(t, res)
	defer result.Rows.Close()

	fds := result.Rows.FieldDescriptions()
	assert.NotNil(t, fds)

	expectedColumns := []string{
		"schema",
		"name",
		"internal_name",
		"size",
		"elements",
		"access_privileges",
		"description",
	}
	assert.Equal(t, expectedColumns, getColumnNames(fds), "Column names do not match expected")
	// expecting 7 columns
	assert.Len(t, fds, 7)
	var allRows []map[string]interface{}
	allRows, err = RowsToMaps(result.Rows)
	if err != nil {
		t.Fatalf("Failed to read rows: %v", err)
	}
	for _, typeName := range typeNames {
		assert.True(t, containsByField(allRows, "name", typeName))
	}
}

func TestListDatatypesVerboseWithNoMatchingPattern(t *testing.T) {
	db := connectTestDB(t)
	defer db.(*pgxpool.Pool).Close()
	ctx := context.Background()
	typeNames := []string{"mood_enum", "status_enum", "priority_enum"}

	for _, typeName := range typeNames {
		CreateDatatype(t, ctx, db.(*pgxpool.Pool), typeName)
		defer DropDatatype(t, ctx, db.(*pgxpool.Pool), typeName)
	}

	pattern := "type_xenum"
	verbose := true
	res, err := dbcommands.ListDatatypes(context.Background(), db, pattern, verbose)
	if err != nil {
		t.Fatalf("ListDatatypes failed: %v", err)
	}
	result := RequiresRowResult(t, res)
	defer result.Rows.Close()

	fds := result.Rows.FieldDescriptions()
	assert.NotNil(t, fds)

	expectedColumns := []string{
		"schema",
		"name",
		"internal_name",
		"size",
		"elements",
		"access_privileges",
		"description",
	}
	assert.Equal(t, expectedColumns, getColumnNames(fds), "Column names do not match expected")
	// expecting 7 columns
	assert.Len(t, fds, 7)
	var allRows []map[string]interface{}
	allRows, err = RowsToMaps(result.Rows)
	if err != nil {
		t.Fatalf("Failed to read rows: %v", err)
	}
	assert.Len(t, allRows, 0)
}
