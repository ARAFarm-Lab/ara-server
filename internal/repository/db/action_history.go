package db

import (
	"ara-server/internal/constants"
	"database/sql"
	"fmt"
)

func (repo *Repository) GetLastActions(deviceID int64) ([]ActionHistory, error) {
	query := fmt.Sprintf(queryGetLastAction, "device_id = $1 ")

	var actionHistories []ActionHistory
	err := repo.db.Select(&actionHistories, query, deviceID)
	if err != nil && err != sql.ErrNoRows {
		return actionHistories, err
	}

	return actionHistories, nil
}

func (repo *Repository) GetLastActionByActionType(deviceID int64, actionType constants.ActionType) (ActionHistory, error) {
	query := fmt.Sprintf(queryGetLastAction, "device_id = $1 AND action_type = $2 ")

	var actionHistory ActionHistory
	err := repo.db.Get(&actionHistory, query, deviceID, actionType)
	if err != nil && err != sql.ErrNoRows {
		return actionHistory, err
	}

	return actionHistory, nil
}

func (repo *Repository) InsertActionLog(param ActionHistory) error {
	_, err := repo.db.Exec(queryInsertActionLog, param.DeviceID, param.ActionType, param.Value, param.ActionBy, param.ActionAt)
	return err
}
