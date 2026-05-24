package database

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/balajz/pgxcli/internal/database/result"
	"github.com/balajz/pgxcli/pgxspecial"

	"github.com/jackc/pgx/v5"
)

var (
	ErrConnectionClosed         = errors.New("connection closed unexpectedly")
	ErrConnectionNotEstablished = errors.New("connection not established")
)

type executor struct {
	Host      string
	Port      uint16
	Database  string
	Schema    string
	User      string
	Password  string
	URI       string
	conn      *pgx.Conn
	Connector Connector

	Logger *slog.Logger
}

func newExecutor(ctx context.Context, c Connector, logger *slog.Logger) (*executor, error) {
	conn, err := c.Connect(ctx)
	if err != nil {
		return nil, err
	}

	err = conn.Ping(ctx)
	if err != nil {
		logger.Error("Connection ping failed", "error", err)
		if err := conn.Close(ctx); err != nil {
			logger.Error("Failed to close connection", "error", err)
		}

		return nil, err
	}

	return &executor{
		Host:      conn.Config().Host,
		Port:      conn.Config().Port,
		Database:  conn.Config().Database,
		User:      conn.Config().User,
		Password:  conn.Config().Password,
		URI:       conn.Config().ConnString(),
		conn:      conn,
		Connector: c,
		Logger:    logger,
	}, nil
}

func (e *executor) ensureConn(ctx context.Context) error {
	if e.conn != nil && !e.conn.IsClosed() {
		return nil
	}

	conn, err := e.Connector.Connect(ctx)
	if err != nil {
		return err
	}
	e.conn = conn
	return nil
}

// For executing queries like SELECT, SHOW etc.
func (e *executor) query(ctx context.Context, sql string, isMulti bool, args ...any) (Rows, bool, error) {
	if err := e.ensureConn(ctx); err != nil {
		return nil, false, err
	}
	e.Logger.Debug("Executing query", "sql", sql)
	if isMulti {
		mrr := e.conn.PgConn().Exec(ctx, sql)
		if e.conn.IsClosed() {
			mrr.Close()
			return nil, false, ErrConnectionClosed
		}

		rs := &sqlRowsMultiResultSet{
			rows:    mrr,
			typeMap: e.conn.TypeMap(),
			exec:    e,
		}
		return rs, true, nil
	}

	rows, err := e.conn.Query(ctx, sql, args...)
	if err != nil {
		if e.conn.IsClosed() {
			return nil, false, ErrConnectionClosed
		}
		e.Logger.Error("Query failed", "error", err, "sql", sql)
		return nil, false, err
	}

	return &sqlRows{
		rows:    rows,
		typeMap: e.conn.TypeMap(),
		exec:    e,
	}, false, nil
}

// execute determines whether to run query or exec based on SQL type.
func (e *executor) execute(ctx context.Context, sql string, isMulti bool, args ...any) (Rows, bool, error) {
	return e.query(ctx, sql, isMulti, args...)
}

func (e *executor) executeSpecial(ctx context.Context, cmd string) (pgxspecial.SpecialCommandResult, bool, error) {
	if err := e.ensureConn(ctx); err != nil {
		return nil, false, err
	}
	specialResult, ok, err := pgxspecial.ExecuteSpecialCommand(ctx, e.conn, cmd)
	if err != nil {
		e.Logger.Error("Special command execution failed", "error", err, "command", cmd)
		return nil, ok, err
	}

	if !ok || specialResult == nil || specialResult.ResultKind() != pgxspecial.ResultKindRows {
		return specialResult, ok, nil
	}

	rowResult, isRowResult := specialResult.(pgxspecial.RowResult)
	if !isRowResult {
		return nil, ok, fmt.Errorf("invalid row special result type: %T", specialResult)
	}

	normalizedRows, err := result.NewSpecialRow(rowResult.Rows)
	if err != nil {
		e.Logger.Error("Failed to materialize special command rows", "error", err, "command", cmd)
		return nil, ok, err
	}

	return normalizedRows, ok, nil
}

func (e *executor) cancel(ctx context.Context) error {
	if e.conn == nil || e.conn.IsClosed() {
		return ErrConnectionNotEstablished
	}
	return e.conn.PgConn().CancelRequest(ctx)
}

func (e *executor) close(ctx context.Context) error {
	if e.conn != nil && !e.conn.IsClosed() {
		return e.conn.Close(ctx)
	}
	return nil
}

func (e *executor) ping(ctx context.Context) error {
	if err := e.ensureConn(ctx); err != nil {
		return err
	}
	return e.conn.Ping(ctx)
}

func (e *executor) isConnected() bool {
	return e.conn != nil && !e.conn.IsClosed()
}
