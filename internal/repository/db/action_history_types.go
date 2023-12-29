package db

import (
	"ara-server/internal/constants"
	"time"
)

type ActionHistory struct {
	DeviceID   int64                  `db:"device_id"`
	ActionType constants.ActionType   `db:"action_type"`
	Value      interface{}            `db:"value"`
	ActionBy   constants.ActionSource `db:"action_by"`
	ActionAt   time.Time              `db:"action_at"`
}
