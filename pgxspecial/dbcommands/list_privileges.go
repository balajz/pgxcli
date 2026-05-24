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
		Cmd:           "\\dp",
		Alias:         []string{"\\z"},
		Description:   "List privileges.",
		Syntax:        "\\dp [pattern]",
		Handler:       ListPrivileges,
		CaseSensitive: true,
	})
}

func ListPrivileges(ctx context.Context, db database.Queryer, pattern string, verbose bool) (pgxspecial.SpecialCommandResult, error) {
	var sb strings.Builder
	args := []any{}
	argIndex := 1

	sb.WriteString(`
	        SELECT n.nspname as schema,
          c.relname as name,
          CASE c.relkind WHEN 'r' THEN 'table'
                         WHEN 'v' THEN 'view'
                         WHEN 'm' THEN 'materialized view'
                         WHEN 'S' THEN 'sequence'
                         WHEN 'f' THEN 'foreign table'
                         WHEN 'p' THEN 'partitioned table' END as type,
          pg_catalog.array_to_string(c.relacl, E'\n') AS access_privileges,

          pg_catalog.array_to_string(ARRAY(
            SELECT attname || E':\n  ' || pg_catalog.array_to_string(attacl, E'\n  ')
            FROM pg_catalog.pg_attribute a
            WHERE attrelid = c.oid AND NOT attisdropped AND attacl IS NOT NULL
          ), E'\n') AS column_privileges,
          pg_catalog.array_to_string(ARRAY(
            SELECT polname
            || CASE WHEN NOT polpermissive THEN
               E' (RESTRICTIVE)'
               ELSE '' END
            || CASE WHEN polcmd != '*' THEN
                   E' (' || polcmd::pg_catalog.text || E'):'
               ELSE E':'
               END
            || CASE WHEN polqual IS NOT NULL THEN
                   E'\n  (u): ' || pg_catalog.pg_get_expr(polqual, polrelid)
               ELSE E''
               END
            || CASE WHEN polwithcheck IS NOT NULL THEN
                   E'\n  (c): ' || pg_catalog.pg_get_expr(polwithcheck, polrelid)
               ELSE E''
               END    || CASE WHEN polroles <> '{0}' THEN
                   E'\n  to: ' || pg_catalog.array_to_string(
                       ARRAY(
                           SELECT rolname
                           FROM pg_catalog.pg_roles
                           WHERE oid = ANY (polroles)
                           ORDER BY 1
                       ), E', ')
               ELSE E''
               END
            FROM pg_catalog.pg_policy pol
            WHERE polrelid = c.oid), E'\n')
            AS policies
        FROM pg_catalog.pg_class c
             LEFT JOIN pg_catalog.pg_namespace n ON n.oid = c.relnamespace
		  WHERE c.relkind IN ('r','v','m','S','f','p')
	`)

	if pattern != "" {
		schema, table := sqlNamePattern(pattern)
		if table != "" {
			sb.WriteString(" AND c.relname OPERATOR(pg_catalog.~) $" + strconv.Itoa(argIndex) + " COLLATE pg_catalog.default ")
			args = append(args, table)
			argIndex++
		}
		if schema != "" {
			sb.WriteString(" AND n.nspname OPERATOR(pg_catalog.~) $" + strconv.Itoa(argIndex) + " COLLATE pg_catalog.default ")
			args = append(args, schema)
		}
	} else {
		sb.WriteString(" AND pg_catalog.pg_table_is_visible(c.oid) ")
	}

	sb.WriteString("  AND n.nspname !~ '^pg_'")
	sb.WriteString(" ORDER BY 1, 2")
	rows, err := db.Query(ctx, sb.String(), args...)
	return pgxspecial.RowResult{Rows: rows}, err
}
