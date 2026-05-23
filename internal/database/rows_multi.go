package database

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"io"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
)

type sqlRowsMultiResultSet struct {
	rows     *pgconn.MultiResultReader
	typeMap  *pgtype.Map
	exec     *executor
	colNames []string
}

var _ Rows = (*sqlRowsMultiResultSet)(nil)

func (r *sqlRowsMultiResultSet) Columns() []string {
	if r.colNames == nil {
		rr := r.rows.ResultReader()
		if rr == nil {
			// ResultReader may be nil if an empty query was executed.
			return nil
		}
		fields := rr.FieldDescriptions()
		r.colNames = make([]string, len(fields))
		for i, fd := range fields {
			r.colNames[i] = fd.Name
		}
	}
	return r.colNames
}

func (r *sqlRowsMultiResultSet) Tag() (CommandTag, error) {
	if rr := r.rows.ResultReader(); rr != nil {
		// ResultReader may be nil if an empty query was executed.
		return r.rows.ResultReader().Close()
	}
	return pgconn.CommandTag{}, nil
}

func (r *sqlRowsMultiResultSet) Close() (retErr error) {
	if rr := r.rows.ResultReader(); rr != nil {
		// ResultReader may be nil if an empty query was executed.
		_, retErr = r.rows.ResultReader().Close()
	}
	retErr = errors.Join(retErr, r.rows.Close())
	if r.exec.conn.IsClosed() {
		return retErr
	}
	return retErr
}

// Next implements the Rows interface.
func (r *sqlRowsMultiResultSet) Next(values []driver.Value) error {
	if r.exec.conn.IsClosed() {
		return ErrConnectionClosed
	}
	rd := r.rows.ResultReader()
	if rd == nil {
		// ResultReader may be nil if an empty query was executed.
		return io.EOF
	}
	if !rd.NextRow() {
		if _, err := rd.Close(); err != nil {
			return err
		}
		return io.EOF
	}
	if len(rd.FieldDescriptions()) != len(values) {
		return fmt.Errorf("number of field descriptions must equal number of destinations, got %d and %d",
			len(rd.FieldDescriptions()), len(values))
	}
	for i := range values {
		rowVal := rd.Values()[i]
		if rowVal == nil {
			values[i] = nil
			continue
		}
		fieldOID := rd.FieldDescriptions()[i].DataTypeOID
		fieldFormat := rd.FieldDescriptions()[i].Format

		// By scanning the value into a string, pgconn will use the pgwire
		// text format to represent the value.
		var s string
		err := r.typeMap.Scan(fieldOID, fieldFormat, rowVal, &s)
		if err != nil {
			return pgx.ScanArgError{ColumnIndex: i, Err: err}
		}
		values[i] = s
	}
	return nil
}

// NextResultSet prepares the next result set for reading.
func (r *sqlRowsMultiResultSet) NextResultSet() (bool, error) {
	r.colNames = nil
	next := r.rows.NextResult()
	if !next {
		if err := r.rows.Close(); err != nil {
			return false, err
		}
	}
	if r.exec.conn.IsClosed() {
		return false, ErrConnectionClosed
	}
	return next, nil
}

func (r *sqlRowsMultiResultSet) ColumnTypeDatabaseTypeName(index int) string {
	rd := r.rows.ResultReader()
	fieldOID := rd.FieldDescriptions()[index].DataTypeOID
	return databaseTypeName(r.typeMap, fieldOID)
}
