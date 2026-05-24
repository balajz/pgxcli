package dbcommands

import (
	"context"
	"fmt"
	"strings"

	"github.com/balajz/pgxcli/pgxspecial"
	"github.com/balajz/pgxcli/pgxspecial/database"
)

func init() {
	pgxspecial.RegisterCommand(pgxspecial.SpecialCommandRegistry{
		Cmd:           "\\ddp",
		Description:   "Lists default access privilege settings.",
		Syntax:        "\\ddp [pattern]",
		Handler:       ListDefaultPrivileges,
		CaseSensitive: true,
	})
}

func ListDefaultPrivileges(ctx context.Context, db database.Queryer, pattern string, verbose bool) (pgxspecial.SpecialCommandResult, error) {
	var sb strings.Builder
	args := []any{}

	sb.WriteString(`
	 SELECT pg_catalog.pg_get_userbyid(d.defaclrole) AS owner,
    n.nspname AS schema,
    CASE d.defaclobjtype WHEN 'r' THEN 'table'
                         WHEN 'S' THEN 'sequence'
                         WHEN 'f' THEN 'function'
                         WHEN 'T' THEN 'type'
                         WHEN 'n' THEN 'schema' END as type,
    pg_catalog.array_to_string(d.defaclacl, E'\n') AS access_privileges
    FROM pg_catalog.pg_default_acl d
        LEFT JOIN pg_catalog.pg_namespace n ON n.oid = d.defaclnamespace
	`)
	if pattern != "" {
		sb.WriteString(`
		 WHERE (n.nspname OPERATOR(pg_catalog.~) $1 COLLATE pg_catalog.default
            OR pg_catalog.pg_get_userbyid(d.defaclrole) OPERATOR(pg_catalog.~) $1 COLLATE pg_catalog.default)
		`)
		args = append(args, fmt.Sprintf("^(%s)$", pattern))
	}
	sb.WriteString("ORDER BY 1, 2, 3;")
	rows, err := db.Query(ctx, sb.String(), args...)
	return pgxspecial.RowResult{Rows: rows}, err
}
