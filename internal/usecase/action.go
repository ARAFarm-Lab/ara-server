package usecase

import (
	"ara-server/internal/constants"
	"ara-server/internal/repository/db"
)

func (uc *Usecase) DispatchAction(param DispatcherParam) error {
	return uc.actionDispatcherMap[param.ActionType](param)
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
		Value:    history.Value,
		ActionAt: &history.ActionAt,
	}, nil
}

func (uc *Usecase) insertActionLog(param InsertActionLogParam) error {
	return uc.db.InsertActionLog(db.ActionHistory(param))
}
