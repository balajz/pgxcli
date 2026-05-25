package database

import (
	"context"
	"fmt"
	"strings"

	"github.com/balajz/pgxcli/internal/completer"
)

func (c *Client) Schemas() ([]string, error) {
	query := `
SELECT schema_name
FROM information_schema.schemata
WHERE schema_name NOT IN ('information_schema')
  AND schema_name NOT LIKE 'pg_toast%'
ORDER BY schema_name`

	rows, err := c.queryRows(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]string, 0)
	for rows.Next() {
		var schema string
		if err := rows.Scan(&schema); err != nil {
			return nil, fmt.Errorf("scanning schema metadata: %w", err)
		}
		out = append(out, schema)
	}

	return out, rows.Err()
}

func (c *Client) Tables() ([]completer.Relation, error) {
	query := `
SELECT table_schema, table_name
FROM information_schema.tables
WHERE table_type = 'BASE TABLE'
  AND table_schema NOT IN ('information_schema', 'pg_catalog')
ORDER BY table_schema, table_name`

	rows, err := c.queryRows(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]completer.Relation, 0)
	for rows.Next() {
		var schema, name string
		if err := rows.Scan(&schema, &name); err != nil {
			return nil, fmt.Errorf("scanning table metadata: %w", err)
		}
		out = append(out, completer.Relation{
			Schema: schema,
			Name:   name,
			Kind:   completer.RelationKindTable,
		})
	}

	return out, rows.Err()
}

func (c *Client) Views() ([]completer.Relation, error) {
	query := `
SELECT table_schema, table_name
FROM information_schema.views
WHERE table_schema NOT IN ('information_schema', 'pg_catalog')
ORDER BY table_schema, table_name`

	rows, err := c.queryRows(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]completer.Relation, 0)
	for rows.Next() {
		var schema, name string
		if err := rows.Scan(&schema, &name); err != nil {
			return nil, fmt.Errorf("scanning view metadata: %w", err)
		}
		out = append(out, completer.Relation{
			Schema: schema,
			Name:   name,
			Kind:   completer.RelationKindView,
		})
	}

	return out, rows.Err()
}

func (c *Client) TableColumns() ([]completer.ColumnInfo, error) {
	return c.columnsByTableType("BASE TABLE")
}

func (c *Client) ViewColumns() ([]completer.ColumnInfo, error) {
	return c.columnsByTableType("VIEW")
}

func (c *Client) columnsByTableType(tableType string) ([]completer.ColumnInfo, error) {
	query := `
SELECT c.table_schema, c.table_name, c.column_name, c.data_type, c.column_default IS NOT NULL, c.column_default
FROM information_schema.columns c
JOIN information_schema.tables t
  ON t.table_schema = c.table_schema
 AND t.table_name = c.table_name
WHERE t.table_type = $1
  AND c.table_schema NOT IN ('information_schema', 'pg_catalog')
ORDER BY c.table_schema, c.table_name, c.ordinal_position`

	rows, err := c.queryRows(query, tableType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]completer.ColumnInfo, 0)
	for rows.Next() {
		var schema, table, column, dataType string
		var hasDefault bool
		var defaultVal *string
		if err := rows.Scan(&schema, &table, &column, &dataType, &hasDefault, &defaultVal); err != nil {
			return nil, fmt.Errorf("scanning column metadata: %w", err)
		}
		out = append(out, completer.ColumnInfo{
			Schema:     schema,
			Table:      table,
			Column:     column,
			DataType:   dataType,
			HasDefault: hasDefault,
			Default:    defaultVal,
		})
	}

	return out, rows.Err()
}

func (c *Client) Functions() ([]*completer.FunctionMetadata, error) {
	return []*completer.FunctionMetadata{}, nil
}

func (c *Client) DataTypes() ([]completer.DatatypeName, error) {
	return []completer.DatatypeName{}, nil
}

func (c *Client) ForeignKeys() ([]completer.ForeignKey, error) {
	query := `
SELECT
    tc.table_schema AS child_schema,
    tc.table_name AS child_table,
    kcu.column_name AS child_column,
    ccu.table_schema AS parent_schema,
    ccu.table_name AS parent_table,
    ccu.column_name AS parent_column
FROM information_schema.table_constraints tc
JOIN information_schema.key_column_usage kcu
  ON tc.constraint_name = kcu.constraint_name
 AND tc.table_schema = kcu.table_schema
JOIN information_schema.constraint_column_usage ccu
  ON ccu.constraint_name = tc.constraint_name
 AND ccu.constraint_schema = tc.table_schema
WHERE tc.constraint_type = 'FOREIGN KEY'
ORDER BY child_schema, child_table, child_column`

	rows, err := c.queryRows(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]completer.ForeignKey, 0)
	for rows.Next() {
		var fk completer.ForeignKey
		if err := rows.Scan(
			&fk.ChildSchema,
			&fk.ChildTable,
			&fk.ChildColumn,
			&fk.ParentSchema,
			&fk.ParentTable,
			&fk.ParentColumn,
		); err != nil {
			return nil, fmt.Errorf("scanning foreign key metadata: %w", err)
		}
		out = append(out, fk)
	}

	return out, rows.Err()
}

func (c *Client) Databases() ([]string, error) {
	query := `
SELECT datname
FROM pg_database
WHERE datistemplate = false
ORDER BY datname`

	rows, err := c.queryRows(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]string, 0)
	for rows.Next() {
		var database string
		if err := rows.Scan(&database); err != nil {
			return nil, fmt.Errorf("scanning database metadata: %w", err)
		}
		out = append(out, database)
	}

	return out, rows.Err()
}

func (c *Client) SearchPath() ([]string, error) {
	rows, err := c.queryRows("SHOW search_path")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var raw string
	if rows.Next() {
		if err := rows.Scan(&raw); err != nil {
			return nil, fmt.Errorf("scanning search_path metadata: %w", err)
		}
	}

	path := strings.Split(raw, ",")
	out := make([]string, 0, len(path))
	for _, p := range path {
		p = strings.TrimSpace(strings.Trim(p, `"`))
		if p == "" {
			continue
		}
		out = append(out, p)
	}

	return out, rows.Err()
}

func (c *Client) queryRows(query string, args ...any) (rows pgxRows, err error) {
	if !c.IsConnected() {
		return nil, ErrConnectionNotEstablished
	}
	rows, err = c.executor.conn.Query(context.Background(), query, args...)
	if err != nil {
		return nil, fmt.Errorf("querying metadata: %w", err)
	}
	return rows, nil
}

type pgxRows interface {
	Next() bool
	Scan(dest ...any) error
	Close()
	Err() error
}
