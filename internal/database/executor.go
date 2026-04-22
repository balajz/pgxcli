package database

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/balaji01-4d/pgxcli/internal/parser"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// Result marks values returned by SQL execution paths.
type Result interface {
	isResult()
}

type conn interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
	Config() *pgx.ConnConfig
	Ping(ctx context.Context) error
	Close(ctx context.Context) error
}

type executor struct {
	Host     string
	Port     uint16
	Database string
	Schema   string
	User     string
	Password string
	URI      string
	Conn     conn

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
		Host:     conn.Config().Host,
		Port:     conn.Config().Port,
		Database: conn.Config().Database,
		User:     conn.Config().User,
		Password: conn.Config().Password,
		URI:      conn.Config().ConnString(),
		Conn:     conn,
		Logger:   logger,
	}, nil
}

// For executing queries like SELECT, SHOW etc.
func (e *executor) query(ctx context.Context, sql string, args ...any) (*QueryResult, error) {
	e.Logger.Debug("Executing query", "sql", sql)
	start := time.Now()
	rows, err := e.Conn.Query(ctx, sql, args...)
	if err != nil {
		e.Logger.Error("Query failed", "error", err, "sql", sql)
		return nil, err
	}
	dur := time.Since(start)
	fds := rows.FieldDescriptions()
	columns := make([]string, len(fds))
	for i, fd := range fds {
		columns[i] = fd.Name
	}
	e.Logger.Info("Query completed", "duration_ms", dur.Milliseconds(), "columns", len(columns))
	return &QueryResult{
		rowStreamer: rowStreamer{
			rows:     rows,
			columns:  columns,
			duration: dur,
		},
	}, nil
}

// For executing commands like INSERT, UPDATE, DELETE etc.
func (e *executor) exec(ctx context.Context, sql string, args ...any) (*ExecResult, error) {
	e.Logger.Debug("Executing command", "sql", sql)
	start := time.Now()
	tag, err := e.Conn.Exec(ctx, sql, args...)
	if err != nil {
		e.Logger.Error("Command failed", "error", err, "sql", sql)
		return nil, err
	}
	dur := time.Since(start)
	e.Logger.Info("Command completed", "duration_ms", dur.Milliseconds(), "rows_affected", tag.RowsAffected(), "status", tag.String())
	return &ExecResult{
		RowsAffected: tag.RowsAffected(),
		Status:       tag.String(),
		Duration:     dur,
	}, nil
}

// execute determines whether to run query or exec based on SQL type.
func (e *executor) execute(ctx context.Context, sql string, args ...any) (Result, error) {
	if parser.IsQuery(sql) {
		return e.query(ctx, sql, args...)
	}
	return e.exec(ctx, sql, args...)
}

func (e *executor) close(ctx context.Context) error {
	if e.Conn != nil {
		return e.Conn.Close(ctx)
	}
	return nil
}

func (e *executor) ping(ctx context.Context) error {
	if e.Conn == nil {
		return fmt.Errorf("database not connected")
	}
	return e.Conn.Ping(ctx)
}

func (e *executor) isConnected() bool {
	return e.Conn != nil
}
