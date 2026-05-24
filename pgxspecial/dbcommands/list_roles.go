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
		Cmd:           "\\du",
		Description:   "List roles.",
		Syntax:        "\\du[+] [pattern]",
		Handler:       ListRoles,
		CaseSensitive: true,
	})
}

func ListRoles(ctx context.Context, db database.Queryer, pattern string, verbose bool) (pgxspecial.SpecialCommandResult, error) {
	var sb strings.Builder
	args := []any{}
	argIndex := 1

	sb.WriteString(`
	   SELECT r.rolname,
                r.rolsuper,
                r.rolinherit,
                r.rolcreaterole,
                r.rolcreatedb,
                r.rolcanlogin,
                r.rolconnlimit,
                r.rolvaliduntil,
                ARRAY(SELECT b.rolname FROM pg_catalog.pg_auth_members m JOIN pg_catalog.pg_roles b ON (m.roleid = b.oid) WHERE m.member = r.oid) as memberof,
	`)

	if verbose {
		sb.WriteString("pg_catalog.shobj_description(r.oid, 'pg_authid') AS description, ")
	}
	sb.WriteString(`
	 	r.rolreplication
			FROM pg_catalog.pg_roles r
	`)

	if pattern != "" {
		_, tablePattern := sqlNamePattern(pattern)
		if tablePattern != "" {
			sb.WriteString(" WHERE r.rolname ~ $" + strconv.Itoa(argIndex) + " ")
			args = append(args, tablePattern)
		}
	}

	sb.WriteString(" ORDER BY 1;")
	rows, err := db.Query(ctx, sb.String(), args...)
	return pgxspecial.RowResult{Rows: rows}, err
}
