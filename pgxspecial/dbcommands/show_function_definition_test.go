package dbcommands_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/balajz/pgxcli/pgxspecial/dbcommands"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupFunction creates a function for testing and registers a cleanup to drop it.
func setupFunction(t *testing.T, ctx context.Context, db *pgxpool.Pool, name, args, body string) {
	t.Helper()

	createSQL := fmt.Sprintf(`
		CREATE OR REPLACE FUNCTION %s(%s)
		RETURNS integer AS $$
		BEGIN
			%s
		END;
		$$ LANGUAGE plpgsql;
	`, name, args, body)

	_, err := db.Exec(ctx, createSQL)
	require.NoError(t, err, "Failed to create function")
}

func teardownFunction(t *testing.T, ctx context.Context, db *pgxpool.Pool, name, args string) {
	dropSQL := fmt.Sprintf("DROP FUNCTION IF EXISTS %s(%s);", name, args)
	_, err := db.Exec(ctx, dropSQL)
	require.NoError(t, err, "Failed to drop function")
}

func TestShowFunctionDefinition(t *testing.T) {
	db := connectTestDB(t)
	defer db.(*pgxpool.Pool).Close()

	ctx := context.Background()
	setupFunction(t, ctx, db.(*pgxpool.Pool), "add_numbers", "a integer, b integer", "RETURN a + b;")
	setupFunction(t, ctx, db.(*pgxpool.Pool), "get_answer", "", "RETURN 42;")
	defer teardownFunction(t, ctx, db.(*pgxpool.Pool), "add_numbers", "integer, integer")
	defer teardownFunction(t, ctx, db.(*pgxpool.Pool), "get_answer", "")
	pattern := "add_numbers(integer, integer)"
	verbose := false

	res, err := dbcommands.ShowFunctionDefinition(ctx, db, pattern, verbose)
	if err != nil {
		t.Fatalf("ShowFunctionDefinition failed: %v", err)
	}
	result := RequiresRowResult(t, res)
	defer result.Rows.Close()

	var source string
	if result.Rows.Next() {
		err = result.Rows.Scan(&source)
		if err != nil {
			t.Fatalf("Failed to scan result: %v", err)
		}
	} else {
		t.Fatalf("No rows returned")
	}

	assert.Contains(t, strings.TrimSpace(source), "add_numbers", "Function definition does not match expected")
	assert.Contains(t, strings.TrimSpace(source), "RETURN a + b;", "Function body does not match expected")
	assert.False(t, result.Rows.Next(), "Expected only one row")
}

func TestShowFunctionDefinitionVerbose(t *testing.T) {
	db := connectTestDB(t).(*pgxpool.Pool)
	defer db.Close()

	ctx := context.Background()

	// Create a sample function
	_, err := db.Exec(ctx, `
	CREATE OR REPLACE FUNCTION add_numbers_verbose(a integer, b integer)
	RETURNS integer AS $$
	BEGIN
		RETURN a + b;
	END;
	$$ LANGUAGE plpgsql;
	`)
	if err != nil {
		t.Fatalf("Failed to create function: %v", err)
	}
	defer db.Exec(ctx, `DROP FUNCTION IF EXISTS add_numbers_verbose(integer, integer);`)

	pattern := "add_numbers_verbose(integer, integer)"
	verbose := true

	res, err := dbcommands.ShowFunctionDefinition(ctx, db, pattern, verbose)
	if err != nil {
		t.Fatalf("ShowFunctionDefinition with verbose failed: %v", err)
	}
	result := RequiresRowResult(t, res)
	defer result.Rows.Close()

	var source string
	if result.Rows.Next() {
		err = result.Rows.Scan(&source)
		if err != nil {
			t.Fatalf("Failed to scan result: %v", err)
		}
	} else {
		t.Fatalf("No rows returned")
	}

	if !strings.Contains(source, "1      ") {
		t.Errorf("Verbose output does not contain expected line numbers. Got: %s", source)
	}

	if !strings.Contains(source, "BEGIN") {
		t.Errorf("Verbose output does not contain function body. Got: %s", source)
	}
}
