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
	dispatch, found := uc.actionDispatcherMap[param.ActionType]
	if !found {
		return errorInvalidActionType
	}

	defer uc.insertActionLog(InsertActionLogParam{
		DeviceID:   param.DeviceID,
		ActionType: param.ActionType,
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
	}

	result := make([]ActionHistory, 0, len(histories))
	for _, history := range histories {
		actionTime := history.ActionAt
		result = append(result, ActionHistory{
			ActionType: history.ActionType,
			Value:      parseActionValue(history.ActionType, history.Value),
			ActionBy:   history.ActionBy,
			ActionAt:   &actionTime,
		})
	}

	return result, nil
}

func (uc *Usecase) GetAvailableActions(ctx context.Context) ([]DispatcherAction, error) {
	return []DispatcherAction{
		{Name: "Built In LED", Action: constants.ActionTypeBuiltInLED},
		{Name: "Relay", Action: constants.ActionTypeRelay},
	}, nil
}

func (uc *Usecase) GetLastAction(deviceID int64, actionType constants.ActionType) (ActionHistory, error) {
	history, err := uc.db.GetLastActionByActionType(deviceID, actionType)
	if err != nil {
		return ActionHistory{}, err
	}

	if history.ActionAt.IsZero() {
		return ActionHistory{}, nil
	}

	return ActionHistory{
		Value:    parseActionValue(actionType, history.Value),
		ActionAt: &history.ActionAt,
	}, nil
}

func (uc *Usecase) insertActionLog(param InsertActionLogParam) error {
	return uc.db.InsertActionLog(db.ActionHistory(param))
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
