package dbcommands

import (
	"context"
	"strconv"
	"strings"

	"github.com/balajz/pgxcli/pgxspecial"
	"github.com/balajz/pgxcli/pgxspecial/database"
)

func init() {
	pgxspecial.RegisterCommand(pgxspecial.SpecialCommandRegistry{
		Cmd:           "\\dE",
		Description:   "List foreign tables.",
		Syntax:        "\\dE[+] [pattern]",
		Handler:       ListForeignTables,
		CaseSensitive: true,
	})
}

func ListForeignTables(ctx context.Context, db database.Queryer, pattern string, verbose bool) (pgxspecial.SpecialCommandResult, error) {
	var sb strings.Builder
	args := []any{}
	argIndex := 1

	sb.WriteString(`
SELECT 
    n.nspname AS schema,
    c.relname AS name,
    CASE c.relkind 
        WHEN 'r' THEN 'table'
        WHEN 'v' THEN 'view'
        WHEN 'm' THEN 'materialized view'
        WHEN 'i' THEN 'index'
        WHEN 'S' THEN 'sequence'
        WHEN 's' THEN 'special'
        WHEN 'f' THEN 'foreign table'
        WHEN 'p' THEN 'table'
        WHEN 'I' THEN 'index'
    END AS type,
    pg_catalog.pg_get_userbyid(c.relowner) AS owner
`)

	if verbose {
		sb.WriteString(`
  , pg_catalog.pg_size_pretty(pg_catalog.pg_table_size(c.oid)) AS size
  , pg_catalog.obj_description(c.oid, 'pg_class') AS description
`)
	}

	sb.WriteString(`
FROM pg_catalog.pg_class c
LEFT JOIN pg_catalog.pg_namespace n ON n.oid = c.relnamespace
WHERE c.relkind IN ('f','')
  AND n.nspname <> 'pg_catalog'
  AND n.nspname <> 'information_schema'
  AND n.nspname !~ '^pg_toast'
  AND pg_catalog.pg_table_is_visible(c.oid)
`)

	if pattern != "" {
		_, tblPattern := sqlNamePattern(pattern)
		sb.WriteString("  AND c.relname OPERATOR(pg_catalog.~) $" + strconv.Itoa(argIndex) + "\n")
		args = append(args, tblPattern)
	}

	sb.WriteString("ORDER BY 1,2;")

	rows, err := db.Query(ctx, sb.String(), args...)
	return pgxspecial.RowResult{Rows: rows}, err
}
