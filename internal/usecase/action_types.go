package usecase

import (
	"ara-server/internal/constants"
	"time"
)

type ActionHistory struct {
	ActuatorID int64                  `json:"actuator_id,omitempty"`
	Value      interface{}            `json:"value,omitempty"`
	ActionBy   constants.ActionSource `json:"action_by,omitempty"`
	ActionAt   *time.Time             `json:"action_at,omitempty"`
	Action     DispatcherAction       `json:"action,omitempty"`
}

type DispatcherAction struct {
	ID   int64                `json:"id,omitempty"`
	Type constants.ActionType `json:"type,omitempty"`
	Name string               `json:"name,omitempty"`
	Icon string               `json:"icon,omitempty"`
}

type DispatcherParam struct {
	DeviceID   int64                  `json:"device_id,omitempty"`
	ActuatorID int64                  `json:"actuator_id,omitempty"`
	Value      interface{}            `json:"value,omitempty"`
	ActionBy   constants.ActionSource `json:"action_by,omitempty"`
}

type InsertActionLogParam struct {
	ActuatorID int64
	Value      interface{}
	ActionBy   constants.ActionSource
	ActionAt   time.Time
}
