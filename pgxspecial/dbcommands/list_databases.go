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
		Cmd:           "\\l",
		Alias:         []string{"\\list"},
		Description:   "List Databases",
		Syntax:        "\\l[+] [pattern]",
		Handler:       ListDatabases,
		CaseSensitive: true,
	})
}

func ListDatabases(ctx context.Context, db database.Queryer, pattern string, verbose bool) (pgxspecial.SpecialCommandResult, error) {
	var sb strings.Builder
	args := []any{}
	argIndex := 1

	sb.WriteString(
		`SELECT d.datname as name,
        pg_catalog.pg_get_userbyid(d.datdba) as owner,
        pg_catalog.pg_encoding_to_char(d.encoding) as encoding,
        d.datcollate as collate,
        d.datctype as ctype,
        pg_catalog.array_to_string(d.datacl, E'\n') AS access_privileges
		`)

	if verbose {
		sb.WriteString(
			`, 
			CASE WHEN pg_catalog.has_database_privilege(d.datname, 'CONNECT')
				THEN pg_catalog.pg_size_pretty(pg_catalog.pg_database_size(d.datname))
				ELSE 'No Access'
            END as size,
            t.spcname as "Tablespace",
            pg_catalog.shobj_description(d.oid, 'pg_database') as description
	`)
	}

	sb.WriteString(`
	FROM pg_catalog.pg_database d
	`)

	if verbose {
		sb.WriteString(`JOIN pg_catalog.pg_tablespace t on d.dattablespace = t.oid`)
	}

	if pattern != "" {
		_, tablePattern := sqlNamePattern(pattern)

		if tablePattern != "" {
			args = append(args, tablePattern)
		}
		sb.WriteString("\nWHERE d.datname ~ $" + strconv.Itoa(argIndex) + " ")
	}

	sb.WriteString("\nORDER BY 1;")
	res, err := db.Query(ctx, sb.String(), args...)
	if err != nil {
		return nil, err
	}

	return pgxspecial.RowResult{Rows: res}, nil
}
