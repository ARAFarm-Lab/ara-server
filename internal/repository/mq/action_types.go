package mq

import "ara-server/internal/constants"

type PublishJSONPayload struct {
	ActionType constants.ActionType
	Value      interface{}
}
