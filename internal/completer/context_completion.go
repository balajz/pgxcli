package completer

import (
	"fmt"
	"sort"
	"strings"
	"unicode"
)

type CompletionType int

const (
	CompletionTypeKeyword CompletionType = iota
	CompletionTypeColName
	CompletionTypeTableReference
	CompletionTypeWhereCondition
	CompletionTypeJoinClause
	CompletionTypeJoinOn
	CompletionTypeInsertColumn
)

type CompletionContext struct {
	Types      []CompletionType
	Prefix     string
	Qualifier  string
	Statement  string
	TableRefs  []tableRef
	LastTable  *tableRef
	IsJoinOnly bool
}

type tableRef struct {
	Schema string
	Name   string
	Alias  string
}

// EnableSmartCompletion sets up metadata-backed autocompletion.
func (c *Completer) EnableSmartCompletion(executor DatabaseExecutor) error {
	c.executor = executor
	c.smartCompletion = executor != nil
	if executor == nil {
		return nil
	}

	if err := c.RefreshMetaData(); err != nil {
		return err
	}

	return nil
}

// RefreshMetaData loads table, column, schema and foreign-key metadata for completion.
func (c *Completer) RefreshMetaData() error {
	if c.executor == nil {
		return nil
	}

	schemas, err := c.executor.Schemas()
	if err != nil {
		return fmt.Errorf("loading schemas for completion: %w", err)
	}
	c.ExtendSchemas(schemas)

	tables, err := c.executor.Tables()
	if err != nil {
		return fmt.Errorf("loading tables for completion: %w", err)
	}
	c.ExtendTables(tables)

	tableCols, err := c.executor.TableColumns()
	if err != nil {
		return fmt.Errorf("loading table columns for completion: %w", err)
	}
	c.ExtendColumns(tableCols, false)

	viewCols, err := c.executor.ViewColumns()
	if err == nil {
		c.ExtendColumns(viewCols, true)
	}

	fks, err := c.executor.ForeignKeys()
	if err == nil {
		c.ExtendForeignKeys(fks)
	}

	return nil
}

// Complete returns context-aware suggestions for SQL input at a cursor location.
func (c *Completer) Complete(input string, cursor int) (string, []string) {
	if cursor < 0 || cursor > len(input) {
		cursor = len(input)
	}

	ctx := c.buildContext(input, cursor)
	candidates := c.candidatesForContext(ctx)
	candidates = filterAndUniqueByPrefix(candidates, ctx.Prefix)
	if len(candidates) == 0 {
		return "Keywords", nil
	}

	return completionCategory(ctx.Types), candidates
}

func (c *Completer) buildContext(input string, cursor int) CompletionContext {
	statement, beforeCursorStatement := currentStatementAtCursor(input, cursor)
	prefix, qualifier := identifierPrefix(beforeCursorStatement)
	refs := extractTableReferences(statement)
	types := getCompletionTypes(beforeCursorStatement, qualifier)

	var last *tableRef
	if len(refs) > 0 {
		lastRef := refs[len(refs)-1]
		last = &lastRef
	}

	return CompletionContext{
		Types:      types,
		Prefix:     prefix,
		Qualifier:  qualifier,
		Statement:  statement,
		TableRefs:  refs,
		LastTable:  last,
		IsJoinOnly: hasCompletionType(types, CompletionTypeJoinClause),
	}
}

func currentStatementAtCursor(input string, cursor int) (string, string) {
	lastSemicolon := strings.LastIndex(input[:cursor], ";")
	start := 0
	if lastSemicolon >= 0 {
		start = lastSemicolon + 1
	}

	nextSemicolonOffset := strings.Index(input[cursor:], ";")
	end := len(input)
	if nextSemicolonOffset >= 0 {
		end = cursor + nextSemicolonOffset
	}

	return input[start:end], input[start:cursor]
}

