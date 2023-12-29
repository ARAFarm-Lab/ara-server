package http

import "ara-server/internal/constants"

type DispatcherRequest struct {
	DeviceID   int64                `json:"device_id"`
	ActionType constants.ActionType `json:"action_type"`
	Value      interface{}          `json:"value"`
}
