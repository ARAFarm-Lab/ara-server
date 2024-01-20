package usecase

import (
	"context"
	"database/sql"
	"errors"
	"strconv"
	"time"

	"github.com/robfig/cron/v3"

	"ara-server/internal/constants"
	"ara-server/internal/infrastructure"
	"ara-server/internal/repository/db"
	"ara-server/internal/repository/mq"
)

var (
	errorDispatchSchedulerFailed = errors.New("dispatch scheduler failed")
	errorInvalidActionType       = errors.New("invalid action type")
	errorInvalidActionValue      = errors.New("invalid action value")
)

type Usecase struct {
	infra *infrastructure.Infrastructure
	db    *db.Repository
	mq    *mq.Repository

	cronParser cron.Parser

	actionDispatcherMap map[constants.ActionType]func(context.Context, DispatcherParam) error
}

func NewUsecase(infra *infrastructure.Infrastructure, db *db.Repository, mq *mq.Repository) *Usecase {
	uc := &Usecase{
		infra:      infra,
		db:         db,
		mq:         mq,
		cronParser: cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow),
	}

	// define action dispatcher map
	uc.actionDispatcherMap = map[constants.ActionType]func(context.Context, DispatcherParam) error{
		constants.ActionTypeBuiltInLED: uc.toggleBuiltInLED,
		constants.ActionTypeRelay:      uc.toggleRelay,
	}

	return uc
}

func assignSQLNullTime(time *time.Time) sql.NullTime {
	if time == nil {
		return sql.NullTime{}
	}

	return sql.NullTime{
		Time:  *time,
		Valid: true,
	}
}

func assignSQLNullInt(value int) sql.NullInt32 {
	if value == 0 {
		return sql.NullInt32{}
	}

	return sql.NullInt32{
		Int32: int32(value),
		Valid: true,
	}
}

func assignSQLNullString(value string) sql.NullString {
	if value == "" {
		return sql.NullString{}
	}

	return sql.NullString{
		String: value,
		Valid:  true,
	}
}

func generateDeviceTopic(deviceID int64) string {
	return "d-" + strconv.FormatInt(deviceID, 10)
}

func getSQLNullString(str sql.NullString) string {
	if !str.Valid {
		return ""
	}

	return str.String
}