func completionCategory(types []CompletionType) string {
	if len(types) == 0 {
		return "Keywords"
	}

	switch types[0] {
	case CompletionTypeColName, CompletionTypeWhereCondition, CompletionTypeJoinOn, CompletionTypeInsertColumn:
		return "Columns"
	case CompletionTypeTableReference, CompletionTypeJoinClause:
		return "Tables"
	default:
		return "Keywords"
	}
}

func hasCompletionType(types []CompletionType, target CompletionType) bool {
	for _, t := range types {
		if t == target {
			return true
		}
	}
	return false
}

func getCompletionTypes(statement, qualifier string) []CompletionType {
	trimmed := strings.TrimRightFunc(statement, unicode.IsSpace)
	upper := strings.ToUpper(trimmed)
	clause, tail := lastClauseAndTail(upper)
	tail = strings.TrimSpace(tail)

	if clause == "JOIN" && tail == "" {
		return []CompletionType{CompletionTypeJoinClause, CompletionTypeTableReference, CompletionTypeKeyword}
	}

	if clause == "FROM" || clause == "UPDATE" || clause == "INTO" || clause == "DELETE FROM" {
		return []CompletionType{CompletionTypeTableReference, CompletionTypeKeyword}
	}

	if clause == "WHERE" || clause == "HAVING" {
		return []CompletionType{CompletionTypeWhereCondition, CompletionTypeColName, CompletionTypeKeyword}
	}

	if clause == "ON" {
		return []CompletionType{CompletionTypeJoinOn, CompletionTypeColName, CompletionTypeKeyword}
	}

	if strings.Contains(upper, "INSERT INTO") && strings.LastIndex(upper, "(") > strings.LastIndex(upper, ")") {
		return []CompletionType{CompletionTypeInsertColumn, CompletionTypeColName, CompletionTypeKeyword}
	}

	switch clause {
	case "SELECT", "SET", "GROUP BY", "ORDER BY":
		return []CompletionType{CompletionTypeColName, CompletionTypeKeyword}
	case "FROM", "UPDATE", "INTO", "DELETE FROM":
		return []CompletionType{CompletionTypeTableReference, CompletionTypeKeyword}
	case "JOIN":
		return []CompletionType{CompletionTypeJoinClause, CompletionTypeTableReference, CompletionTypeKeyword}
	case "WHERE", "HAVING", "ON":
		return []CompletionType{CompletionTypeWhereCondition, CompletionTypeColName, CompletionTypeKeyword}
	}

	if qualifier != "" {
		return []CompletionType{CompletionTypeColName, CompletionTypeKeyword}
	}

	return []CompletionType{CompletionTypeKeyword}
}

func lastClauseAndTail(upper string) (string, string) {
	keywords := []string{
		"GROUP BY", "ORDER BY", "DELETE FROM", "INSERT INTO", "SELECT", "FROM", "WHERE", "HAVING", "JOIN", "ON", "UPDATE", "SET", "INTO",
	}

	lastPos := -1
	lastKeyword := ""
	for _, kw := range keywords {
		pos := strings.LastIndex(upper, kw)
		if pos > lastPos {
			lastPos = pos
			lastKeyword = kw
		}
	}

	if lastPos == -1 {
		return "", ""
	}
	return lastKeyword, upper[lastPos+len(lastKeyword):]
}

func identifierPrefix(statement string) (prefix, qualifier string) {
	if statement == "" {
		return "", ""
	}

	idx := len(statement) - 1
	for idx >= 0 {
		r := rune(statement[idx])
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' || r == '.' || r == '"' {
			idx--
			continue
		}
		break
	}

	token := strings.TrimSpace(statement[idx+1:])
	if token == "" {
		return "", ""
	}

	lastDot := strings.LastIndex(token, ".")
	if lastDot >= 0 {
		qualifier = strings.Trim(token[:lastDot], `"`)
		prefix = strings.Trim(token[lastDot+1:], `"`)
		return prefix, qualifier
	}

	return strings.Trim(token, `"`), ""
}

