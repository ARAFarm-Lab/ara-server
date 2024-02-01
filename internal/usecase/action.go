package usecase

import (
	"ara-server/internal/constants"
	"ara-server/internal/repository/db"
	"context"
	"strconv"
)

func (uc *Usecase) DispatchAction(ctx context.Context, param DispatcherParam) error {
	dispatch, found := uc.actionDispatcherMap[param.ActionType]
	if !found {
		return errorInvalidActionType
	}

	return dispatch(ctx, param)
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
