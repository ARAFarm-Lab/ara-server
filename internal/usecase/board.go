package usecase

import (
	"ara-server/internal/constants"
	"ara-server/util/log"
	"context"
	"time"
)

func (uc *Usecase) toggleBuiltInLED(ctx context.Context, param DispatcherParam) error {
	value, ok := param.Value.(bool)
	if !ok {
		log.Error(ctx, param.Value, errorInvalidActionValue, "invalid built in LED action value")
		return errorInvalidActionValue
	}

	defer uc.insertActionLog(InsertActionLogParam{
		DeviceID:   param.DeviceID,
		ActionType: param.ActionType,
		Value:      param.Value,
		ActionBy:   constants.ActionSourceUser,
		ActionAt:   time.Now(),
	})

	return uc.mq.PublishJSON(generateDeviceTopic(param.DeviceID), []interface{}{constants.ActionTypeBuiltInLED, value})
}
