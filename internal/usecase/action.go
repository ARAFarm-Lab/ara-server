package usecase

import (
	"ara-server/internal/constants"
	"ara-server/internal/repository/db"
	"ara-server/util/log"
	"context"
	"strconv"
	"time"
)

func (uc *Usecase) DispatchAction(ctx context.Context, param DispatcherParam) error {
	actuator, err := uc.db.GetActuatorByID(ctx, param.ActuatorID)
	if err != nil {
		log.Error(ctx, param, err, "failed getting actuator")
		return err
	}

	if actuator.ID == 0 {
		return errorActuatorNotFound
	}

	dispatch, found := uc.actionDispatcherMap[actuator.ActionType]
	if !found {
		return errorInvalidActionType
	}

	defer uc.insertActionLog(InsertActionLogParam{
		ActuatorID: param.ActuatorID,
		Value:      param.Value,
		ActionBy:   param.ActionBy,
		ActionAt:   time.Now(),
	})

	return dispatch(ctx, param)
}

func (uc *Usecase) GetActionHistories(ctx context.Context, deviceID int64) ([]ActionHistory, error) {
	histories, err := uc.db.GetActionHistories(ctx, deviceID)
	if err != nil {
		log.Error(ctx, deviceID, err, "failed to get action histories")
		return nil, err
	}

	result := make([]ActionHistory, 0, len(histories))
	for _, history := range histories {
		actionTime := history.ActionAt
		result = append(result, ActionHistory{
			Value:          parseActionValue(history.ActionType, history.Value),
			ActionBy:       history.ActionBy,
			ActionExecutor: history.ActionExecutor,
			ActionAt:       &actionTime,
			Action: DispatcherAction{
				ID:   history.ActuatorID,
				Type: history.ActionType,
				Name: history.Name,
				Icon: history.Icon,
			},
		})
	}

	return result, nil
}

func (uc *Usecase) GetAvailableActions(ctx context.Context, deviceID int64) ([]DispatcherAction, error) {
	actuators, err := uc.db.GetActuatorsByFilter(ctx, []db.GetActuatorsFilter{
		{
			Name:  "is_active",
			Value: true,
		},
	})
	if err != nil {
		log.Error(ctx, deviceID, err, "failed getting actuators")
		return nil, err
	}

	result := make([]DispatcherAction, 0, len(actuators))
	for _, actuator := range actuators {
		result = append(result, DispatcherAction{
			ID:   actuator.ID,
			Type: actuator.ActionType,
			Name: actuator.Name,
			Icon: actuator.Icon,
		})
	}
	return result, nil
}

func (uc *Usecase) GetLastAction(deviceID int64, actuatorID int64) (ActionHistory, error) {
	history, err := uc.db.GetLastActionByActuatorID(deviceID, actuatorID)
	if err != nil {
		return ActionHistory{}, err
	}

	if history.ActionAt.IsZero() {
		return ActionHistory{}, nil
	}

	return ActionHistory{
		Value:    parseActionValue(history.ActionType, history.Value),
		ActionAt: &history.ActionAt,
	}, nil
}

func (uc *Usecase) insertActionLog(param InsertActionLogParam) error {
	return uc.db.InsertActionLog(db.ActionHistory{
		ActuatorID: param.ActuatorID,
		Value:      param.Value,
		ActionBy:   param.ActionBy,
		ActionAt:   param.ActionAt,
	})
}

func parseActionValue(actionType constants.ActionType, value interface{}) interface{} {
	switch actionType {
	case constants.ActionTypeBuiltInLED, constants.ActionTypeRelay:
		val, ok := value.(string)
		if !ok {
			return false
		}
		if v, err := strconv.ParseBool(val); err == nil {
			return v
		}
		return false
	}

	return nil
}
