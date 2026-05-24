package dbcommands_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/balajz/pgxcli/pgxspecial"
	"github.com/balajz/pgxcli/pgxspecial/dbcommands"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
)

func TestDescribeOneTableDetails(t *testing.T) {
	db := connectTestDB(t)
	defer db.(*pgxpool.Pool).Close()

	tables := []struct {
		name    string
		columns map[string]string
	}{
		{
			name: "test_table_1",
			columns: map[string]string{
				"id":   "SERIAL PRIMARY KEY",
				"name": "VARCHAR(100)",
			},
		},
		{
			name: "test_table_2",
			columns: map[string]string{
				"id":      "SERIAL PRIMARY KEY",
				"age":     "INT",
				"address": "TEXT",
			},
		},
	}

	// pattern := ""
	verbose := false

	for _, table := range tables {
		oid := CreateTable(t, context.Background(), db.(*pgxpool.Pool), table.name, table.columns)
		defer DropTable(t, context.Background(), db.(*pgxpool.Pool), table.name)
		result, err := dbcommands.DescribeOneTableDetails(context.Background(), db, "public", table.name, oid, verbose)
		if err != nil {
			t.Fatalf("DescribeTables failed: %v", err)
		}

		columnsExpected := []string{
			"Column",
			"Type",
			"Modifiers",
		}
		assert.Equal(t, columnsExpected, result.Columns, "Column names do not match expected")
		// expecting 3 columns
		assert.Len(t, result.Columns, 3)

		// Check for columns from both tables
		for col_name := range table.columns {
			found := false
			for _, row := range result.Data {
				if len(row) > 0 && row[0] == col_name {
					found = true
					break
				}
			}
			assert.True(t, found, fmt.Sprintf("Expected column %s not found", col_name))
		}
	}
}

func TestDescribeOneTableDetailsVerbose(t *testing.T) {
	db := connectTestDB(t)
	defer db.(*pgxpool.Pool).Close()

	tables := []struct {
		name    string
		columns map[string]string
	}{
		{
			name: "test_table_1",
			columns: map[string]string{
				"id":         "SERIAL PRIMARY KEY",
				"name":       "VARCHAR(100)",
				"age":        "INT",
				"email":      "VARCHAR(100)",
				"created_at": "TIMESTAMP DEFAULT NOW()",
			},
		},
		{
			name: "test_table_2",
			columns: map[string]string{
				"id":      "SERIAL PRIMARY KEY",
				"age":     "INT",
				"address": "TEXT",
			},
		},
	}

	// pattern := ""
	verbose := true

	for _, table := range tables {
		oid := CreateTable(t, context.Background(), db.(*pgxpool.Pool), table.name, table.columns)
		defer DropTable(t, context.Background(), db.(*pgxpool.Pool), table.name)
		result, err := dbcommands.DescribeOneTableDetails(context.Background(), db, "public", table.name, oid, verbose)
		if err != nil {
			t.Fatalf("DescribeTables failed: %v", err)
		}

		columnsExpected := []string{
			"Column",
			"Type",
			"Modifiers",
			"Storage",
			"Stats target",
			"Description",
		}
		assert.Equal(t, columnsExpected, result.Columns, "Column names do not match expected")
		// expecting 6 columns
		assert.Len(t, result.Columns, 6)

		// Check for columns from both tables
		for col_name := range table.columns {
			found := false
			for _, row := range result.Data {
				if len(row) > 0 && row[0] == col_name {
					found = true
					break
				}
			}
			assert.True(t, found, fmt.Sprintf("Expected column %s not found", col_name))
		}
	}
}

