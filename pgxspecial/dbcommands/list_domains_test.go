package dbcommands_test

import (
	"context"
	"testing"

	"github.com/balajz/pgxcli/pgxspecial/dbcommands"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
)

func CreateDomain(
	t *testing.T,
	ctx context.Context,
	pool *pgxpool.Pool,
	domain string,
	baseType string,
) {
	t.Helper()

	sql := `CREATE DOMAIN ` + domain + ` AS ` + baseType

	if _, err := pool.Exec(ctx, sql); err != nil {
		t.Fatalf("create domain %q failed: %v", domain, err)
	}
}

func DropDomain(
	t *testing.T,
	ctx context.Context,
	pool *pgxpool.Pool,
	domain string,
) {
	t.Helper()

	sql := `DROP DOMAIN IF EXISTS ` + domain + ` CASCADE`

	if _, err := pool.Exec(ctx, sql); err != nil {
		t.Fatalf("drop domain %q failed: %v", domain, err)
	}
}

func TestListDomains(t *testing.T) {
	db := connectTestDB(t)
	defer db.(*pgxpool.Pool).Close()

	pattern := ""
	verbose := false

	// Setup: Create domains
	ctx := context.Background()
	domains := []struct {
		name     string
		baseType string
	}{
		{"mood_domain", "TEXT"},
		{"age_domain", "INT"},
		{"email_domain", "VARCHAR(255)"},
	}

	for _, domain := range domains {
		CreateDomain(t, ctx, db.(*pgxpool.Pool), domain.name, domain.baseType)
		defer DropDomain(t, ctx, db.(*pgxpool.Pool), domain.name)
	}
	res, err := dbcommands.ListDomains(context.Background(), db, pattern, verbose)
	if err != nil {
		t.Fatalf("ListDomains failed: %v", err)
	}
	result := RequiresRowResult(t, res)

	defer result.Rows.Close()

	fds := result.Rows.FieldDescriptions()

	expectedColumns := []string{
		"schema",
		"name",
		"type",
		"modifier",
		"check",
	}
	assert.Equal(t, expectedColumns, getColumnNames(fds), "Column names do not match expected")
	// expecting 5 columns
	assert.Len(t, fds, 5)
	var allRows []map[string]interface{}
	allRows, err = RowsToMaps(result.Rows)
	if err != nil {
		t.Fatalf("Failed to read rows: %v", err)
	}
	assert.GreaterOrEqual(t, len(allRows), 3, "Expected at least 3 domains in the result")
	assert.True(t, containsByField(allRows, "name", "mood_domain"))
	assert.True(t, containsByField(allRows, "name", "age_domain"))
	assert.True(t, containsByField(allRows, "name", "email_domain"))
}

func TestListDomainsWithPattern(t *testing.T) {
	db := connectTestDB(t)
	defer db.(*pgxpool.Pool).Close()

	pattern := "*_domain"
	verbose := false

	// Setup: Create domains
	ctx := context.Background()
	domains := []struct {
		name     string
		baseType string
	}{
		{"mood_domain", "TEXT"},
		{"age_domain", "INT"},
		{"email_domain", "VARCHAR(255)"},
	}

	for _, domain := range domains {
		CreateDomain(t, ctx, db.(*pgxpool.Pool), domain.name, domain.baseType)
		defer DropDomain(t, ctx, db.(*pgxpool.Pool), domain.name)
	}
	res, err := dbcommands.ListDomains(context.Background(), db, pattern, verbose)
	if err != nil {
		t.Fatalf("ListDomains failed: %v", err)
	}
	result := RequiresRowResult(t, res)

	defer result.Rows.Close()

	fds := result.Rows.FieldDescriptions()

	expectedColumns := []string{
		"schema",
		"name",
		"type",
		"modifier",
		"check",
	}
	assert.Equal(t, expectedColumns, getColumnNames(fds), "Column names do not match expected")
	// expecting 5 columns
	assert.Len(t, fds, 5)
	var allRows []map[string]interface{}
	allRows, err = RowsToMaps(result.Rows)
	if err != nil {
		t.Fatalf("Failed to read rows: %v", err)
	}
	assert.GreaterOrEqual(t, len(allRows), 3, "Expected at least 3 domains in the result")
	assert.True(t, containsByField(allRows, "name", "mood_domain"))
	assert.True(t, containsByField(allRows, "name", "age_domain"))
	assert.True(t, containsByField(allRows, "name", "email_domain"))
}

