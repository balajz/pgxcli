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
		Cmd:           "\\dT",
		Description:   "List data types",
		Syntax:        "\\dT[+] [pattern]",
		Handler:       ListDatatypes,
		CaseSensitive: true,
	})
}

func ListDatatypes(ctx context.Context, db database.Queryer, pattern string, verbose bool) (pgxspecial.SpecialCommandResult, error) {
	var sb strings.Builder
	args := []any{}
	argIndex := 1

	sb.WriteString(`
SELECT n.nspname AS schema,
       pg_catalog.format_type(t.oid, NULL) AS name,
`)

	if verbose {
		sb.WriteString(`
       t.typname AS internal_name,
       CASE
           WHEN t.typrelid != 0 THEN 'tuple'
           WHEN t.typlen < 0 THEN 'var'
           ELSE t.typlen::text
       END AS size,
       pg_catalog.array_to_string(
           ARRAY(
               SELECT e.enumlabel
               FROM pg_catalog.pg_enum e
               WHERE e.enumtypid = t.oid
               ORDER BY e.enumsortorder
           ), E'\n') AS elements,
       pg_catalog.array_to_string(t.typacl, E'\n') AS access_privileges,
       pg_catalog.obj_description(t.oid, 'pg_type') AS description
`)
	} else {
		sb.WriteString(`
       pg_catalog.obj_description(t.oid, 'pg_type') AS description
`)
	}

	sb.WriteString(`
FROM pg_catalog.pg_type t
LEFT JOIN pg_catalog.pg_namespace n ON n.oid = t.typnamespace
WHERE (t.typrelid = 0 OR
      (SELECT c.relkind='c' FROM pg_catalog.pg_class c WHERE c.oid = t.typrelid))
  AND NOT EXISTS(
      SELECT 1 FROM pg_catalog.pg_type el
      WHERE el.oid = t.typelem AND el.typarray = t.oid
  )
`)

	schemaPattern, typePattern := sqlNamePattern(pattern)

	if schemaPattern != "" {
		sb.WriteString("  AND n.nspname ~ $" + strconv.Itoa(argIndex) + "\n")
		args = append(args, schemaPattern)
		argIndex++
	} else {
		sb.WriteString("  AND pg_catalog.pg_type_is_visible(t.oid)\n")
	}

	if typePattern != "" {
		sb.WriteString("  AND (t.typname ~ $" + strconv.Itoa(argIndex) +
			" OR pg_catalog.format_type(t.oid, NULL) ~ $" + strconv.Itoa(argIndex) + ")\n")
		args = append(args, typePattern)
	}

	if schemaPattern == "" && typePattern == "" {
		sb.WriteString(`
  AND n.nspname <> 'pg_catalog'
  AND n.nspname <> 'information_schema'
`)
	}

	sb.WriteString("ORDER BY 1, 2;")

	rows, err := db.Query(ctx, sb.String(), args...)

	return pgxspecial.RowResult{Rows: rows}, err
}
