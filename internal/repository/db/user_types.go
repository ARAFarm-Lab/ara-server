package db

import (
	"ara-server/internal/constants"
	"database/sql"
	"time"
)

type InsertUserParam struct {
	Email    string `db:"email"`
	Password string `db:"password"`
}

type InsertProfileParam struct {
	UserID int                `db:"user_id"`
	Name   string             `db:"name"`
	Role   constants.UserRole `db:"role"`
}

type User struct {
	ID        int          `db:"id"`
	Email     string       `db:"email"`
	Password  string       `db:"password"`
	IsActive  bool         `db:"is_active"`
	CreatedAt time.Time    `db:"created_at"`
	UpdatedAt sql.NullTime `db:"updated_at"`
}

type UserInfo struct {
	UserID   int                `db:"user_id"`
	IsActive bool               `db:"is_active"`
	Name     string             `db:"name"`
	Role     constants.UserRole `db:"role"`
}
