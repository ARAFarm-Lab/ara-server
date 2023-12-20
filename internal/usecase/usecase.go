package usecase

import (
	"ara-server/internal/infrastructure"
	"ara-server/internal/repository/db"
	"ara-server/internal/repository/mq"
)

type Usecase struct {
	infra *infrastructure.Infrastructure
	db    *db.Repository
	mq    *mq.Repository
}

func NewUsecase(infra *infrastructure.Infrastructure, db *db.Repository, mq *mq.Repository) *Usecase {
	return &Usecase{infra: infra, db: db, mq: mq}
}
