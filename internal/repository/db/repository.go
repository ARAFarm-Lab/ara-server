package db

import (
	"ara-server/internal/infrastructure"
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type Repository struct {
	db    *sqlx.DB
	infra *infrastructure.Infrastructure
}

func NewRepository(db *sqlx.DB, infra *infrastructure.Infrastructure) *Repository {
	return &Repository{
		db:    db,
		infra: infra,
	}
}

func (repo *Repository) BeginTx(ctx context.Context) (*sql.Tx, error) {
	return repo.db.BeginTx(ctx, nil)
}

func (repo *Repository) Rebind(query string) string {
	return sqlx.Rebind(sqlx.DOLLAR, query)
}
