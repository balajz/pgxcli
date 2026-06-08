//go:build integration

package dbcommands_test

// this file contains utility functions for setting up and tearing down database objects
import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/balajz/pgxcli/pgxspecial"
	"github.com/balajz/pgxcli/pgxspecial/database"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

func connectTestDB(t *testing.T) database.Queryer {
	t.Helper()
	ctx := context.Background()
	db_url := os.Getenv("PGXSPECIAL_TEST_DSN")
	db, err := pgxpool.New(ctx, db_url)
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	return db
}

func CreateForeignTable(t *testing.T, ctx context.Context, pool *pgxpool.Pool, tableName string) {
	t.Helper()

	// Create extension
	_, err := pool.Exec(ctx, `CREATE EXTENSION IF NOT EXISTS postgres_fdw;`)
	if err != nil {
		t.Fatalf("failed to create extension: %v", err)
	}

	// Create server
	_, err = pool.Exec(ctx, `
        CREATE SERVER IF NOT EXISTS test_remote_server
        FOREIGN DATA WRAPPER postgres_fdw
        OPTIONS (host 'localhost', dbname 'remotedb', port '5432');
    `)
	if err != nil {
		t.Fatalf("failed to create server: %v", err)
	}

	// Create user mapping
	_, err = pool.Exec(ctx, `
        CREATE USER MAPPING IF NOT EXISTS FOR CURRENT_USER
        SERVER test_remote_server
        OPTIONS (user 'remote_user', password 'remote_pass');
    `)
	if err != nil {
		t.Fatalf("failed to create user mapping: %v", err)
	}

	// Create FOREIGN TABLE
	query := fmt.Sprintf(`
        CREATE FOREIGN TABLE IF NOT EXISTS %s (
            id    integer,
            name  text,
            email text
        )
        SERVER test_remote_server
        OPTIONS (schema_name 'public', table_name 'users');
    `, tableName)

	_, err = pool.Exec(ctx, query)
	if err != nil {
		t.Fatalf("failed to create foreign table %s: %v", tableName, err)
	}
}

func DropForeignTable(t *testing.T, ctx context.Context, pool *pgxpool.Pool, tableName string) {
	t.Helper()

	query := fmt.Sprintf(`
        DROP FOREIGN TABLE IF EXISTS %s CASCADE;
    `, tableName)

	if _, err := pool.Exec(ctx, query); err != nil {
		t.Fatalf("failed to drop foreign table %s: %v", tableName, err)
	}
}

func CreateDatatype(t *testing.T, ctx context.Context, pool *pgxpool.Pool, typeName string) {
	t.Helper()

	// create an ENUM datatype
	query := fmt.Sprintf(`
        DO $$
        BEGIN
            IF NOT EXISTS (
                SELECT 1 FROM pg_type WHERE typname = '%s'
            ) THEN
                CREATE TYPE %s AS ENUM ('a', 'b', 'c');
            END IF;
        END$$;
    `, typeName, typeName)

	if _, err := pool.Exec(ctx, query); err != nil {
		t.Fatalf("failed to create datatype %s: %v", typeName, err)
	}
}

func DropDatatype(t *testing.T, ctx context.Context, pool *pgxpool.Pool, typeName string) {
	t.Helper()

	query := fmt.Sprintf(`
        DO $$
        BEGIN
            IF EXISTS (
                SELECT 1 FROM pg_type WHERE typname = '%s'
            ) THEN
                DROP TYPE %s;
            END IF;
        END$$;
    `, typeName, typeName)

	if _, err := pool.Exec(ctx, query); err != nil {
		t.Fatalf("failed to drop datatype %s: %v", typeName, err)
	}
}

func CreateFunction(t *testing.T, ctx context.Context, pool *pgxpool.Pool, funcName string) {
	t.Helper()

	// Simple example function: returns integer 42
	query := fmt.Sprintf(`
        CREATE OR REPLACE FUNCTION %s()
        RETURNS int
        LANGUAGE plpgsql
        AS $$
        BEGIN
            RETURN 42;
        END;
        $$;
    `, funcName)

	if _, err := pool.Exec(ctx, query); err != nil {
		t.Fatalf("failed to create function %s: %v", funcName, err)
	}
}

