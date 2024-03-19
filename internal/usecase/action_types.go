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
	ID   int64                `json:"id"`
	Type constants.ActionType `json:"type"`
	Name string               `json:"name"`
	Icon string               `json:"icon"`
}

type DispatcherParam struct {
	DeviceID   int64       `json:"device_id"`
	ActuatorID int64       `json:"actuator_id"`
	Value      interface{} `json:"value"`
	ActionBy   constants.ActionSource
}

type InsertActionLogParam struct {
	ActuatorID int64
	Value      interface{}
	ActionBy   constants.ActionSource
	ActionAt   time.Time
}
