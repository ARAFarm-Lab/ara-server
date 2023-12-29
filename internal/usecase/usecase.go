package usecase

import (
	"ara-server/internal/constants"
	"ara-server/internal/infrastructure"
	"ara-server/internal/repository/db"
	"ara-server/internal/repository/mq"
	"errors"
	"strconv"
)

var (
	errorInvalidActionValue = errors.New("invalid action value")
)

type Usecase struct {
	infra *infrastructure.Infrastructure
	db    *db.Repository
	mq    *mq.Repository

	actionDispatcherMap map[constants.ActionType]func(DispatcherParam) error
}

func NewUsecase(infra *infrastructure.Infrastructure, db *db.Repository, mq *mq.Repository) *Usecase {
	uc := &Usecase{infra: infra, db: db, mq: mq}

	// define action dispatcher map
	uc.actionDispatcherMap = map[constants.ActionType]func(DispatcherParam) error{
		constants.ActionTypeBuiltInLED: uc.toggleBuiltInLED,
		constants.ActionTypeRelay:      uc.toggleRelay,
	}

	return uc
}

func generateDeviceTopic(deviceID int64) string {
	return "d-" + strconv.FormatInt(deviceID, 10)
}
