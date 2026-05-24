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
		Cmd:           "\\db",
		Description:   "List tablespaces.",
		Syntax:        "\\db[+] [pattern]",
		Handler:       ListTablespaces,
		CaseSensitive: true,
	})
}

func ListTablespaces(ctx context.Context, db database.Queryer, pattern string, verbose bool) (pgxspecial.SpecialCommandResult, error) {
	var sb strings.Builder
	args := []any{}
	argIndex := 1

	var isLocationSupported bool
	rows := db.QueryRow(ctx, `
	  SELECT EXISTS (
            SELECT * FROM pg_proc
            WHERE proname = 'pg_tablespace_location'
        )
	`)
	if err := rows.Scan(&isLocationSupported); err != nil {
		return nil, err
	}

	sb.WriteString(`
	SELECT
		n.spcname AS name,
		pg_catalog.pg_get_userbyid(n.spcowner) AS owner,
	`)
	if isLocationSupported {
		sb.WriteString("    pg_catalog.pg_tablespace_location(n.oid) AS location\n")
	} else {
		sb.WriteString("    'Not supported' AS location\n")
	}

	sb.WriteString(`
	FROM pg_catalog.pg_tablespace n
	`)

	if pattern != "" {
		_, tablePattern := sqlNamePattern(pattern)
		if tablePattern != "" {
			sb.WriteString(" WHERE n.spcname ~ $" + strconv.Itoa(argIndex) + " COLLATE pg_catalog.default ")
			args = append(args, tablePattern)
		}
	}

	sb.WriteString(" ORDER BY 1;")
	rowResult, err := db.Query(ctx, sb.String(), args...)
	return pgxspecial.RowResult{Rows: rowResult}, err

}
