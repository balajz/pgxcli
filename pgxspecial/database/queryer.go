// this package provide the interface to interact with the database
package database

import (
	"context"

	"github.com/jackc/pgx/v5"
)

// Queryer is an interface that defines methods for querying a database.
// It includes methods for executing queries that return multiple rows, single rows, and closing the connection.
// test with pgx.Conn and pgxpool.Pool.
// suitable pgx ecosystem
type Queryer interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}
