package usecase

import (
	"ara-server/internal/constants"
	"ara-server/util/log"
	"context"
	"strconv"
	"time"
)

func (uc *Usecase) InitiateDeviceState(ctx context.Context, deviceID int64) {
	histories, err := uc.db.GetLastActions(deviceID)
	if err != nil {
		log.Error(ctx, deviceID, err, "error getting last action")
		return
	}

	result := make([][]interface{}, 0, len(histories))
	for _, history := range histories {
		result = append(result, []interface{}{history.ActionType, history.Value})
	}

	uc.mq.PublishJSON("dcs-"+strconv.FormatInt(deviceID, 10), result)
}

func (uc *Usecase) toggleRelay(ctx context.Context, param DispatcherParam) error {
	value, ok := param.Value.(bool)
	if !ok {
		log.Error(ctx, param.Value, errorInvalidActionValue, "invalid relay action value")
		return errorInvalidActionValue
	}

	defer uc.insertActionLog(InsertActionLogParam{
		DeviceID:   param.DeviceID,
		ActionType: param.ActionType,
		Value:      param.Value,
		ActionBy:   constants.ActionSourceUser,
		ActionAt:   time.Now(),
	})

	return uc.mq.PublishJSON(generateDeviceTopic(param.DeviceID), []interface{}{constants.ActionTypeRelay, value})
}
