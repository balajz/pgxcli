package main

import (
	"context"
	"fmt"

	"github.com/balajz/pgxcli/pgxspecial"
	_ "github.com/balajz/pgxcli/pgxspecial/dbcommands"
	"github.com/jackc/pgx/v5/pgxpool"
)

//nolint:gocyclo
func main() {
	// This example demonstrates how to list all databases in a PostgreSQL server using the pgx library.

	ctx := context.Background()
	dbpool, err := pgxpool.New(ctx, "postgres://user:password@localhost:5432/postgres")
	if err != nil {
		panic(err)
	}
	defer dbpool.Close()

	// list databases

	result, ok, err := pgxspecial.ExecuteSpecialCommand(ctx, dbpool, "\\l")
	if err != nil {
		fmt.Println("error occurred: ", err)
		panic(err)
	}
	if !ok {
		panic("command did not execute successfully")
	}
	if result.ResultKind() != pgxspecial.ResultKindRows {
		panic("expected rows result")
	}

	rows := result.(pgxspecial.RowResult).Rows

	columns := rows.FieldDescriptions()
	for _, col := range columns {
		print("| ", col.Name)
	}
	println()

	for rows.Next() {
		val, err := rows.Values()
		if err != nil {
			panic(err)
		}
		fmt.Printf("%s | %s | %s | %s | %s | %s\n",
			val[0], val[1], val[2], val[3], val[4], val[5])
	}

	// Output:
	// | name| owner| encoding| collate| ctype| access_privileges
	// postgres | postgres | UTF8 | en_US.UTF-8 | en_US.UTF-8 | %!s(<nil>)
	// template0 | postgres | UTF8 | en_US.UTF-8 | en_US.UTF-8 | =c/postgres
	// template1 | postgres | UTF8 | en_US.UTF-8 | en_US.UTF-8 | =c/postgres

	// list database with verbose

	verboseResult, ok, err := pgxspecial.ExecuteSpecialCommand(ctx, dbpool, "\\l+")
	if err != nil {
		fmt.Println("error occurred: ", err)
		panic(err)
	}
	if !ok {
		panic("command did not execute successfully")
	}
	if verboseResult.ResultKind() != pgxspecial.ResultKindRows {
		panic("expected rows result")
	}

	rows = verboseResult.(pgxspecial.RowResult).Rows

	columns = rows.FieldDescriptions()
	for _, col := range columns {
		print("| ", col.Name)
	}
	println()

	for rows.Next() {
		val, err := rows.Values()
		if err != nil {
			panic(err)
		}
		fmt.Printf("%s | %s | %s | %s | %s | %s | %s | %s | %s\n",
			val[0], val[1], val[2], val[3], val[4], val[5], val[6], val[7], val[8])
	}

	// Output:
	// | name| owner| encoding| collate| ctype| access_privileges| size| Tablespace| description
	// template0 | postgres | UTF8 | en_US.UTF-8 | en_US.UTF-8 | =c/postgres
	// postgres=CTc/postgres | 7521 kB | pg_default | unmodifiable empty database
	// template1 | postgres | UTF8 | en_US.UTF-8 | en_US.UTF-8 | =c/postgres
	// postgres=CTc/postgres | 7750 kB | pg_default | default template for new databases

	// list database with pattern
	// here the pattern is tem* which will match template0 and template1

	patternResult, ok, err := pgxspecial.ExecuteSpecialCommand(ctx, dbpool, "\\l tem*")
	if err != nil {
		fmt.Println("error occurred: ", err)
		panic(err)
	}
	if !ok {
		panic("command did not execute successfully")
	}
	if patternResult.ResultKind() != pgxspecial.ResultKindRows {
		panic("expected rows result")
	}

	rows = patternResult.(pgxspecial.RowResult).Rows

	columns = rows.FieldDescriptions()
	for _, col := range columns {
		print("| ", col.Name)
	}
	println()

	for rows.Next() {
		val, err := rows.Values()
		if err != nil {
			panic(err)
		}
		fmt.Printf("%s | %s | %s | %s | %s | %s\n",
			val[0], val[1], val[2], val[3], val[4], val[5])
	}

	// Output:
	// | name| owner| encoding| collate| ctype| access_privileges
	// template0 | postgres | UTF8 | en_US.UTF-8 | en_US.UTF-8 | =c/postgres
	// template1 | postgres | UTF8 | en_US.UTF-8 | en_US.UTF-8 | =c/postgres
}
