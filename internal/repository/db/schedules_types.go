package db

import (
	"database/sql"
	"time"
)

type ActionSchedule struct {
	ID            int            `db:"id"`
	Name          string         `db:"name"`
	Description   sql.NullString `db:"description"`
	Actions       string         `db:"actions"`
	Schedule      sql.NullString `db:"schedule"`
	IsActive      bool           `db:"is_active"`
	NextRunAt     time.Time      `db:"next_run_at"`
	LastLockAt    sql.NullTime   `db:"last_lock_at"`
	LastRunAt     sql.NullTime   `db:"last_run_at"`
	LastRunStatus sql.NullInt32  `db:"last_run_status"`
	LastError     sql.NullString `db:"last_error"`
	Attempts      int            `db:"attempts"`
	CreatedAt     time.Time      `db:"created_at"`
}
