package dbcommands

import (
	"context"
	"strconv"
	"strings"

	"github.com/balajz/pgxcli/pgxspecial"
	"github.com/balajz/pgxcli/pgxspecial/database"
	"github.com/jackc/pgx/v5"
)

func init() {
	pgxspecial.RegisterCommand(pgxspecial.SpecialCommandRegistry{
		Cmd:           "\\dx",
		Description:   "List extensions.",
		Syntax:        "\\dx[+] [pattern]",
		Handler:       ListExtensions,
		CaseSensitive: true,
	})
}

func ListExtensions(ctx context.Context, db database.Queryer, pattern string, verbose bool) (pgxspecial.SpecialCommandResult, error) {

	if verbose {
		extensions, err := findExtension(ctx, db, pattern)
		if err != nil {
			return nil, err
		}
		defer extensions.Close()
		var allExtensions []extensionDetail
		for extensions.Next() {
			var ext extensionDetail
			err := extensions.Scan(&ext.name, &ext.oid)
			if err != nil {
				return nil, err
			}
			allExtensions = append(allExtensions, ext)
		}

		var extDescriptions []pgxspecial.ExtensionVerboseResult

		for _, ext := range allExtensions {

			detailResult, err := describeExtension(ctx, db, ext.oid)
			if err != nil {
				return nil, err
			}
			defer detailResult.Close()

			var descriptions []string

			for detailResult.Next() {
				var desc string
				err := detailResult.Scan(&desc)
				if err != nil {
					return nil, err
				}
				descriptions = append(descriptions, desc)
			}

			extDescriptions = append(extDescriptions, pgxspecial.ExtensionVerboseResult{
				Name:        ext.name,
				Description: descriptions,
			})
		}
		return pgxspecial.ExtensionVerboseListResult{Results: extDescriptions}, nil
	}

	var sb strings.Builder
	args := []any{}
	argIndex := 1

	sb.WriteString(`
	 SELECT e.extname AS name,
             e.extversion AS version,
             n.nspname AS schema,
             c.description AS description
      FROM pg_catalog.pg_extension e
           LEFT JOIN pg_catalog.pg_namespace n
             ON n.oid = e.extnamespace
           LEFT JOIN pg_catalog.pg_description c
             ON c.objoid = e.oid
                AND c.classoid = 'pg_catalog.pg_extension'::pg_catalog.regclass
	`)

	if pattern != "" {
		_, tablePattern := sqlNamePattern(pattern)
		sb.WriteString(" WHERE e.extname ~ $" + strconv.Itoa(argIndex) + " ")
		args = append(args, tablePattern)
	}

	sb.WriteString(" ORDER BY 1, 2;")
	rows, err := db.Query(ctx, sb.String(), args...)
	return pgxspecial.RowResult{Rows: rows}, err
}

func findExtension(ctx context.Context, db database.Queryer, extName string) (pgx.Rows, error) {
	var sb strings.Builder
	var args []any

	sb.WriteString(`
			SELECT e.extname, e.oid
            FROM pg_catalog.pg_extension e
	`)

	if extName != "" {
		sb.WriteString(" WHERE e.extname = $1 ")
		args = append(args, extName)
	}

	sb.WriteString(" ORDER BY 1, 2;")

	rows, err := db.Query(ctx, sb.String(), args...)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func describeExtension(ctx context.Context, db database.Queryer, oid uint32) (pgx.Rows, error) {
	var sb strings.Builder

	sb.WriteString(`
	SELECT  pg_catalog.pg_describe_object(classid, objid, 0)
                    AS object_description
            FROM    pg_catalog.pg_depend
            WHERE   refclassid = 'pg_catalog.pg_extension'::pg_catalog.regclass
                    AND refobjid = $1
                    AND deptype = 'e'
            ORDER BY 1;
	`)

	rows, err := db.Query(ctx, sb.String(), oid)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

type extensionDetail struct {
	name string
	oid  uint32
}
