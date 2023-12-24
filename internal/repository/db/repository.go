package db

import (
	"ara-server/internal/infrastructure"

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
