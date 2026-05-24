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
		Cmd:           "\\df",
		Description:   "List functions.",
		Syntax:        "\\df[+] [pattern]",
		Handler:       ListFunctions,
		CaseSensitive: true,
	})
}

func ListFunctions(ctx context.Context, db database.Queryer, pattern string, verbose bool) (pgxspecial.SpecialCommandResult, error) {
	var sb strings.Builder
	args := []any{}
	argIndex := 1

	sb.WriteString(`
	            SELECT  n.nspname as schema,
                    p.proname as name,
                    pg_catalog.pg_get_function_result(p.oid)
                      as "Result data type",
                    pg_catalog.pg_get_function_arguments(p.oid)
                      as "Argument data types",
                     CASE
                        WHEN p.prokind = 'a' THEN 'agg'
                        WHEN p.prokind = 'w' THEN 'window'
                        WHEN p.prorettype = 'pg_catalog.trigger'::pg_catalog.regtype
                            THEN 'trigger'
                        ELSE 'normal'
                    END as type 
	`)

	if verbose {
		sb.WriteString(`
		 ,CASE
                 WHEN p.provolatile = 'i' THEN 'immutable'
                 WHEN p.provolatile = 's' THEN 'stable'
                 WHEN p.provolatile = 'v' THEN 'volatile'
            END as "Volatility",
            pg_catalog.pg_get_userbyid(p.proowner) as owner,
          l.lanname as "Language",
          p.prosrc as "Source code",
          pg_catalog.obj_description(p.oid, 'pg_proc') as description 
		`)
	}

	sb.WriteString(`
	   FROM    pg_catalog.pg_proc p
            LEFT JOIN pg_catalog.pg_namespace n ON n.oid = p.pronamespace
	`)

	if verbose {
		sb.WriteString(`
		LEFT JOIN pg_catalog.pg_language l
			ON l.oid = p.prolang
		`)
	}

	sb.WriteString(`
	 WHERE  
	`)

	schemaPattern, funcPattern := sqlNamePattern(pattern)

	if schemaPattern != "" {
		sb.WriteString("  n.nspname ~ $" + strconv.Itoa(argIndex) + " ")
		args = append(args, schemaPattern)
		argIndex++
	} else {
		sb.WriteString(" pg_catalog.pg_function_is_visible(p.oid) ")
	}

	if funcPattern != "" {
		sb.WriteString(" AND p.proname ~ $" + strconv.Itoa(argIndex) + " ")
		args = append(args, funcPattern)
	}

	if schemaPattern == "" && funcPattern == "" {
		sb.WriteString(`
		AND n.nspname <> 'pg_catalog'
		AND n.nspname <> 'information_schema' 	
		`)
	}

	sb.WriteString(" ORDER BY 1, 2, 4;")

	rows, err := db.Query(ctx, sb.String(), args...)
	return pgxspecial.RowResult{Rows: rows}, err
}
