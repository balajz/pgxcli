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
		Cmd:           "\\sf",
		Description:   "Show a function's definition.",
		Syntax:        "\\sf[+] FUNCNAME",
		Handler:       ShowFunctionDefinition,
		CaseSensitive: true,
	})
}

func ShowFunctionDefinition(ctx context.Context, db database.Queryer, pattern string, verbose bool) (pgxspecial.SpecialCommandResult, error) {
	var sql string
	if strings.Contains(pattern, "(") {
		sql = "SELECT $1::pg_catalog.regprocedure::pg_catalog.oid"
	} else {
		sql = "SELECT $1::pg_catalog.regproc::pg_catalog.oid"
	}

	var foid uint32
	err := db.QueryRow(ctx, sql, pattern).Scan(&foid)
	if err != nil {
		return nil, err
	}

	sql = "SELECT pg_catalog.pg_get_functiondef($1) as source"
	if !verbose {
		rows, err := db.Query(ctx, sql, foid)
		if err != nil {
			return nil, err
		}
		return pgxspecial.RowResult{Rows: rows}, nil

	}

	var source string
	err = db.QueryRow(ctx, sql, foid).Scan(&source)
	if err != nil {
		return nil, err
	}

	var sb strings.Builder
	lines := strings.Split(source, "\n")
	var rown int
	started := false

	for _, row := range lines {
		var prefix string
		if !started {
			if strings.HasPrefix(row, "AS ") {
				started = true
				rown = 1
				prefix = fmt.Sprintf("%-7d", rown)
			} else {
				prefix = fmt.Sprintf("%-7s", "")
			}
		} else {
			rown++
			prefix = fmt.Sprintf("%-7d", rown)
		}
		sb.WriteString(prefix + " " + row + "\n")
	}
	rows, err := db.Query(ctx, "SELECT $1 as source", sb.String())
	if err != nil {
		return nil, err
	}
	return pgxspecial.RowResult{Rows: rows}, nil

}