func TestListDomainsWithNoMatchingPattern(t *testing.T) {
	db := connectTestDB(t)
	defer db.(*pgxpool.Pool).Close()

	pattern := "nonexistent_domain"
	verbose := false

	// Setup: Create domains
	ctx := context.Background()
	domains := []struct {
		name     string
		baseType string
	}{
		{"mood_domain", "TEXT"},
		{"age_domain", "INT"},
		{"email_domain", "VARCHAR(255)"},
	}

	for _, domain := range domains {
		CreateDomain(t, ctx, db.(*pgxpool.Pool), domain.name, domain.baseType)
		defer DropDomain(t, ctx, db.(*pgxpool.Pool), domain.name)
	}
	res, err := dbcommands.ListDomains(context.Background(), db, pattern, verbose)
	if err != nil {
		t.Fatalf("ListDomains failed: %v", err)
	}
	result := RequiresRowResult(t, res)

	defer result.Rows.Close()

	fds := result.Rows.FieldDescriptions()

	expectedColumns := []string{
		"schema",
		"name",
		"type",
		"modifier",
		"check",
	}
	assert.Equal(t, expectedColumns, getColumnNames(fds), "Column names do not match expected")
	// expecting 5 columns
	assert.Len(t, fds, 5)
	var allRows []map[string]interface{}
	allRows, err = RowsToMaps(result.Rows)
	if err != nil {
		t.Fatalf("Failed to read rows: %v", err)
	}
	assert.Len(t, allRows, 0, "Expected no domains in the result")
}

func TestListDomainsVerbose(t *testing.T) {
	db := connectTestDB(t)
	defer db.(*pgxpool.Pool).Close()

	pattern := ""
	verbose := true

	// Setup: Create domains
	ctx := context.Background()
	domains := []struct {
		name     string
		baseType string
	}{
		{"mood_domain", "TEXT"},
		{"age_domain", "INT"},
		{"email_domain", "VARCHAR(255)"},
	}

	for _, domain := range domains {
		CreateDomain(t, ctx, db.(*pgxpool.Pool), domain.name, domain.baseType)
		defer DropDomain(t, ctx, db.(*pgxpool.Pool), domain.name)
	}
	res, err := dbcommands.ListDomains(context.Background(), db, pattern, verbose)
	if err != nil {
		t.Fatalf("ListDomains failed: %v", err)
	}
	result := RequiresRowResult(t, res)

	defer result.Rows.Close()

	fds := result.Rows.FieldDescriptions()

	expectedColumns := []string{
		"schema",
		"name",
		"type",
		"modifier",
		"check",
		"access_privileges",
		"description",
	}
	assert.Equal(t, expectedColumns, getColumnNames(fds), "Column names do not match expected")
	// expecting 7 columns
	assert.Len(t, fds, 7)
	var allRows []map[string]interface{}
	allRows, err = RowsToMaps(result.Rows)
	if err != nil {
		t.Fatalf("Failed to read rows: %v", err)
	}
	assert.GreaterOrEqual(t, len(allRows), 3, "Expected at least 3 domains in the result")
	assert.True(t, containsByField(allRows, "name", "mood_domain"))
	assert.True(t, containsByField(allRows, "name", "age_domain"))
	assert.True(t, containsByField(allRows, "name", "email_domain"))
}

func TestListDomainsVerboseWithPattern(t *testing.T) {
	db := connectTestDB(t)
	defer db.(*pgxpool.Pool).Close()

	pattern := "*_domain"
	verbose := true

	// Setup: Create domains
	ctx := context.Background()
	domains := []struct {
		name     string
		baseType string
	}{
		{"mood_domain", "TEXT"},
		{"age_domain", "INT"},
		{"email_domain", "VARCHAR(255)"},
	}

	for _, domain := range domains {
		CreateDomain(t, ctx, db.(*pgxpool.Pool), domain.name, domain.baseType)
		defer DropDomain(t, ctx, db.(*pgxpool.Pool), domain.name)
	}
	res, err := dbcommands.ListDomains(context.Background(), db, pattern, verbose)
	if err != nil {
		t.Fatalf("ListDomains failed: %v", err)
	}
	result := RequiresRowResult(t, res)

	defer result.Rows.Close()

	fds := result.Rows.FieldDescriptions()

	expectedColumns := []string{
		"schema",
		"name",
		"type",
		"modifier",
		"check",
		"access_privileges",
		"description",
	}
	assert.Equal(t, expectedColumns, getColumnNames(fds), "Column names do not match expected")
	// expecting 7 columns
	assert.Len(t, fds, 7)
	var allRows []map[string]interface{}
	allRows, err = RowsToMaps(result.Rows)
	if err != nil {
		t.Fatalf("Failed to read rows: %v", err)
	}
	assert.GreaterOrEqual(t, len(allRows), 3, "Expected at least 3 domains in the result")
	assert.True(t, containsByField(allRows, "name", "mood_domain"))
	assert.True(t, containsByField(allRows, "name", "age_domain"))
	assert.True(t, containsByField(allRows, "name", "email_domain"))
}

