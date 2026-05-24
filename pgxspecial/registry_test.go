package pgxspecial_test

import (
	"context"
	"os"
	"testing"

	"github.com/balajz/pgxcli/pgxspecial"
	"github.com/balajz/pgxcli/pgxspecial/database"
	_ "github.com/balajz/pgxcli/pgxspecial/dbcommands" // to register commands
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
)

func connectTestDB(t *testing.T) database.Queryer {
	t.Helper()
	ctx := context.Background()
	db_url := os.Getenv("PGXSPECIAL_TEST_DSN")
	db, err := pgx.Connect(ctx, db_url)

	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	return db
}

func RowsToMaps(rows pgx.Rows) ([]map[string]interface{}, error) {
	cols := rows.FieldDescriptions()
	colCount := len(cols)

	var result []map[string]interface{}

	for rows.Next() {
		values := make([]interface{}, colCount)
		scanArgs := make([]interface{}, colCount)
		for i := range values {
			scanArgs[i] = &values[i]
		}

		if err := rows.Scan(scanArgs...); err != nil {
			return nil, err
		}

		m := make(map[string]interface{})
		for i, fd := range cols {
			m[string(fd.Name)] = values[i]
		}

		result = append(result, m)
	}

	return result, rows.Err()
}

func containsByField(rows []map[string]interface{}, field, expected string) bool {
	for _, row := range rows {
		v := row[field]
		switch x := v.(type) {
		case string:
			if x == expected {
				return true
			}
		case []byte:
			if string(x) == expected {
				return true
			}
		}
	}
	return false
}

func getColumnNames(fds []pgconn.FieldDescription) []string {
	columns := make([]string, len(fds))
	for i, fd := range fds {
		columns[i] = string(fd.Name)
	}
	return columns
}
func isValidListDatabasesResult(t *testing.T, testingResult pgx.Rows) {
	t.Helper()

	fds := testingResult.FieldDescriptions()
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

	allRows, err := RowsToMaps(testingResult)
	if err != nil {
		t.Fatalf("Failed to read rows: %v", err)
	}
	assert.True(t, containsByField(allRows, "name", "template0"))
	assert.True(t, containsByField(allRows, "name", "template1"))
	assert.True(t, containsByField(allRows, "name", "postgres"))
}

func TestExecuteSpecialCommandWithUnknownCommand(t *testing.T) {
	db := connectTestDB(t)
	defer db.(*pgx.Conn).Close(t.Context())
	ctx := context.Background()

	// Example test for an unknown command
	_, isSpecial, err := pgxspecial.ExecuteSpecialCommand(ctx, db, "\\unknowncmd arg1 arg2")
	if err == nil || !isSpecial {
		t.Errorf("Expected error for unknown command, got nil")
	}
}

func TestExecuteSpecialCommandWithKnownCommand(t *testing.T) {
	db := connectTestDB(t)
	defer db.(*pgx.Conn).Close(t.Context())
	ctx := context.Background()

	// Register a test command
	pgxspecial.RegisterCommand(pgxspecial.SpecialCommandRegistry{
		Cmd:         "\\testcmd",
		Description: "A test command",
		Syntax:      "\\testcmd [args]",
		Handler: func(ctx context.Context, db database.Queryer, args string, verbose bool) (pgxspecial.SpecialCommandResult, error) {
			return nil, nil
		},
	})

	// Example test for the registered command
	_, isSpecial, err := pgxspecial.ExecuteSpecialCommand(ctx, db, "\\testcmd arg1 arg2")
	if err != nil {
		t.Errorf("Expected no error for known command, got: %v", err)
	}
	if !isSpecial {
		t.Errorf("Expected isSpecial to be true for known command")
	}
}

func TestExecuteCommand(t *testing.T) {
	db := connectTestDB(t)
	defer db.(*pgx.Conn).Close(t.Context())
	ctx := context.Background()

	// test for registered command
	result, isSpecial, err := pgxspecial.ExecuteSpecialCommand(ctx, db, "\\l")
	if err != nil {
		t.Errorf("Expected no error for known command, got: %v", err)
	}
	if !isSpecial {
		t.Errorf("Expected isSpecial to be true for known command")
	}

	if result.ResultKind() != pgxspecial.ResultKindRows {
		t.Errorf("Expected result kind to be rows, got: %v", result.ResultKind())
	}

	rows := result.(pgxspecial.RowResult).Rows

	defer rows.Close()

	isValidListDatabasesResult(t, rows)
}

func TestExecuteSpecialCommandNonSpecial(t *testing.T) {
	db := connectTestDB(t)
	defer db.(*pgx.Conn).Close(t.Context())
	ctx := context.Background()

	// Example test for a non-special command
	_, isSpecial, err := pgxspecial.ExecuteSpecialCommand(ctx, db, "SELECT * FROM users;")
	if err != nil {
		t.Errorf("Expected no error for non-special command, got: %v", err)
	}
	if isSpecial {
		t.Errorf("Expected isSpecial to be false for non-special command")
	}
}

func TestRegisterCommandAlias(t *testing.T) {
	db := connectTestDB(t)
	defer db.(*pgx.Conn).Close(t.Context())
	ctx := context.Background()

	result, isSpecial, err := pgxspecial.ExecuteSpecialCommand(ctx, db, "\\list")
	if err != nil {
		t.Errorf("Expected no error for known command alias, got: %v", err)
	}
	if !isSpecial {
		t.Errorf("Expected isSpecial to be true for known command alias")
	}
	if result.ResultKind() != pgxspecial.ResultKindRows {
		t.Errorf("Expected result kind to be rows, got: %v", result.ResultKind())
	}

	rows := result.(pgxspecial.RowResult).Rows
	defer rows.Close()

	isValidListDatabasesResult(t, rows)
}
