package db

import "ara-server/internal/constants"

type Actuator struct {
	ID             int64                `db:"id"`
	DeviceID       int64                `db:"device_id"`
	PinNumber      int                  `db:"pin_number"`
	ActionType     constants.ActionType `db:"action_type"`
	TerminalNumber int                  `db:"terminal_number"`
	Name           string               `db:"name"`
	Icon           string               `db:"icon"`
	IsActive       bool                 `db:"is_active"`
}
