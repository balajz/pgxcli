package completer

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestComplete_SelectAliasDotSuggestsColumns(t *testing.T) {
	c := setupTestCompleter(t)
	category, candidates := c.Complete("SELECT u. FROM users u", len("SELECT u."))

	require.Equal(t, "Columns", category)
	assert.Contains(t, candidates, "id")
	assert.Contains(t, candidates, "name")
}

func TestComplete_FromSuggestsTables(t *testing.T) {
	c := setupTestCompleter(t)
	category, candidates := c.Complete("SELECT * FROM ", len("SELECT * FROM "))

	require.Equal(t, "Tables", category)
	assert.Contains(t, candidates, "users")
	assert.Contains(t, candidates, "orders")
}

func TestComplete_WhereSuggestsReferencedColumns(t *testing.T) {
	c := setupTestCompleter(t)
	query := "SELECT * FROM users u WHERE "
	category, candidates := c.Complete(query, len(query))

	require.Equal(t, "Columns", category)
	assert.Contains(t, candidates, "u.id")
	assert.Contains(t, candidates, "u.name")
}

func TestComplete_JoinSuggestsForeignKeyRelatedTable(t *testing.T) {
	c := setupTestCompleter(t)
	query := "SELECT * FROM users u JOIN "
	category, candidates := c.Complete(query, len(query))

	require.Equal(t, "Tables", category)
	assert.Contains(t, candidates, "orders")
}

func setupTestCompleter(t *testing.T) *Completer {
	t.Helper()

	c := New(nil)
	c.ExtendSchemas([]string{"public"})
	c.ExtendTables([]Relation{
		{Schema: "public", Name: "users", Kind: RelationKindTable},
		{Schema: "public", Name: "orders", Kind: RelationKindTable},
	})
	c.ExtendColumns([]ColumnInfo{
		{Schema: "public", Table: "users", Column: "id"},
		{Schema: "public", Table: "users", Column: "name"},
		{Schema: "public", Table: "orders", Column: "id"},
		{Schema: "public", Table: "orders", Column: "user_id"},
	}, false)
	c.ExtendForeignKeys([]ForeignKey{
		{
			ParentSchema: "public",
			ParentTable:  "users",
			ParentColumn: "id",
			ChildSchema:  "public",
			ChildTable:   "orders",
			ChildColumn:  "user_id",
		},
	})

	return c
}