func extractTableReferences(statement string) []tableRef {
	tokens := tokenizeSQL(statement)
	refs := make([]tableRef, 0)

	for i := 0; i < len(tokens); i++ {
		if !tokens[i].isWord {
			continue
		}

		kw := strings.ToUpper(tokens[i].value)
		if kw != "FROM" && kw != "JOIN" && kw != "UPDATE" && kw != "INTO" {
			continue
		}

		ref, consumed := parseTableRef(tokens, i+1)
		if consumed == 0 {
			continue
		}
		i += consumed
		refs = append(refs, ref)
	}

	return refs
}

type sqlToken struct {
	value  string
	isWord bool
}

func tokenizeSQL(input string) []sqlToken {
	tokens := make([]sqlToken, 0)
	for i := 0; i < len(input); {
		r := rune(input[i])
		if unicode.IsSpace(r) {
			i++
			continue
		}

		if r == '"' {
			j := i + 1
			for j < len(input) && rune(input[j]) != '"' {
				j++
			}
			if j < len(input) {
				j++
			}
			tokens = append(tokens, sqlToken{value: input[i:j], isWord: true})
			i = j
			continue
		}

		if unicode.IsLetter(r) || r == '_' {
			j := i + 1
			for j < len(input) {
				nr := rune(input[j])
				if unicode.IsLetter(nr) || unicode.IsDigit(nr) || nr == '_' {
					j++
					continue
				}
				break
			}
			tokens = append(tokens, sqlToken{value: input[i:j], isWord: true})
			i = j
			continue
		}

		tokens = append(tokens, sqlToken{value: string(r), isWord: false})
		i++
	}

	return tokens
}

func parseTableRef(tokens []sqlToken, start int) (tableRef, int) {
	if start >= len(tokens) {
		return tableRef{}, 0
	}

	i := start
	if tokens[i].value == "(" {
		depth := 1
		i++
		for i < len(tokens) && depth > 0 {
			switch tokens[i].value {
			case "(":
				depth++
			case ")":
				depth--
			}
			i++
		}
		if i < len(tokens) && tokens[i].isWord {
			return tableRef{Name: trimIdentifier(tokens[i].value), Alias: trimIdentifier(tokens[i].value)}, i - start
		}
		return tableRef{}, 0
	}

	if !tokens[i].isWord {
		return tableRef{}, 0
	}

	first := trimIdentifier(tokens[i].value)
	i++

	schema := ""
	table := first
	if i+1 < len(tokens) && tokens[i].value == "." && tokens[i+1].isWord {
		schema = first
		table = trimIdentifier(tokens[i+1].value)
		i += 2
	}

	alias := table
	if i < len(tokens) && tokens[i].isWord && strings.EqualFold(tokens[i].value, "AS") {
		i++
	}
	if i < len(tokens) && tokens[i].isWord && !isSQLKeyword(tokens[i].value) {
		alias = trimIdentifier(tokens[i].value)
		i++
	}

	return tableRef{Schema: schema, Name: table, Alias: alias}, i - start
}

func trimIdentifier(s string) string {
	return strings.Trim(s, `"`)
}

func isSQLKeyword(token string) bool {
	switch strings.ToUpper(token) {
	case "SELECT", "FROM", "WHERE", "JOIN", "ON", "AND", "OR", "GROUP", "ORDER", "BY", "LEFT", "RIGHT", "INNER", "OUTER", "FULL", "CROSS", "LIMIT", "OFFSET", "SET", "VALUES":
		return true
	default:
		return false
	}
}

func (c *Completer) candidatesForContext(ctx CompletionContext) []string {
	c.metadata.mu.RLock()
	defer c.metadata.mu.RUnlock()

	candidates := make([]string, 0)
	for _, typ := range ctx.Types {
		switch typ {
		case CompletionTypeColName, CompletionTypeWhereCondition, CompletionTypeJoinOn, CompletionTypeInsertColumn:
			candidates = append(candidates, c.columnCandidates(ctx)...)
		case CompletionTypeTableReference:
			candidates = append(candidates, c.tableCandidates(ctx)...)
		case CompletionTypeJoinClause:
			candidates = append(candidates, c.joinCandidates(ctx)...)
		case CompletionTypeKeyword:
			candidates = append(candidates, c.metadata.KeyWords...)
		}
	}

	return candidates
}