func DropFunction(t *testing.T, ctx context.Context, pool *pgxpool.Pool, funcName string) {
	t.Helper()

	query := fmt.Sprintf(`
        DROP FUNCTION IF EXISTS %s() CASCADE;
    `, funcName)

	if _, err := pool.Exec(ctx, query); err != nil {
		t.Fatalf("failed to drop function %s: %v", funcName, err)
	}
}

func CreateDefaultPrivileges(
	t *testing.T,
	ctx context.Context,
	pool *pgxpool.Pool,
	role string,
) {
	t.Helper()

	sql := `
		ALTER DEFAULT PRIVILEGES
		FOR ROLE current_user
		IN SCHEMA public
		GRANT SELECT ON TABLES TO ` + role + `;
	`

	if _, err := pool.Exec(ctx, sql); err != nil {
		t.Fatalf("create default privileges failed: %v", err)
	}
}

func DropDefaultPrivileges(
	t *testing.T,
	ctx context.Context,
	pool *pgxpool.Pool,
	role string,
) {
	t.Helper()

	sql := `
		ALTER DEFAULT PRIVILEGES
		FOR ROLE current_user
		IN SCHEMA public
		REVOKE SELECT ON TABLES FROM ` + role + `;
	`

	if _, err := pool.Exec(ctx, sql); err != nil {
		t.Fatalf("drop default privileges failed: %v", err)
	}
}

func CreateSchema(
	t *testing.T,
	ctx context.Context,
	pool *pgxpool.Pool,
	schema string,
) {
	t.Helper()

	sql := `CREATE SCHEMA ` + schema

	if _, err := pool.Exec(ctx, sql); err != nil {
		t.Fatalf("create schema %q failed: %v", schema, err)
	}
}

func DropSchema(
	t *testing.T,
	ctx context.Context,
	pool *pgxpool.Pool,
	schema string,
) {
	t.Helper()

	sql := `DROP SCHEMA ` + schema + ` CASCADE`

	if _, err := pool.Exec(ctx, sql); err != nil {
		t.Fatalf("drop schema %q failed: %v", schema, err)
	}
}

func GrantPrivilege(
	t *testing.T,
	ctx context.Context,
	pool *pgxpool.Pool,
	privilege string,
	object string,
	role string,
) {
	t.Helper()

	sql := `GRANT ` + privilege + ` ON ` + object + ` TO ` + role

	if _, err := pool.Exec(ctx, sql); err != nil {
		t.Fatalf("grant privilege failed: %v", err)
	}
}

func RevokePrivilege(
	t *testing.T,
	ctx context.Context,
	pool *pgxpool.Pool,
	privilege string,
	object string,
	role string,
) {
	t.Helper()

	sql := `REVOKE ` + privilege + ` ON ` + object + ` FROM ` + role

	if _, err := pool.Exec(ctx, sql); err != nil {
		t.Fatalf("revoke privilege failed: %v", err)
	}
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

func getColumnNames(fds []pgconn.FieldDescription) []string {
	columns := make([]string, len(fds))
	for i, fd := range fds {
		columns[i] = string(fd.Name)
	}
	return columns
}

func containsDB(rows []map[string]interface{}, name string) bool {
	for _, r := range rows {
		if n, ok := r["name"].(string); ok && n == name {
			return true
		}
	}
	return false
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

func RequiresRowResult(t *testing.T, r pgxspecial.SpecialCommandResult) pgxspecial.RowResult {
	t.Helper()

	_, isRow := r.(pgxspecial.RowResult)
	if !isRow {
		t.Fatalf("expected rows result, got %T", r)
	}

	rowsResult, ok := r.(pgxspecial.RowResult)
	if !ok {
		t.Fatalf("expected RowsResult, got %T", r)
	}

	return rowsResult
}
