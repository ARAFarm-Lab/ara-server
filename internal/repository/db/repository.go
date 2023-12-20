package db

import "ara-server/internal/infrastructure"

type Repository struct {
	infra *infrastructure.Infrastructure
}

func NewRepository(infra *infrastructure.Infrastructure) *Repository {
	return &Repository{infra: infra}
}