func TestDescribeTableMetadata(t *testing.T) {
	db := connectTestDB(t)
	defer db.(*pgxpool.Pool).Close()

	ctx := context.Background()
	tableName := "test_metadata_table"

	// Create a complex table with constraints, indexes, etc.
	setupSQL := `
		CREATE TABLE test_metadata_table (
			id SERIAL PRIMARY KEY,
			email VARCHAR(100) UNIQUE NOT NULL,
			age INT CHECK (age >= 18),
			parent_id INT REFERENCES test_metadata_table(id)
		);
		CREATE INDEX idx_metadata_age ON test_metadata_table(age);
		COMMENT ON COLUMN test_metadata_table.email IS 'User email address';
	`

	_, err := db.(*pgxpool.Pool).Exec(ctx, setupSQL)
	if err != nil {
		t.Fatalf("Failed to setup table: %v", err)
	}
	defer func() {
		_, _ = db.(*pgxpool.Pool).Exec(ctx, "DROP TABLE IF EXISTS test_metadata_table CASCADE")
	}()

	// Get OID
	var oid uint32
	err = db.(*pgxpool.Pool).QueryRow(ctx, "SELECT oid FROM pg_class WHERE relname = $1", tableName).Scan(&oid)
	if err != nil {
		t.Fatalf("Failed to get OID: %v", err)
	}

	// Test
	result, err := dbcommands.DescribeOneTableDetails(ctx, db, "public", tableName, oid, true)
	if err != nil {
		t.Fatalf("DescribeOneTableDetails failed: %v", err)
	}

	// Verify Metadata
	meta := result.TableMetaData

	// Indexes
	assert.NotEmpty(t, meta.Indexes, "Should have indexes")
	foundIndex := false
	for _, idx := range meta.Indexes {
		if strings.Contains(idx, "idx_metadata_age") {
			foundIndex = true
			break
		}
	}
	assert.True(t, foundIndex, "Expected index idx_metadata_age not found")

	// Check Constraints
	assert.NotEmpty(t, meta.CheckConstraints, "Should have check constraints")
	foundCheck := false
	for _, chk := range meta.CheckConstraints {
		if strings.Contains(chk, "age >= 18") {
			foundCheck = true
			break
		}
	}
	assert.True(t, foundCheck, "Expected check constraint on age not found")

	// Foreign Keys
	assert.NotEmpty(t, meta.ForeignKeys, "Should have foreign keys")

	// Referenced By
	assert.NotEmpty(t, meta.ReferencedBy, "Should have referenced by (self-reference)")

	// Column Description (Verbose)
	foundDesc := false
	for _, row := range result.Data {
		if row[0] == "email" {
			// Description is the last column in verbose mode
			if len(row) > 0 && row[len(row)-1] == "User email address" {
				foundDesc = true
			}
		}
	}
	assert.True(t, foundDesc, "Expected column comment not found")
}

func TestDescribeTableDetails_Patterns(t *testing.T) {
	db := connectTestDB(t)
	defer db.(*pgxpool.Pool).Close()
	ctx := context.Background()

	// Setup tables
	tables := []string{"pattern_test_1", "pattern_test_2", "other_table"}
	for _, name := range tables {
		CreateTable(t, ctx, db.(*pgxpool.Pool), name, map[string]string{"id": "int"})
		defer DropTable(t, ctx, db.(*pgxpool.Pool), name)
	}

	t.Run("Multiple Matches", func(t *testing.T) {
		res, err := dbcommands.DescribeTableDetails(ctx, db, "pattern_test*", false)
		assert.NoError(t, err)

		descRes, ok := res.(pgxspecial.DescribeTableListResult)
		assert.True(t, ok, "Expected DescribeTableListResult")
		assert.Len(t, descRes.Results, 2, "Expected 2 tables matching pattern")
	})

	t.Run("No Pattern", func(t *testing.T) {
		res, err := dbcommands.DescribeTableDetails(ctx, db, "", false)
		assert.NoError(t, err)

		rowRes, ok := res.(pgxspecial.RowResult)
		assert.True(t, ok, "Expected RowResult for empty pattern")
		if ok {
			rowRes.Rows.Close()
		}
	})

	t.Run("Verbose", func(t *testing.T) {
		res, err := dbcommands.DescribeTableDetails(ctx, db, "other_table", true)
		assert.NoError(t, err)

		descRes, ok := res.(pgxspecial.DescribeTableListResult)
		assert.True(t, ok)
		assert.Len(t, descRes.Results, 1)

		// Check for verbose columns
		assert.Contains(t, descRes.Results[0].Columns, "Storage")
		assert.Contains(t, descRes.Results[0].Columns, "Description")
	})
}

func CreateTable(
	t *testing.T,
	ctx context.Context,
	pool *pgxpool.Pool,
	tableName string,
	columns map[string]string,
) uint32 {
	t.Helper()

	if len(columns) == 0 {
		t.Fatal("columns map cannot be empty")
	}

	defs := make([]string, 0, len(columns))
	for col, typ := range columns {
		defs = append(defs, fmt.Sprintf("%s %s", col, typ))
	}

	createSQL := fmt.Sprintf(
		"CREATE TABLE %s (%s)",
		pgx.Identifier{tableName}.Sanitize(),
		strings.Join(defs, ", "),
	)

	if _, err := pool.Exec(ctx, createSQL); err != nil {
		t.Fatalf("failed to create table %s: %v", tableName, err)
	}

	var oid uint32
	err := pool.QueryRow(
		ctx,
		`SELECT oid FROM pg_class WHERE relname = $1 AND relkind = 'r'`,
		tableName,
	).Scan(&oid)

	if err != nil {
		t.Fatalf("failed to fetch OID for table %s: %v", tableName, err)
	}

	return oid
}

func DropTable(
	t *testing.T,
	ctx context.Context,
	pool *pgxpool.Pool,
	tableName string,
) {
	t.Helper()

	dropSQL := fmt.Sprintf(
		"DROP TABLE IF EXISTS %s",
		pgx.Identifier{tableName}.Sanitize(),
	)

	if _, err := pool.Exec(ctx, dropSQL); err != nil {
		t.Fatalf("failed to drop table %s: %v", tableName, err)
	}
}
