package usecase

import (
	"ara-server/internal/constants"
	"ara-server/internal/repository/mq"
	"ara-server/util/log"
	"time"
)

func (uc *Usecase) toggleRelay(param DispatcherParam) error {
	value, ok := param.Value.(bool)
	if !ok {
		log.Error(errorInvalidActionValue, "invalid relay action value", param.Value)
		return errorInvalidActionValue
	}

	defer uc.insertActionLog(InsertActionLogParam{
		DeviceID:   param.DeviceID,
		ActionType: param.ActionType,
		Value:      param.Value,
		ActionBy:   constants.ActionSourceUser,
		ActionAt:   time.Now(),
	})

	return uc.mq.PublishJSON(generateDeviceTopic(param.DeviceID), mq.PublishJSONPayload{
		ActionType: constants.ActionTypeRelay,
		Value:      value,
	})
}
