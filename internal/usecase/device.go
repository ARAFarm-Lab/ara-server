package usecase

import (
	"ara-server/internal/constants"
	"ara-server/util/log"
	"strconv"
	"time"
)

func (uc *Usecase) InitiateDeviceState(deviceID int64) {
	histories, err := uc.db.GetLastActions(deviceID)
	if err != nil {
		log.Error(err, "error getting last action", deviceID)
		return
	}

	result := make([][]interface{}, 0, len(histories))
	for _, history := range histories {
		result = append(result, []interface{}{history.ActionType, history.Value})
	}

	uc.mq.PublishJSON("dcs-"+strconv.FormatInt(deviceID, 10), result)
}

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

	return uc.mq.PublishJSON(generateDeviceTopic(param.DeviceID), []interface{}{constants.ActionTypeRelay, value})
}
