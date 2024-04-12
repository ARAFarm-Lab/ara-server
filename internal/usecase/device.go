package usecase

import (
	"ara-server/internal/constants"
	"ara-server/util/log"
	"context"
	"strconv"
)

const (
	mqTopicAction                     = "action/"
	mqTopicInitialDeviceStateResponse = "device-initial-state-response/"
	mqTopicRestartDevice              = "restart-device/"
	mqTopicHeartbeatResponse          = "heartbeat-request/"
)

func (uc *Usecase) InitiateDeviceState(ctx context.Context, deviceID int64) {
	histories, err := uc.db.GetLastActions(deviceID)
	if err != nil {
		log.Error(ctx, deviceID, err, "failed getting last action")
		return
	}

	actuators, err := uc.db.GetActuatorsByFilter(ctx, nil)
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

	uc.mq.PublishJSON(mqTopicInitialDeviceStateResponse+strconv.FormatInt(deviceID, 10), result)
}

func (uc *Usecase) RestartDevice(ctx context.Context, deviceID int64) {
	uc.mq.PublishJSON(mqTopicRestartDevice+strconv.FormatInt(deviceID, 10), nil)
}

func (uc *Usecase) SendHeartbeat(ctx context.Context, deviceID int64) {
	uc.mq.PublishJSON(mqTopicHeartbeatResponse+strconv.FormatInt(deviceID, 10), nil)
}

func (uc *Usecase) toggleBuiltInLED(ctx context.Context, param DispatcherParam) error {
	value, ok := param.Value.(bool)
	if !ok {
		log.Error(ctx, param.Value, errorInvalidActionValue, "invalid built in LED action value")
		return errorInvalidActionValue
	}

	return uc.mq.PublishJSON(mqTopicAction+strconv.FormatInt(param.DeviceID, 10), []interface{}{constants.ActionTypeBuiltInLED, value})
}

func (uc *Usecase) toggleRelay(ctx context.Context, param DispatcherParam) error {
	value, ok := param.Value.(bool)
	if !ok {
		log.Error(ctx, param.Value, errorInvalidActionValue, "invalid relay action value")
		return errorInvalidActionValue
	}

	return uc.mq.PublishJSON(mqTopicAction+strconv.FormatInt(param.DeviceID, 10), []interface{}{param.ActuatorID, value})
}
