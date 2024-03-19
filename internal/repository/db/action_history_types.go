package db

import (
	"ara-server/internal/constants"
	"time"
)

type ActionHistory struct {
	ActuatorID int64                  `db:"actuator_id"`
	ActionType constants.ActionType   `db:"action_type"`
	Name       string                 `db:"name"`
	Icon       string                 `db:"icon"`
	Value      interface{}            `db:"value"`
	ActionBy   constants.ActionSource `db:"action_by"`
	ActionAt   time.Time              `db:"action_at"`
}
