package dbcommands

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/balajz/pgxcli/pgxspecial"
	"github.com/balajz/pgxcli/pgxspecial/database"
)

func init() {
	pgxspecial.RegisterCommand(pgxspecial.SpecialCommandRegistry{
		Cmd:           "\\d",
		Description:   "List or describe tables, views and sequences.",
		Syntax:        "\\d[+] [pattern]",
		Handler:       DescribeTableDetails,
		CaseSensitive: true,
	})

	pgxspecial.RegisterCommand(pgxspecial.SpecialCommandRegistry{
		Cmd:           "DESCRIBE",
		Description:   "List or describe tables, views and sequences.",
		Syntax:        "DESCRIBE [pattern]",
		Handler:       DescribeTableDetails,
		CaseSensitive: false,
	})
}

func DescribeTableDetails(ctx context.Context, db database.Queryer, pattern string, verbose bool) (pgxspecial.SpecialCommandResult, error) {
	if pattern == "" {
		return ListObjects(ctx, db, "", verbose, []string{"r", "p", "v", "m", "S", "f", ""})
	}

	schema, relname := sqlNamePattern(pattern)
	var sb strings.Builder
	args := []any{}
	argIndex := 1

	sb.WriteString(`
		SELECT c.oid, n.nspname, c.relname
		FROM pg_catalog.pg_class c
		LEFT JOIN pg_catalog.pg_namespace n ON n.oid = c.relnamespace
		WHERE 1=1
	`)

	if schema != "" {
		sb.WriteString(" AND n.nspname ~ $" + strconv.Itoa(argIndex))
		args = append(args, schema)
		argIndex++
	} else {
		sb.WriteString(" AND pg_catalog.pg_table_is_visible(c.oid)")
	}

	if relname != "" {
		sb.WriteString(" AND c.relname OPERATOR(pg_catalog.~) $" + strconv.Itoa(argIndex))
		args = append(args, relname)
	}

	sb.WriteString(" ORDER BY 2, 3")

	rows, err := db.Query(ctx, sb.String(), args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []pgxspecial.DescribeTableResult
	var targets []struct {
		oid     uint32
		schema  string
		relname string
	}

	for rows.Next() {
		var t struct {
			oid     uint32
			schema  string
			relname string
		}
		if err := rows.Scan(&t.oid, &t.schema, &t.relname); err != nil {
			return nil, err
		}
		targets = append(targets, t)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if len(targets) == 0 {
		return nil, fmt.Errorf("did not find any relation named %s", pattern)
	}

	for _, t := range targets {
		res, err := DescribeOneTableDetails(ctx, db, t.schema, t.relname, t.oid, verbose)
		if err != nil {
			return nil, err
		}
		results = append(results, res)
	}

	return pgxspecial.DescribeTableListResult{Results: results}, nil
}

type tableInfo struct {
	RelChecks     int
	RelKind       string
	HasIndex      bool
	HasRules      bool
	HasTriggers   bool
	HasOids       bool
	RelOptions    *string
	Tablespace    *string
	RelOfType     string
	Persistence   string
	IsPartition   bool
	RelToastRelID uint32
}

func DescribeOneTableDetails(ctx context.Context, db database.Queryer, schema, name string, oid uint32, verbose bool) (pgxspecial.DescribeTableResult, error) {
	var ti tableInfo
	// Assuming PG >= 12
	sql := `SELECT c.relchecks, c.relkind::text, c.relhasindex,
		c.relhasrules, c.relhastriggers, false as relhasoids,
		pg_catalog.array_to_string(c.reloptions || array(select 'toast.' || x from pg_catalog.unnest(tc.reloptions) x), ', '),
		c.reltablespace::text,
		CASE WHEN c.reloftype = 0 THEN '' ELSE c.reloftype::pg_catalog.regtype::pg_catalog.text END,
		c.relpersistence::text,
		c.relispartition
		FROM pg_catalog.pg_class c
		LEFT JOIN pg_catalog.pg_class tc ON (c.reltoastrelid = tc.oid)
		WHERE c.oid = $1`

	err := db.QueryRow(ctx, sql, oid).Scan(
		&ti.RelChecks, &ti.RelKind, &ti.HasIndex,
		&ti.HasRules, &ti.HasTriggers, &ti.HasOids,
		&ti.RelOptions, &ti.Tablespace, &ti.RelOfType,
		&ti.Persistence, &ti.IsPartition,
	)
	if err != nil {
		return pgxspecial.DescribeTableResult{}, err
	}

	// Columns
	headers, data, err := getTableColumns(ctx, db, oid, ti, verbose, schema, name)
	if err != nil {
		return pgxspecial.DescribeTableResult{}, err
	}

	// Footer
	meta, err := getTableFooter(ctx, db, oid, ti, verbose, schema)
	if err != nil {
		return pgxspecial.DescribeTableResult{}, err
	}

	return pgxspecial.DescribeTableResult{
		Columns:       headers,
		Data:          data,
		TableMetaData: meta,
	}, nil
}

//nolint:gocyclo
func getTableColumns(ctx context.Context, db database.Queryer, oid uint32, ti tableInfo, verbose bool, schema, name string) ([]string, [][]string, error) {
	var sb strings.Builder
	sb.WriteString(`SELECT a.attname,
    pg_catalog.format_type(a.atttypid, a.atttypmod),
    (SELECT substring(pg_catalog.pg_get_expr(d.adbin, d.adrelid, true) for 128)
                     FROM pg_catalog.pg_attrdef d
                     WHERE d.adrelid = a.attrelid AND d.adnum = a.attnum AND a.atthasdef),
    a.attnotnull,
    (SELECT c.collname FROM pg_catalog.pg_collation c, pg_catalog.pg_type t
                    WHERE c.oid = a.attcollation
                    AND t.oid = a.atttypid AND a.attcollation <> t.typcollation) AS attcollation,
    a.attidentity::text,
    a.attgenerated::text`)

	if ti.RelKind == "i" || ti.RelKind == "I" {
		fmt.Fprintf(&sb, `
		, CASE WHEN a.attnum <= (SELECT i.indnkeyatts FROM pg_catalog.pg_index i WHERE i.indexrelid = '%d') THEN 'yes' ELSE 'no' END AS is_key
		, pg_catalog.pg_get_indexdef(a.attrelid, a.attnum, TRUE) AS indexdef`, oid)
	} else {
		sb.WriteString(`, NULL AS is_key, NULL AS indexdef`)
	}

	if ti.RelKind == "f" {
		sb.WriteString(`, CASE WHEN attfdwoptions IS NULL THEN '' ELSE '(' ||
                array_to_string(ARRAY(SELECT quote_ident(option_name) ||  ' '
                || quote_literal(option_value)  FROM
                pg_options_to_table(attfdwoptions)), ', ') || ')' END AS attfdwoptions`)
	} else {
		sb.WriteString(`, NULL AS attfdwoptions`)
	}

	if verbose {
		sb.WriteString(`, a.attstorage::text`)
		if ti.RelKind == "r" || ti.RelKind == "i" || ti.RelKind == "I" || ti.RelKind == "m" || ti.RelKind == "f" || ti.RelKind == "p" {
			sb.WriteString(`, CASE WHEN a.attstattarget=-1 THEN NULL ELSE a.attstattarget END AS attstattarget`)
		} else {
			sb.WriteString(`, NULL AS attstattarget`)
		}
		if ti.RelKind == "r" || ti.RelKind == "v" || ti.RelKind == "m" || ti.RelKind == "f" || ti.RelKind == "p" || ti.RelKind == "c" {
			sb.WriteString(`, pg_catalog.col_description(a.attrelid, a.attnum)`)
		} else {
			sb.WriteString(`, NULL AS attdescr`)
		}
	} else {
		sb.WriteString(`, NULL AS attstorage, NULL AS attstattarget, NULL AS attdescr`)
	}

	fmt.Fprintf(&sb, ` FROM pg_catalog.pg_attribute a WHERE a.attrelid = '%d' AND
    a.attnum > 0 AND NOT a.attisdropped ORDER BY a.attnum;`, oid)

	rows, err := db.Query(ctx, sb.String())
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	var headers []string
	headers = append(headers, "Column", "Type")

	showModifiers := false
	if ti.RelKind == "r" || ti.RelKind == "p" || ti.RelKind == "v" || ti.RelKind == "m" || ti.RelKind == "f" || ti.RelKind == "c" {
		headers = append(headers, "Modifiers")
		showModifiers = true
	}

	// Fetch sequence values BEFORE iterating rows to avoid "conn busy" error
	var seqValues []any
	if ti.RelKind == "S" {
		headers = append(headers, "Value")
		// Close the rows first to free the connection
		rows.Close()

		// Now fetch sequence values
		seqRows, err := db.Query(ctx, fmt.Sprintf(`SELECT * FROM "%s"."%s"`, schema, name))
		if err != nil {
			return nil, nil, err
		}
		defer seqRows.Close()
		if seqRows.Next() {
			seqValues, err = seqRows.Values()
			if err != nil {
				return nil, nil, err
			}
		}
		seqRows.Close()

		// Re-open rows query for the main iteration
		rows, err = db.Query(ctx, sb.String())
		if err != nil {
			return nil, nil, err
		}
		defer rows.Close()
	}

	if ti.RelKind == "i" {
		headers = append(headers, "Definition")
	}

	if ti.RelKind == "f" {
		headers = append(headers, "FDW Options")
	}

	if verbose {
		headers = append(headers, "Storage")
		if ti.RelKind == "r" || ti.RelKind == "m" || ti.RelKind == "f" {
			headers = append(headers, "Stats target")
		}
		if ti.RelKind == "r" || ti.RelKind == "v" || ti.RelKind == "m" || ti.RelKind == "c" || ti.RelKind == "f" {
			headers = append(headers, "Description")
		}
	}

	var data [][]string
	rowIndex := 0
	for rows.Next() {
		var attname, atttype string
		var attrdef *string
		var attnotnull bool
		var attcollation *string
		var attidentity, attgenerated string
		var isKey, indexdef, attfdwoptions *string
		var attstorage *string
		var attstattarget *int32
		var attdescr *string

		err := rows.Scan(&attname, &atttype, &attrdef, &attnotnull, &attcollation, &attidentity, &attgenerated,
			&isKey, &indexdef, &attfdwoptions, &attstorage, &attstattarget, &attdescr)
		if err != nil {
			return nil, nil, err
		}

		var row []string
		row = append(row, attname, atttype)

		if showModifiers {
			modifier := ""
			if attcollation != nil {
				modifier += fmt.Sprintf(" collate %s", *attcollation)
			}
			if attnotnull {
				modifier += " not null"
			}
			if attrdef != nil {
				modifier += fmt.Sprintf(" default %s", *attrdef)
			}
			if attidentity == "a" {
				modifier += " generated always as identity"
			} else if attidentity == "d" {
				modifier += " generated by default as identity"
			} else if attgenerated == "s" {
				if attrdef != nil {
					modifier += fmt.Sprintf(" generated always as (%s) stored", *attrdef)
				}
			}
			row = append(row, modifier)
		}

		if ti.RelKind == "S" {
			if rowIndex < len(seqValues) {
				row = append(row, fmt.Sprintf("%v", seqValues[rowIndex]))
			} else {
				row = append(row, "")
			}
		}

		if ti.RelKind == "i" {
			if indexdef != nil {
				row = append(row, *indexdef)
			} else {
				row = append(row, "")
			}
		}

		if ti.RelKind == "f" {
			if attfdwoptions != nil {
				row = append(row, *attfdwoptions)
			} else {
				row = append(row, "")
			}
		}

		if verbose {
			if attstorage != nil {
				switch (*attstorage)[0] {
				case 'p':
					row = append(row, "plain")
				case 'm':
					row = append(row, "main")
				case 'x':
					row = append(row, "extended")
				case 'e':
					row = append(row, "external")
				default:
					row = append(row, "???")
				}
			} else {
				row = append(row, "")
			}

			if ti.RelKind == "r" || ti.RelKind == "m" || ti.RelKind == "f" {
				if attstattarget != nil {
					row = append(row, fmt.Sprintf("%d", *attstattarget))
				} else {
					row = append(row, "")
				}
			}

			if ti.RelKind == "r" || ti.RelKind == "v" || ti.RelKind == "m" || ti.RelKind == "c" || ti.RelKind == "f" {
				if attdescr != nil {
					row = append(row, *attdescr)
				} else {
					row = append(row, "")
				}
			}
		}
		data = append(data, row)
		rowIndex++
	}
	return headers, data, nil
}

//nolint:gocyclo
func getTableFooter(ctx context.Context, db database.Queryer, oid uint32, ti tableInfo, verbose bool, schema string) (pgxspecial.TableFooterMeta, error) {
	var meta pgxspecial.TableFooterMeta
	var err error

	if (ti.RelKind == "v" || ti.RelKind == "m") && verbose {
		meta.ViewDefinition = getViewDefinition(ctx, db, oid)
	}

	switch ti.RelKind {
	case "i":
		meta.Options, err = getIndexFooter(ctx, db, oid, schema)
		if err != nil {
			return meta, err
		}
	case "S":
		meta.OwnedBy = getSequenceOwner(ctx, db, oid)
	case "r", "p", "m", "f":
		if ti.HasIndex {
			meta.Indexes, err = getIndexes(ctx, db, oid)
			if err != nil {
				return meta, err
			}
		}
		if ti.RelChecks > 0 {
			meta.CheckConstraints, err = getCheckConstraints(ctx, db, oid)
			if err != nil {
				return meta, err
			}
		}
		if ti.HasTriggers {
			meta.ForeignKeys, err = getForeignKeys(ctx, db, oid)
			if err != nil {
				return meta, err
			}
			meta.ReferencedBy, err = getReferencedBy(ctx, db, oid)
			if err != nil {
				return meta, err
			}
		}
		if ti.HasRules && ti.RelKind != "m" {
			meta.RulesEnabled, meta.RulesDisabled, meta.RulesAlways, meta.RulesReplica, err = getRules(ctx, db, oid)
			if err != nil {
				return meta, err
			}
		}
		if ti.IsPartition {
			meta.PartitionOf, meta.PartitionConstraints, err = getPartitionInfo(ctx, db, oid)
			if err != nil {
				return meta, err
			}
		}
		if ti.RelKind == "p" {
			meta.PartitionKey, meta.Partitions, meta.PartitionsSummary, err = getPartitionDetails(ctx, db, oid, verbose)
			if err != nil {
				return meta, err
			}
		}
	}

	if ti.HasTriggers {
		meta.TriggersEnabled, meta.TriggersDisabled, meta.TriggersAlways, meta.TriggersReplica, err = getTriggers(ctx, db, oid)
		if err != nil {
			return meta, err
		}
	}

	if ti.RelKind == "r" || ti.RelKind == "m" || ti.RelKind == "f" {
		if ti.RelKind == "f" {
			meta.Server, meta.FDWOptions = getForeignTableInfo(ctx, db, oid)
		}
		if !ti.IsPartition {
			meta.Inherits, err = getInherits(ctx, db, oid)
			if err != nil {
				return meta, err
			}
		}
		meta.ChildTables, meta.ChildTablesSummary, err = getChildTables(ctx, db, oid, verbose)
		if err != nil {
			return meta, err
		}
		if ti.RelOfType != "" {
			meta.TypedTableOf = &ti.RelOfType
		}
		if verbose && ti.RelKind != "m" {
			meta.HasOIDs = &ti.HasOids
		}
	}

	if verbose && ti.RelOptions != nil && *ti.RelOptions != "" {
		meta.Options = ti.RelOptions
	}

	return meta, nil
}

func getViewDefinition(ctx context.Context, db database.Queryer, oid uint32) *string {
	var viewDef string
	err := db.QueryRow(ctx, fmt.Sprintf("SELECT pg_catalog.pg_get_viewdef('%d'::pg_catalog.oid, true)", oid)).Scan(&viewDef)
	if err != nil {
		return nil // Ignore error?
	}
	return &viewDef
}

func getIndexFooter(ctx context.Context, db database.Queryer, oid uint32, schema string) (*string, error) {
	sql := `SELECT i.indisunique, i.indisprimary, i.indisclustered, i.indisvalid,
			(NOT i.indimmediate) AND EXISTS (SELECT 1 FROM pg_catalog.pg_constraint WHERE conrelid = i.indrelid AND conindid = i.indexrelid AND contype IN ('p','u','x') AND condeferrable) AS condeferrable,
			(NOT i.indimmediate) AND EXISTS (SELECT 1 FROM pg_catalog.pg_constraint WHERE conrelid = i.indrelid AND conindid = i.indexrelid AND contype IN ('p','u','x') AND condeferred) AS condeferred,
			a.amname, c2.relname, pg_catalog.pg_get_expr(i.indpred, i.indrelid, true)
			FROM pg_catalog.pg_index i, pg_catalog.pg_class c, pg_catalog.pg_class c2, pg_catalog.pg_am a
			WHERE i.indexrelid = c.oid AND c.oid = $1 AND c.relam = a.oid AND i.indrelid = c2.oid`

	var indisunique, indisprimary, indisclustered, indisvalid, deferrable, deferred bool
	var indamname, indtable string
	var indpred *string
	err := db.QueryRow(ctx, sql, oid).Scan(&indisunique, &indisprimary, &indisclustered, &indisvalid, &deferrable, &deferred, &indamname, &indtable, &indpred)
	if err != nil {
		return nil, err
	}

	var statusParts []string
	if indisprimary {
		statusParts = append(statusParts, "primary key")
	} else if indisunique {
		statusParts = append(statusParts, "unique")
	}
	statusParts = append(statusParts, indamname)
	statusParts = append(statusParts, fmt.Sprintf(`for table "%s.%s"`, schema, indtable))
	if indpred != nil {
		statusParts = append(statusParts, fmt.Sprintf("predicate (%s)", *indpred))
	}
	if indisclustered {
		statusParts = append(statusParts, "clustered")
	}
	if !indisvalid {
		statusParts = append(statusParts, "invalid")
	}
	if deferrable {
		statusParts = append(statusParts, "deferrable")
	}
	if deferred {
		statusParts = append(statusParts, "initially deferred")
	}
	summary := strings.Join(statusParts, ", ")
	return &summary, nil
}

func getSequenceOwner(ctx context.Context, db database.Queryer, oid uint32) *string {
	sql := `SELECT pg_catalog.quote_ident(nspname) || '.' || pg_catalog.quote_ident(relname) || '.' || pg_catalog.quote_ident(attname)
		FROM pg_catalog.pg_class c
		INNER JOIN pg_catalog.pg_depend d ON c.oid=d.refobjid
		INNER JOIN pg_catalog.pg_namespace n ON n.oid=c.relnamespace
		INNER JOIN pg_catalog.pg_attribute a ON (a.attrelid=c.oid AND a.attnum=d.refobjsubid)
		WHERE d.classid='pg_catalog.pg_class'::pg_catalog.regclass
		AND d.refclassid='pg_catalog.pg_class'::pg_catalog.regclass
		AND d.objid=$1 AND d.deptype='a'`
	var seqOwner string
	if err := db.QueryRow(ctx, sql, oid).Scan(&seqOwner); err != nil {
		return nil
	}

	return &seqOwner
}

//nolint:gocyclo
func getIndexes(ctx context.Context, db database.Queryer, oid uint32) ([]string, error) {
	sql := `SELECT c2.relname, i.indisprimary, i.indisunique, i.indisclustered, i.indisvalid,
			pg_catalog.pg_get_indexdef(i.indexrelid, 0, true),
			pg_catalog.pg_get_constraintdef(con.oid, true),
			contype::text, condeferrable, condeferred, c2.reltablespace
			FROM pg_catalog.pg_class c, pg_catalog.pg_class c2, pg_catalog.pg_index i
			LEFT JOIN pg_catalog.pg_constraint con ON conrelid = i.indrelid AND conindid = i.indexrelid AND contype IN ('p','u','x')
			WHERE c.oid = $1 AND c.oid = i.indrelid AND i.indexrelid = c2.oid
			ORDER BY i.indisprimary DESC, i.indisunique DESC, c2.relname`

	rows, err := db.Query(ctx, sql, oid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var indexes []string
	for rows.Next() {
		var relname string
		var primary, unique, clustered, valid bool
		var indexDef string
		var constraintDef *string
		var conType *string
		var deferrable, deferred *bool
		var tablespace *uint32

		if err := rows.Scan(&relname, &primary, &unique, &clustered, &valid, &indexDef, &constraintDef, &conType, &deferrable, &deferred, &tablespace); err != nil {
			return nil, err
		}

		entry := fmt.Sprintf(`"%s"`, relname)
		if conType != nil && *conType == "x" {
			entry += " " + *constraintDef
		} else {
			if primary {
				entry += " PRIMARY KEY,"
			} else if unique {
				if conType != nil && *conType == "u" {
					entry += " UNIQUE CONSTRAINT,"
				} else {
					entry += " UNIQUE,"
				}
			}

			usingPos := strings.Index(indexDef, " USING ")
			if usingPos >= 0 {
				entry += " " + indexDef[usingPos+7:]
			}

			if deferrable != nil && *deferrable {
				entry += " DEFERRABLE"
			}
			if deferred != nil && *deferred {
				entry += " INITIALLY DEFERRED"
			}
		}
		if clustered {
			entry += " CLUSTER"
		}
		if !valid {
			entry += " INVALID"
		}
		indexes = append(indexes, entry)
	}
	return indexes, nil
}

func getCheckConstraints(ctx context.Context, db database.Queryer, oid uint32) ([]string, error) {
	sql := `SELECT r.conname, pg_catalog.pg_get_constraintdef(r.oid, true)
		FROM pg_catalog.pg_constraint r
		WHERE r.conrelid = $1 AND r.contype = 'c' ORDER BY 1`
	rows, err := db.Query(ctx, sql, oid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var constraints []string
	for rows.Next() {
		var conname, condef string
		_ = rows.Scan(&conname, &condef)
		constraints = append(constraints, fmt.Sprintf(`"%s" %s`, conname, condef))
	}
	return constraints, nil
}

func getForeignKeys(ctx context.Context, db database.Queryer, oid uint32) ([]string, error) {
	sql := `SELECT conname, pg_catalog.pg_get_constraintdef(r.oid, true)
		FROM pg_catalog.pg_constraint r
		WHERE r.conrelid = $1 AND r.contype = 'f' ORDER BY 1`
	rows, err := db.Query(ctx, sql, oid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var fks []string
	for rows.Next() {
		var conname, condef string
		_ = rows.Scan(&conname, &condef)
		fks = append(fks, fmt.Sprintf(`"%s" %s`, conname, condef))
	}
	return fks, nil
}

func getReferencedBy(ctx context.Context, db database.Queryer, oid uint32) ([]string, error) {
	sql := `SELECT conrelid::pg_catalog.regclass::text, conname, pg_catalog.pg_get_constraintdef(c.oid, true)
		FROM pg_catalog.pg_constraint c
		WHERE c.confrelid = $1 AND c.contype = 'f' ORDER BY 1`
	rows, err := db.Query(ctx, sql, oid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var refs []string
	for rows.Next() {
		var conrelid, conname, condef string
		_ = rows.Scan(&conrelid, &conname, &condef)
		refs = append(refs, fmt.Sprintf(`TABLE "%s" CONSTRAINT "%s" %s`, conrelid, conname, condef))
	}
	return refs, nil
}

func getRules(ctx context.Context, db database.Queryer, oid uint32) (enabled, disabled, always, replica []string, err error) {
	sql := `SELECT r.rulename, trim(trailing ';' from pg_catalog.pg_get_ruledef(r.oid, true)), ev_enabled::text
		FROM pg_catalog.pg_rewrite r WHERE r.ev_class = $1 ORDER BY 1`
	rows, err := db.Query(ctx, sql, oid)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var rulename, ruledef, evEnabled string
		if err := rows.Scan(&rulename, &ruledef, &evEnabled); err != nil {
			return nil, nil, nil, nil, err
		}
		switch evEnabled {
		case "O":
			enabled = append(enabled, ruledef)
		case "D":
			disabled = append(disabled, ruledef)
		case "A":
			always = append(always, ruledef)
		case "R":
			replica = append(replica, ruledef)
		}
	}
	return
}

func getTriggers(ctx context.Context, db database.Queryer, oid uint32) (enabled, disabled, always, replica []string, err error) {
	sql := `SELECT t.tgname, pg_catalog.pg_get_triggerdef(t.oid, true), t.tgenabled::text
		FROM pg_catalog.pg_trigger t
		WHERE t.tgrelid = $1 AND NOT t.tgisinternal ORDER BY 1`
	rows, err := db.Query(ctx, sql, oid)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var tgname, tgdef, tgenabled string
		if err := rows.Scan(&tgname, &tgdef, &tgenabled); err != nil {
			return nil, nil, nil, nil, err
		}

		triggerPos := strings.Index(tgdef, " TRIGGER ")
		if triggerPos >= 0 {
			tgdef = tgdef[triggerPos+9:]
		}

		switch tgenabled {
		case "O":
			enabled = append(enabled, tgdef)
		case "D":
			disabled = append(disabled, tgdef)
		case "A":
			always = append(always, tgdef)
		case "R":
			replica = append(replica, tgdef)
		}
	}
	return
}

func getPartitionInfo(ctx context.Context, db database.Queryer, oid uint32) (partOf, partConstraints []string, err error) {
	sql := `select quote_ident(np.nspname) || '.' || quote_ident(cp.relname) || ' ' || pg_get_expr(cc.relpartbound, cc.oid, true),
		pg_get_partition_constraintdef(cc.oid)
		from pg_inherits i
		inner join pg_class cp on cp.oid = i.inhparent
		inner join pg_namespace np on np.oid = cp.relnamespace
		inner join pg_class cc on cc.oid = i.inhrelid
		where cc.oid = $1`
	rows, err := db.Query(ctx, sql, oid)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var po, pc string
		_ = rows.Scan(&po, &pc)
		partOf = append(partOf, po)
		partConstraints = append(partConstraints, pc)
	}
	return
}

func getPartitionDetails(ctx context.Context, db database.Queryer, oid uint32, verbose bool) (partKey *string, partitions []string, summary *string, err error) {
	// Partition key
	sql := fmt.Sprintf("select pg_get_partkeydef(%d)", oid)
	var pk string
	if err := db.QueryRow(ctx, sql).Scan(&pk); err == nil {
		partKey = &pk
	}

	// Partitions
	sql = fmt.Sprintf(`select quote_ident(n.nspname) || '.' || quote_ident(c.relname) || ' ' || pg_get_expr(c.relpartbound, c.oid, true)
		from pg_inherits i
		inner join pg_class c on c.oid = i.inhrelid
		inner join pg_namespace n on n.oid = c.relnamespace
		where i.inhparent = %d order by 1`, oid)
	rows, err := db.Query(ctx, sql)
	if err != nil {
		return nil, nil, nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var p string
		_ = rows.Scan(&p)
		partitions = append(partitions, p)
	}

	if !verbose {
		s := fmt.Sprintf("Number of partitions %d: (Use \\d+ to list them.)", len(partitions))
		summary = &s
		partitions = nil
	}
	return
}

func getForeignTableInfo(ctx context.Context, db database.Queryer, oid uint32) (server *string, fdwOptions *string) {
	sql := `SELECT s.srvname,
		array_to_string(ARRAY(SELECT quote_ident(option_name) ||  ' ' || quote_literal(option_value)
		FROM pg_options_to_table(ftoptions)),  ', ')
		FROM pg_catalog.pg_foreign_table f, pg_catalog.pg_foreign_server s
		WHERE f.ftrelid = $1 AND s.oid = f.ftserver`
	var srvname string
	var opts *string
	if err := db.QueryRow(ctx, sql, oid).Scan(&srvname, &opts); err == nil {
		server = &srvname
		if opts != nil && *opts != "" {
			o := fmt.Sprintf("(%s)", *opts)
			fdwOptions = &o
		}
	}
	return
}

func getInherits(ctx context.Context, db database.Queryer, oid uint32) ([]string, error) {
	sql := `SELECT c.oid::pg_catalog.regclass::text
		FROM pg_catalog.pg_class c, pg_catalog.pg_inherits i
		WHERE c.oid = i.inhparent AND i.inhrelid = $1 ORDER BY inhseqno`
	rows, err := db.Query(ctx, sql, oid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var inherits []string
	for rows.Next() {
		var inh string
		_ = rows.Scan(&inh)
		inherits = append(inherits, inh)
	}
	return inherits, nil
}

func getChildTables(ctx context.Context, db database.Queryer, oid uint32, verbose bool) (children []string, summary *string, err error) {
	sql := `SELECT c.oid::pg_catalog.regclass::text
		FROM pg_catalog.pg_class c, pg_catalog.pg_inherits i
		WHERE c.oid = i.inhrelid AND i.inhparent = $1
		ORDER BY c.oid::pg_catalog.regclass::pg_catalog.text`
	rows, err := db.Query(ctx, sql, oid)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var child string
		_ = rows.Scan(&child)
		children = append(children, child)
	}

	if !verbose {
		if len(children) > 0 {
			s := fmt.Sprintf("Number of child tables: %d (Use \\d+ to list them.)", len(children))
			summary = &s
			children = nil
		}
	}
	return
}
