package database

import (
	"time"
)

// ExecResult represents a non-row SQL command result (INSERT, UPDATE, DELETE, DDL).
type ExecResult struct {
	RowsAffected int64
	Status       string
	Duration     time.Duration
}

func (e *ExecResult) isResult() {}
