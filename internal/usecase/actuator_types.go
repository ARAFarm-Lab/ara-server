package usecase

import "ara-server/internal/constants"

type Actuator struct {
	ID             int64                `json:"id,omitempty"`
	DeviceID       int64                `json:"device_id"`
	PinNumber      int                  `json:"pin_number"`
	ActionType     constants.ActionType `json:"action_type"`
	TerminalNumber int                  `json:"terminal_number"`
	Name           string               `json:"name"`
	Icon           string               `json:"icon"`
	IsActive       bool                 `json:"is_active,omitempty"`
}
