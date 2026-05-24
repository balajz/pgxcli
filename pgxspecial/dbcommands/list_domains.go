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
		Cmd: "\\dD",
		Syntax: "\\dD[+] [pattern]",
		Description: "List or describe domains.",
		CaseSensitive: true,
		Handler: ListDomains,
	})
}

func ListDomains(ctx context.Context, db database.Queryer, pattern string, verbose bool) (pgxspecial.SpecialCommandResult, error) {
	var sb strings.Builder
	args := []any{}
	argIndex := 1

	sb.WriteString(`
	        SELECT n.nspname AS schema,
               t.typname AS name,
               pg_catalog.format_type(t.typbasetype, t.typtypmod) as type,
               pg_catalog.ltrim((COALESCE((SELECT (' collate ' || c.collname)
                                           FROM pg_catalog.pg_collation AS c,
                                                pg_catalog.pg_type AS bt
                                           WHERE c.oid = t.typcollation
                                             AND bt.oid = t.typbasetype
                                             AND t.typcollation <> bt.typcollation) , '')
                                || CASE
                                     WHEN t.typnotnull
                                       THEN ' not null'
                                     ELSE ''
                                   END) || CASE
                                             WHEN t.typdefault IS NOT NULL
                                               THEN(' default ' || t.typdefault)
                                             ELSE ''
                                           END) AS modifier,
               pg_catalog.array_to_string(ARRAY(
                 SELECT pg_catalog.pg_get_constraintdef(r.oid, TRUE)
                 FROM pg_catalog.pg_constraint AS r
				WHERE t.oid = r.contypid), ' ') AS check 
		`)

	if verbose {
		sb.WriteString(`,
		pg_catalog.array_to_string(t.typacl, E'\n') AS access_privileges,
               d.description as description
		`)
	}
	sb.WriteString(`
	        FROM pg_catalog.pg_type AS t
           LEFT JOIN pg_catalog.pg_namespace AS n ON n.oid = t.typnamespace`)

	if verbose {
		sb.WriteString(`
		LEFT JOIN pg_catalog.pg_description d ON d.classoid = t.tableoid
                                                AND d.objoid = t.oid AND d.objsubid = 0
			`)
	}

	sb.WriteString(` WHERE t.typtype = 'd' `)
	if pattern != "" {
		schemaRe, nameRe := sqlNamePattern(pattern)
		if schemaRe != "" || nameRe != "" {
			if schemaRe != "" {
				sb.WriteString(" AND n.nspname ~ $" + strconv.Itoa(argIndex) + "\n")
				args = append(args, schemaRe)
				argIndex++
			}
			if nameRe != "" {
				sb.WriteString(" AND t.typname ~ $" + strconv.Itoa(argIndex) + "\n")
				args = append(args, nameRe)
			}
		} else {
			sb.WriteString(`
			AND n.nspname <> 'pg_catalog'
			AND n.nspname <> 'information_schema'
			AND pg_catalog.pg_type_is_visible(t.oid)
			`)
		}
	}
	sb.WriteString("ORDER BY 1, 2;")
	rows, err := db.Query(ctx, sb.String(), args...)
	if err != nil {
		return nil, err
	}

	return pgxspecial.RowResult{Rows: rows}, nil
}
