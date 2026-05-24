package dbcommands_test

import (
	"context"
	"testing"

	"github.com/balajz/pgxcli/pgxspecial/dbcommands"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
)

func TestListFunctions(t *testing.T) {
	db := connectTestDB(t)
	defer db.(*pgxpool.Pool).Close()

	ctx := context.Background()
	funcNames := []string{"calculate_discount", "get_user_status", "compute_tax"}

	for _, funcName := range funcNames {
		CreateFunction(t, ctx, db.(*pgxpool.Pool), funcName)
		defer DropFunction(t, ctx, db.(*pgxpool.Pool), funcName)
	}

	pattern := ""
	verbose := false

	res, err := dbcommands.ListFunctions(context.Background(), db, pattern, verbose)
	if err != nil {
		t.Fatalf("ListFunctions failed: %v", err)
	}
	result := RequiresRowResult(t, res)

	defer result.Rows.Close()

	fds := result.Rows.FieldDescriptions()
	assert.NotNil(t, fds)

	columnsExpected := []string{
		"schema",
		"name",
		"Result data type",
		"Argument data types",
		"type",
	}
	assert.Equal(t, columnsExpected, getColumnNames(fds), "Column names do not match expected")
	// expecting 5 columns
	assert.Len(t, fds, 5)

	var allRows []map[string]interface{}
	allRows, err = RowsToMaps(result.Rows)
	if err != nil {
		t.Fatalf("Failed to read rows: %v", err)
	}
	for _, funcName := range funcNames {
		assert.True(t, containsByField(allRows, "name", funcName))
	}
}

func TestListFunctionsWithPattern(t *testing.T) {
	db := connectTestDB(t)
	defer db.(*pgxpool.Pool).Close()

	ctx := context.Background()
	funcNames := []string{"calculate_discount", "get_user_status", "compute_tax"}
	for _, funcName := range funcNames {
		CreateFunction(t, ctx, db.(*pgxpool.Pool), funcName)
		defer DropFunction(t, ctx, db.(*pgxpool.Pool), funcName)
	}

	pattern := "get_*"
	verbose := false

	res, err := dbcommands.ListFunctions(context.Background(), db, pattern, verbose)
	if err != nil {
		t.Fatalf("ListFunctions failed: %v", err)
	}
	result := RequiresRowResult(t, res)

	defer result.Rows.Close()

	fds := result.Rows.FieldDescriptions()
	assert.NotNil(t, fds)

	columnsExpected := []string{
		"schema",
		"name",
		"Result data type",
		"Argument data types",
		"type",
	}
	assert.Equal(t, columnsExpected, getColumnNames(fds), "Column names do not match expected")
	// expecting 5 columns
	assert.Len(t, fds, 5)

	var allRows []map[string]interface{}
	allRows, err = RowsToMaps(result.Rows)
	if err != nil {
		t.Fatalf("Failed to read rows: %v", err)
	}
	assert.True(t, containsByField(allRows, "name", "get_user_status"))
	assert.False(t, containsByField(allRows, "name", "calculate_discount"))
	assert.False(t, containsByField(allRows, "name", "compute_tax"))
}

func TestListFunctionsWithNoMatchingPattern(t *testing.T) {
	db := connectTestDB(t)
	defer db.(*pgxpool.Pool).Close()

	ctx := context.Background()
	funcNames := []string{"calculate_discount", "get_user_status", "compute_tax"}

	for _, funcName := range funcNames {
		CreateFunction(t, ctx, db.(*pgxpool.Pool), funcName)
		defer DropFunction(t, ctx, db.(*pgxpool.Pool), funcName)
	}

	pattern := "fetch_*"
	verbose := false

	res, err := dbcommands.ListFunctions(context.Background(), db, pattern, verbose)
	if err != nil {
		t.Fatalf("ListFunctions failed: %v", err)
	}
	result := RequiresRowResult(t, res)

	defer result.Rows.Close()

	fds := result.Rows.FieldDescriptions()
	assert.NotNil(t, fds)

	columnsExpected := []string{
		"schema",
		"name",
		"Result data type",
		"Argument data types",
		"type",
	}
	assert.Equal(t, columnsExpected, getColumnNames(fds), "Column names do not match expected")
	// expecting 5 columns
	assert.Len(t, fds, 5)

	var allRows []map[string]interface{}
	allRows, err = RowsToMaps(result.Rows)
	if err != nil {
		t.Fatalf("Failed to read rows: %v", err)
	}
	assert.Len(t, allRows, 0)
}