func TestListDomainsVerboseWithNoMatchingPattern(t *testing.T) {
	db := connectTestDB(t)
	defer db.(*pgxpool.Pool).Close()

	pattern := "nonexistent_domain"
	verbose := true

	// Setup: Create domains
	ctx := context.Background()
	domains := []struct {
		name     string
		baseType string
	}{
		{"mood_domain", "TEXT"},
		{"age_domain", "INT"},
		{"email_domain", "VARCHAR(255)"},
	}

	for _, domain := range domains {
		CreateDomain(t, ctx, db.(*pgxpool.Pool), domain.name, domain.baseType)
		defer DropDomain(t, ctx, db.(*pgxpool.Pool), domain.name)
	}
	res, err := dbcommands.ListDomains(context.Background(), db, pattern, verbose)
	if err != nil {
		t.Fatalf("ListDomains failed: %v", err)
	}
	result := RequiresRowResult(t, res)

	defer result.Rows.Close()

	fds := result.Rows.FieldDescriptions()

	expectedColumns := []string{
		"schema",
		"name",
		"type",
		"modifier",
		"check",
		"access_privileges",
		"description",
	}
	assert.Equal(t, expectedColumns, getColumnNames(fds), "Column names do not match expected")
	// expecting 7 columns
	assert.Len(t, fds, 7)
	var allRows []map[string]interface{}
	allRows, err = RowsToMaps(result.Rows)
	if err != nil {
		t.Fatalf("Failed to read rows: %v", err)
	}
	assert.Len(t, allRows, 0, "Expected no domains in the result")
}

func TestListDomainsWithSchema(t *testing.T) {
	db := connectTestDB(t)
	defer db.(*pgxpool.Pool).Close()

	testSchema := "custom_schema"

	CreateSchema(t, context.Background(), db.(*pgxpool.Pool), testSchema)
	defer DropSchema(t, context.Background(), db.(*pgxpool.Pool), testSchema)

	pattern := testSchema + ".*"
	verbose := false

	// Setup: Create domains
	ctx := context.Background()
	domains := []struct {
		name     string
		baseType string
	}{
		{testSchema + ".mood_domain", "TEXT"},
		{testSchema + ".age_domain", "INT"},
		{testSchema + ".email_domain", "VARCHAR(255)"},
	}

	for _, domain := range domains {
		CreateDomain(t, ctx, db.(*pgxpool.Pool), domain.name, domain.baseType)
		defer DropDomain(t, ctx, db.(*pgxpool.Pool), domain.name)
	}
	res, err := dbcommands.ListDomains(context.Background(), db, pattern, verbose)
	if err != nil {
		t.Fatalf("ListDomains failed: %v", err)
	}
	result := RequiresRowResult(t, res)

	defer result.Rows.Close()

	fds := result.Rows.FieldDescriptions()

	expectedColumns := []string{
		"schema",
		"name",
		"type",
		"modifier",
		"check",
	}
	assert.Equal(t, expectedColumns, getColumnNames(fds), "Column names do not match expected")
	// expecting 5 columns
	assert.Len(t, fds, 5)
	var allRows []map[string]interface{}
	allRows, err = RowsToMaps(result.Rows)
	if err != nil {
		t.Fatalf("Failed to read rows: %v", err)
	}
	assert.Len(t, allRows, 3, "Expected exactly 3 domains in the result")
	assert.True(t, containsByField(allRows, "name", "mood_domain"))
}