func (c *Completer) columnCandidates(ctx CompletionContext) []string {
	if len(ctx.TableRefs) == 0 {
		return nil
	}

	out := make([]string, 0)
	if ctx.Qualifier != "" {
		for _, ref := range ctx.TableRefs {
			if strings.EqualFold(ref.Alias, ctx.Qualifier) || strings.EqualFold(ref.Name, ctx.Qualifier) {
				out = append(out, c.columnsForRef(ref)...)
				break
			}
		}
		return out
	}

	for _, ref := range ctx.TableRefs {
		cols := c.columnsForRef(ref)
		for _, col := range cols {
			if ref.Alias != "" {
				out = append(out, ref.Alias+"."+col)
			}
			out = append(out, col)
		}
	}

	return out
}

func (c *Completer) columnsForRef(ref tableRef) []string {
	out := make([]string, 0)

	lookupSchema := ref.Schema
	if lookupSchema != "" {
		if schemaTables, ok := c.metadata.Tables[lookupSchema]; ok {
			if table, ok := schemaTables[ref.Name]; ok {
				for _, col := range sortedColumns(table.Columns) {
					out = append(out, col)
				}
				return out
			}
		}
	}

	for _, schemaTables := range c.metadata.Tables {
		if table, ok := schemaTables[ref.Name]; ok {
			for _, col := range sortedColumns(table.Columns) {
				out = append(out, col)
			}
			break
		}
	}

	return out
}

func sortedColumns(cols map[string]*ColumnMetadata) []string {
	out := make([]string, 0, len(cols))
	for _, col := range cols {
		out = append(out, col.Name)
	}
	sort.Strings(out)
	return out
}

func (c *Completer) tableCandidates(ctx CompletionContext) []string {
	out := make([]string, 0)
	if ctx.Qualifier != "" {
		if schemaTables, ok := c.metadata.Tables[ctx.Qualifier]; ok {
			for _, table := range schemaTables {
				out = append(out, table.Name)
			}
			sort.Strings(out)
			return out
		}
	}

	for schema, schemaTables := range c.metadata.Tables {
		for _, table := range schemaTables {
			out = append(out, table.Name)
			out = append(out, schema+"."+table.Name)
		}
	}

	return out
}

func (c *Completer) joinCandidates(ctx CompletionContext) []string {
	related := make([]string, 0)
	if ctx.LastTable != nil {
		last := ctx.LastTable
		for schemaName, schemaTables := range c.metadata.Tables {
			for _, t := range schemaTables {
				for _, col := range t.Columns {
					for _, fk := range col.ForeignKey {
						if strings.EqualFold(last.Name, fk.ChildTable) {
							related = append(related, fk.ParentTable)
							related = append(related, schemaName+"."+fk.ParentTable)
						}
						if strings.EqualFold(last.Name, fk.ParentTable) {
							related = append(related, fk.ChildTable)
							related = append(related, schemaName+"."+fk.ChildTable)
						}
					}
				}
			}
		}
	}

	return append(related, c.tableCandidates(ctx)...)
}

func filterAndUniqueByPrefix(candidates []string, prefix string) []string {
	if len(candidates) == 0 {
		return nil
	}

	prefixUpper := strings.ToUpper(prefix)
	seen := make(map[string]bool, len(candidates))
	out := make([]string, 0, len(candidates))

	for _, candidate := range candidates {
		if prefix != "" && !strings.HasPrefix(strings.ToUpper(candidate), prefixUpper) {
			continue
		}
		key := strings.ToUpper(candidate)
		if seen[key] {
			continue
		}
		seen[key] = true
		out = append(out, candidate)
	}

	return out
}