func TestListFunctionsVerbose(t *testing.T) {
	db := connectTestDB(t)
	defer db.(*pgxpool.Pool).Close()

	ctx := context.Background()
	funcNames := []string{"calculate_discount", "get_user_status", "compute_tax"}

	for _, funcName := range funcNames {
		CreateFunction(t, ctx, db.(*pgxpool.Pool), funcName)
		defer DropFunction(t, ctx, db.(*pgxpool.Pool), funcName)
	}

	pattern := ""
	verbose := true

	res, err := dbcommands.ListFunctions(context.Background(), db, pattern, verbose)
	if err != nil {
		t.Fatalf("ListFunctions failed: %v", err)
	}
	result := RequiresRowResult(t, res)

	defer result.Rows.Close()

	fds := result.Rows.FieldDescriptions()
	assert.NotNil(t, fds)

	columnsExpected := []string{
		"schema",
		"name",
		"Result data type",
		"Argument data types",
		"type",
		"Volatility",
		"owner",
		"Language",
		"Source code",
		"description",
	}
	assert.Equal(t, columnsExpected, getColumnNames(fds), "Column names do not match expected")
	// expecting 10 columns
	assert.Len(t, fds, 10)

	var allRows []map[string]interface{}
	allRows, err = RowsToMaps(result.Rows)
	if err != nil {
		t.Fatalf("Failed to read rows: %v", err)
	}
	for _, funcName := range funcNames {
		assert.True(t, containsByField(allRows, "name", funcName))
	}
}

func TestListFunctionsVerboseWithPattern(t *testing.T) {
	db := connectTestDB(t)
	defer db.(*pgxpool.Pool).Close()

	ctx := context.Background()
	funcNames := []string{"calculate_discount", "get_user_status", "compute_tax"}

	for _, funcName := range funcNames {
		CreateFunction(t, ctx, db.(*pgxpool.Pool), funcName)
		defer DropFunction(t, ctx, db.(*pgxpool.Pool), funcName)
	}

	pattern := "get_*"
	verbose := true

	res, err := dbcommands.ListFunctions(context.Background(), db, pattern, verbose)
	if err != nil {
		t.Fatalf("ListFunctions failed: %v", err)
	}
	result := RequiresRowResult(t, res)

	defer result.Rows.Close()

	fds := result.Rows.FieldDescriptions()
	assert.NotNil(t, fds)

	columnsExpected := []string{
		"schema",
		"name",
		"Result data type",
		"Argument data types",
		"type",
		"Volatility",
		"owner",
		"Language",
		"Source code",
		"description",
	}
	assert.Equal(t, columnsExpected, getColumnNames(fds), "Column names do not match expected")
	// expecting 10 columns
	assert.Len(t, fds, 10)

	var allRows []map[string]interface{}
	allRows, err = RowsToMaps(result.Rows)
	if err != nil {
		t.Fatalf("Failed to read rows: %v", err)
	}
	assert.True(t, containsByField(allRows, "name", "get_user_status"))
	assert.False(t, containsByField(allRows, "name", "calculate_discount"))
	assert.False(t, containsByField(allRows, "name", "compute_tax"))
}

func TestListFunctionsVerboseWithNoMatchingPattern(t *testing.T) {
	db := connectTestDB(t)
	defer db.(*pgxpool.Pool).Close()

	ctx := context.Background()
	funcNames := []string{"calculate_discount", "get_user_status", "compute_tax"}

	for _, funcName := range funcNames {
		CreateFunction(t, ctx, db.(*pgxpool.Pool), funcName)
		defer DropFunction(t, ctx, db.(*pgxpool.Pool), funcName)
	}

	pattern := "fetch_*"
	verbose := true

	res, err := dbcommands.ListFunctions(context.Background(), db, pattern, verbose)
	if err != nil {
		t.Fatalf("ListFunctions failed: %v", err)
	}
	result := RequiresRowResult(t, res)

	defer result.Rows.Close()

	fds := result.Rows.FieldDescriptions()
	assert.NotNil(t, fds)

	columnsExpected := []string{
		"schema",
		"name",
		"Result data type",
		"Argument data types",
		"type",
		"Volatility",
		"owner",
		"Language",
		"Source code",
		"description",
	}
	assert.Equal(t, columnsExpected, getColumnNames(fds), "Column names do not match expected")
	// expecting 10 columns
	assert.Len(t, fds, 10)

	var allRows []map[string]interface{}
	allRows, err = RowsToMaps(result.Rows)
	if err != nil {
		t.Fatalf("Failed to read rows: %v", err)
	}
	assert.Len(t, allRows, 0)
}
