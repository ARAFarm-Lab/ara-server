package usecase

import (
	"ara-server/internal/repository/db"
	"ara-server/util/log"
	"context"
	"strconv"
)

func (uc *Usecase) InitiateDeviceState(ctx context.Context, deviceID int64) {
	histories, err := uc.db.GetLastActions(deviceID)
	if err != nil {
		log.Error(ctx, deviceID, err, "failed getting last action")
		return
	}

	actuators, err := uc.db.GetActuatorsByFilter(ctx, []db.GetActuatorsFilter{
		{
			Name:  "is_active",
			Value: true,
		},
	})
	if err != nil {
		log.Error(ctx, deviceID, err, "failed getting actuator list")
		return
	}

	states := make([][]interface{}, 0, len(histories))
	for _, history := range histories {
		states = append(states, []interface{}{history.ActuatorID, history.Value})
	}

	pins := make([][]interface{}, 0, len(actuators))
	for _, actuator := range actuators {
		pins = append(pins, []interface{}{actuator.ID, actuator.PinNumber})
	}

	result := deviceConfig{
		Actuators: pins,
		Values:    states,
	}

	uc.mq.PublishJSON("dcs-"+strconv.FormatInt(deviceID, 10), result)
}

func (uc *Usecase) toggleRelay(ctx context.Context, param DispatcherParam) error {
	value, ok := param.Value.(bool)
	if !ok {
		log.Error(ctx, param.Value, errorInvalidActionValue, "invalid relay action value")
		return errorInvalidActionValue
	}

	return uc.mq.PublishJSON(generateDeviceTopic(param.DeviceID), []interface{}{param.ActuatorID, value})
}
