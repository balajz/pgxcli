package database

import (
	"context"
	"database/sql/driver"
)

type QueryFn func(ctx context.Context, query string, isMulti bool) (Rows, bool, error)

// Rows describes a result set.
type Rows interface {
	driver.Rows

	// The caller must call Close() when done with the
	// result and check the error.
	Close() error

	// Columns returns the column labels of the current result set.
	// The implementation of this method should cache the result so that the
	// result does not need to be constructed on each invocation.
	Columns() []string

	// ColumnTypeDatabaseTypeName returns the database type name
	// of the column at the given column index.
	ColumnTypeDatabaseTypeName(index int) string

	// Tag retrieves the statement tag for the current result set.
	Tag() (CommandTag, error)

	// Next populates values with the next row of results. []byte values are copied
	// so that subsequent calls to Next and Close do not mutate values. This
	// makes it slower than theoretically possible but the safety concerns
	// (since this is unobvious and unexpected behavior) outweigh.
	Next(values []driver.Value) error

	// NextResultSet prepares the next result set for reading.
	// Returns false if there is no more result set to read.
	NextResultSet() (bool, error)
}

// CommandTag represents the result of a SQL command.
type CommandTag interface {
	RowsAffected() int64
	String() string
}
