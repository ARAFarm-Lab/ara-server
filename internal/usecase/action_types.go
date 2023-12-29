package usecase

import (
	"ara-server/internal/constants"
	"time"
)

type ActionHistory struct {
	DeviceID   int64                  `json:"device_id,omitempty"`
	ActionType constants.ActionType   `json:"action_type,omitempty"`
	Value      interface{}            `json:"value,omitempty"`
	ActionBy   constants.ActionSource `json:"action_by,omitempty"`
	ActionAt   *time.Time             `json:"action_at,omitempty"`
}

type DispatcherParam struct {
	DeviceID   int64
	ActionType constants.ActionType
	Value      interface{}
}

type InsertActionLogParam struct {
	DeviceID   int64
	ActionType constants.ActionType
	Value      interface{}
	ActionBy   constants.ActionSource
	ActionAt   time.Time
}
